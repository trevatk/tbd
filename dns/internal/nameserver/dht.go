package nameserver

import (
	"cmp"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/trevatk/tbd/lib/protocol"
	pb "github.com/trevatk/tbd/lib/protocol/dns/kademlia/v1"
)

const (
	kademliaK   = 3 // replication factor
	alphaK      = kademliaK
	nodeLength  = sha1.Size
	bitsInBytes = 8
)

type nodeID [nodeLength]byte

func newNodeID(s string) nodeID {
	return nodeID(sha1.Sum([]byte(s)))
}

func nodeIDFromStr(s string) (nodeID, error) {
	decoded, err := hex.DecodeString(s)
	if err != nil {
		return nodeID{}, fmt.Errorf("hex failed to decode %s: %w", s, err)
	}
	if len(decoded) != nodeLength {
		return nodeID{}, fmt.Errorf("unexpected string length %d expected %d", len(decoded), nodeLength)
	}
	var id nodeID
	copy(id[:], decoded)
	return id, nil
}

func (a nodeID) xor(b nodeID) *big.Int {
	distance := make([]byte, nodeLength)
	for i := range a {
		distance[i] = a[i] ^ b[i]
	}
	return new(big.Int).SetBytes(distance)
}

func (n nodeID) toString() string {
	return hex.EncodeToString(n[:])
}

type node struct {
	id       nodeID
	ipOrHost string
	port     uint32
	lastSeen time.Time
}

func newNode(host string, port uint32, lastSeen *time.Time) *node {
	ls := time.Now()
	if lastSeen != nil {
		ls = *lastSeen
	}

	return &node{
		id:       newNodeID(host),
		ipOrHost: host,
		port:     port,
		lastSeen: ls,
	}
}

func (n node) addr() string {
	return net.JoinHostPort(n.ipOrHost, fmt.Sprintf("%d", n.port))
}

type kBucket struct {
	mu sync.RWMutex

	contacts []*node
}

type kademlia struct {
	self *node

	mu           sync.RWMutex
	routingTable []*kBucket

	store kv
}

// interface compliance
var _ dht = (*kademlia)(nil)

// NewDHT return new kademlia implementation of dht
func NewDHT(kv kv, ipOrHost, port string) dht {
	totalKBuckets := nodeLength * bitsInBytes
	routingTable := make([]*kBucket, totalKBuckets)
	for i := range totalKBuckets {
		routingTable[i] = &kBucket{
			mu:       sync.RWMutex{},
			contacts: make([]*node, 0),
		}
	}
	return &kademlia{
		mu:           sync.RWMutex{},
		self:         newNode(ipOrHost, 53, nil),
		store:        kv,
		routingTable: routingTable,
	}
}

func (ka *kademlia) getSelf() *node {
	if ka.self == nil {
		return nil
	}
	return ka.self
}

// bootstrapping a node
// this is the process for a node to join an existing network
// the workflow is entirely reliant on the new node
//
// new node : B
// existing node: A
//
// node A pings node B
// when node a receives a ping message from node B
// node A will update its routing table

// the announcement of the node B joining the network
// is node B triggering a FIND_VALUE gRPC call of itself
func (ka *kademlia) _(ctx context.Context, host string, port uint32) error {
	nn := newNode(host, port, nil)

	// ping existing node
	if err := pingNodeRPC(ctx, ka.self, nn.addr()); err != nil {
		return fmt.Errorf("failed to ping node: %w", err)
	}

	err := ka.addNode(ctx, nn)
	if err != nil {
		return fmt.Errorf("failed to add node: %w", err)
	}

	selfAsTarget := ka.self.id
	if _, err := ka.findNode(ctx, selfAsTarget); err != nil {
		return fmt.Errorf("failed to find node: %w", err)
	}

	return nil
}

// iterative
func (ka *kademlia) findNode(ctx context.Context, targetID nodeID) ([]*node, error) {
	var (
		discoveredNodes   = make(map[nodeID]*node)
		discoveredNodesMu = &sync.Mutex{}

		nodesToQuery = make([]*node, 0)
		queriedNodes = make(map[nodeID]struct{})
	)

	cns := ka.findClosestNodes(targetID)
	for _, cn := range cns {
		discoveredNodes[cn.id] = cn
		nodesToQuery = append(nodesToQuery, cn)
	}

	for i := 0; i <= kademliaK; i++ {
		if len(nodesToQuery) == 0 {
			// no new nodes to query for this iteration
			break
		}

		alphaNodes := make([]*node, 0)
		// iterate over nodes to query
		// verify not previously queried
		// add to alpha list to be queried and set queried with struct
		for _, qn := range nodesToQuery {
			if _, queried := queriedNodes[qn.id]; !queried {
				alphaNodes = append(alphaNodes, qn)
				queriedNodes[qn.id] = struct{}{}
			}
		}

		if len(alphaNodes) == 0 {
			// no new alpha nodes to query with this iteration
			break
		}

		var (
			g, ctx = errgroup.WithContext(ctx)
		)

		for _, an := range alphaNodes {
			n := an
			g.Go(func() error {
				peerNodes, err := findNodeRPC(ctx, ka.self, targetID, n.addr())
				if err != nil {
					return nil
				}

				for _, pn := range peerNodes {
					_ = ka.addNode(ctx, pn)

					discoveredNodesMu.Lock()
					if _, exists := discoveredNodes[pn.id]; !exists {
						discoveredNodes[pn.id] = pn
					}
					discoveredNodesMu.Unlock()
				}

				return nil
			})
		}

		err := g.Wait()
		if err != nil {
			return nil, err
		}

		candidates := make([]*node, 0, len(discoveredNodes))
		discoveredNodesMu.Lock()
		for _, dc := range discoveredNodes {
			candidates = append(candidates, dc)
		}
		discoveredNodesMu.Unlock()

		sort.Slice(candidates, func(i, j int) bool {
			a := candidates[i].id.xor(targetID)
			b := candidates[j].id.xor(targetID)
			return a.Cmp(b) == -1
		})

		nodesToQuery = make([]*node, 0)
		for _, c := range candidates {
			if _, queried := queriedNodes[c.id]; !queried {
				nodesToQuery = append(nodesToQuery, c)
			}
			if len(nodesToQuery) >= kademliaK {
				break
			}
		}
	}

	finalClosestNodes := make([]*node, 0, len(discoveredNodes))
	discoveredNodesMu.Lock()
	for _, n := range discoveredNodes {
		finalClosestNodes = append(finalClosestNodes, n)
	}
	discoveredNodesMu.Unlock()

	sort.Slice(finalClosestNodes, func(i, j int) bool {
		a := finalClosestNodes[i].id.xor(targetID)
		b := finalClosestNodes[j].id.xor(targetID)
		return a.Cmp(b) == -1
	})

	return finalClosestNodes[:minGN(len(finalClosestNodes), kademliaK)], nil
}

func (ka *kademlia) findClosestNodes(targetID nodeID) []*node {
	ka.mu.RLock()
	defer ka.mu.RUnlock()

	var (
		candidates = make([]*node, 0)
	)

	// iterate over all nodes in every bucket
	// add them all as candidates to be sorted later
	for _, kb := range ka.routingTable {
		kb.mu.RLock()
		candidates = append(candidates, kb.contacts...)
		kb.mu.RUnlock()
	}

	if len(candidates) == 0 {
		return []*node{}
	}

	// sort all nodes based on XOR value
	sort.Slice(candidates, func(i, j int) bool {
		a := candidates[i].id.xor(targetID)
		b := candidates[j].id.xor(targetID)
		return a.Cmp(b) == -1
	})

	// determine the nodes of nodes to return
	// if number of nodes is less than K
	// return min number nodes or K
	numNodes := minGN(len(candidates), kademliaK)
	return candidates[:numNodes]
}

func (ka *kademlia) getValue(key string) (*record, error) {
	ka.mu.RLock()
	defer ka.mu.RUnlock()
	return ka.store.get(key)
}

func (ka *kademlia) setValue(key string, value *record) error {
	ka.mu.Lock()
	defer ka.mu.Unlock()
	return ka.store.set(key, value)
}

func (ka *kademlia) addNode(ctx context.Context, node *node) error {
	ka.mu.Lock()
	defer ka.mu.Unlock()

	bucketIndex := getBucketIndex(ka.self.id, node.id)
	if bucketIndex < 0 || bucketIndex > len(ka.routingTable) {
		if bucketIndex == -1 {
			// attempting to insert self
			// kademlia does not allow insertion of
			// self record in kbuckets
			return nil
		}
		return fmt.Errorf("invalid bucket index: %d for node %s", bucketIndex, node.id.toString())
	}

	return ka.routingTable[bucketIndex].addNode(ctx, ka.self, node)
}

func (ka *kademlia) findValue(ctx context.Context, targetID nodeID) (*record, []*node, error) {
	value, err := ka.store.get(targetID.toString())
	if err == nil {
		return value, nil, nil
	} else if !errors.Is(err, errKeyNotFound) {
		return nil, nil, fmt.Errorf("failed to get record from store: %w", err)
	}

	var (
		discoveredNodes   = make(map[nodeID]*node)
		discoveredNodesMu = &sync.Mutex{}
		nodesToQuery      = make([]*node, 0)
		queriedNodes      = make(map[nodeID]struct{})
		foundRecord       *record
		foundRecordMu     = &sync.Mutex{}
	)

	closestNodes := ka.findClosestNodes(targetID)
	for _, cn := range closestNodes {
		discoveredNodes[cn.id] = cn
		nodesToQuery = append(nodesToQuery, cn)
	}

	// begin iterative lookup
	for i := 0; i <= kademliaK; i++ {
		if len(nodesToQuery) == 0 {
			// no nodes to query
			break
		}

		alphaNodes := make([]*node, 0)
		for _, qn := range nodesToQuery {
			if _, queried := queriedNodes[qn.id]; !queried {
				alphaNodes = append(alphaNodes, qn)
				if len(alphaNodes) > alphaK {
					break
				}
			}
		}

		if len(alphaNodes) == 0 {
			break
		}

		g, gCtx := errgroup.WithContext(ctx)

		for _, an := range alphaNodes {
			closestNode := an

			g.Go(func() error {
				// set node as queried
				queriedNodes[an.id] = struct{}{}

				value, closestFromPeer, err := findValueRPC(gCtx, ka.self, targetID.toString(), closestNode.addr())
				if err != nil {
					return fmt.Errorf("failed to execute find_value gRPC: %w", err)
				}

				if value != nil {
					// value found
					foundRecordMu.Lock()
					if foundRecord == nil {
						foundRecord = value
					}
					foundRecordMu.Unlock()
					return nil
				}

				for _, pn := range closestFromPeer {
					err = ka.addNode(gCtx, pn)
					if err != nil {
						return fmt.Errorf("failed to add node: %w", err)
					}
					discoveredNodesMu.Lock()
					if _, exists := discoveredNodes[pn.id]; !exists {
						discoveredNodes[pn.id] = pn
					}
					discoveredNodesMu.Unlock()
				}

				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return nil, nil, fmt.Errorf("go func failed during errgroup execution: %w", err)
		}

		foundRecordMu.Lock()
		if foundRecord != nil {
			return foundRecord, nil, nil
		}
		foundRecordMu.Unlock()

		nodesToQuery = make([]*node, 0)
		discoveredNodesMu.Lock()
		candidates := make([]*node, 0, len(discoveredNodes))
		for _, n := range discoveredNodes {
			candidates = append(candidates, n)
		}
		discoveredNodesMu.Unlock()

		sort.Slice(candidates, func(i, j int) bool {
			a := candidates[i].id.xor(targetID)
			b := candidates[j].id.xor(targetID)
			return a.Cmp(b) == -1
		})

		for _, c := range candidates {
			if _, queried := queriedNodes[c.id]; !queried {
				nodesToQuery = append(nodesToQuery, c)
			}
			if len(nodesToQuery) >= kademliaK {
				break
			}
		}
	}

	return nil, nil, errKeyNotFound
}

func (kb *kBucket) addNode(ctx context.Context, n, newNode *node) error {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	if len(kb.contacts) == 0 {
		kb.contacts = append(kb.contacts, newNode)
		return nil
	}

	// if maximum number of replication is reached per bucket
	// then ping the oldest node
	//
	// if oldest node does not respond
	// then remove the oldest node
	if len(kb.contacts) >= kademliaK {
		sort.Slice(kb.contacts, func(i, j int) bool {
			return kb.contacts[i].lastSeen.Before(kb.contacts[j].lastSeen)
		})

		oldestIndex := len(kb.contacts) - 1
		oldestNode := kb.contacts[oldestIndex]

		if err := pingNodeRPC(ctx, n, oldestNode.addr()); err != nil {
			// failed to ping oldest node
			// replace with new node
			kb.contacts[oldestIndex] = newNode
			return nil
		}

		// oldest node is still active updated last seen
		kb.contacts[oldestIndex].lastSeen = time.Now()

		return nil
	}

	// kbucket is not full
	// new nodes can be appended to list
	kb.contacts = append(kb.contacts, newNode)

	return nil
}

func pingNodeRPC(ctx context.Context, sender *node, target string) error {
	conn, err := protocol.NewConn(target)
	if err != nil {
		return fmt.Errorf("failed to create new client connection: %w", err)
	}
	defer func() { _ = conn.Close() }()

	requestID := uuid.New().String()
	resp, err := pb.NewKademliaServiceClient(conn).Ping(ctx, &pb.PingRequest{
		Sender:    nodeToSender(sender),
		RequestId: requestID,
	})
	if err != nil {
		return fmt.Errorf("failed to execute ping command: %w", err)
	}

	if resp.RequestId != requestID {
		return errInvalidRequestID
	}

	return nil
}

func findValueRPC(ctx context.Context, sender *node, key, target string) (*record, []*node, error) {
	// g, ctx := errgroup.WithContext(ctx)
	conn, err := protocol.NewConn(target)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create client connection: %w", err)
	}
	defer func() { _ = conn.Close() }()

	requestID := uuid.New()

	resp, err := pb.NewKademliaServiceClient(conn).FindValue(ctx, &pb.FindValueRequest{
		Sender:    nodeToSender(sender),
		RequestId: requestID.String(),
		Key:       key,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute find value gRPC call: %w", err)
	}

	switch result := resp.Result.(type) {
	case *pb.FindValueResponse_ClosestNodes:
		ns := make([]*node, 0, len(result.ClosestNodes.Nodes))
		for _, cn := range result.ClosestNodes.Nodes {
			ns = append(ns, newNode(cn.IpOrDomain, cn.Port, nil))
		}
		return nil, ns, nil
	case *pb.FindValueResponse_Record:
		return &record{
			domain:     result.Record.Domain,
			recordType: pbToRecordType(result.Record.RecordType),
			value:      result.Record.Value,
			ttl:        result.Record.Ttl,
		}, nil, nil
	default:
		return nil, nil, errors.New("unsupported response type")
	}
}

func findNodeRPC(ctx context.Context, sender *node, targetID nodeID, target string) ([]*node, error) {
	var (
		discoveredContacts = make([]*node, 0)
	)

	conn, err := protocol.NewConn(target)
	if err != nil {
		return nil, fmt.Errorf("failed to create client connection: %w", err)
	}
	defer func() { _ = conn.Close() }()

	requestID := uuid.New().String()

	request := &pb.FindNodeRequest{
		Sender:       nodeToSender(sender),
		RequestId:    requestID,
		TargetNodeId: targetID.toString(),
	}

	resp, err := pb.NewKademliaServiceClient(conn).FindNode(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to execute find node: %w", err)
	}

	if resp.RequestId != requestID {
		return nil, errInvalidRequestID
	}

	for _, cn := range resp.ClosestNodes {
		nodeID, err := nodeIDFromStr(cn.NodeId)
		if err != nil {
			return nil, fmt.Errorf("node id from string: %w", err)
		}
		discoveredContacts = append(discoveredContacts, &node{
			id:       nodeID,
			ipOrHost: cn.IpOrDomain,
			port:     cn.Port,
			lastSeen: cn.LastSeen.AsTime(),
		})
	}

	return discoveredContacts, nil
}

// calcuate the which kbucket the B nodeID
// belongs int
func getBucketIndex(a, b nodeID) int {
	// distance between A and B
	distance := a.xor(b)
	// leading zeros
	// if distance is zero then node
	// is same as self
	// return -1 to not add self
	if distance.Cmp(big.NewInt(0)) == 0 {
		return -1
	}
	// number of bits in kademlia
	// sha1 hash is 160 bits
	totalBits := nodeLength * bitsInBytes // 160
	// calculate number of leading zeros based on
	// the bits used to represent absolute value of
	// distance
	//
	// the larger the distance the more bits needed
	// to represent value
	// leading to the posssibility of another bucket
	// being used
	//
	// ex
	// x : totalBits = 160
	// y: distance = 1
	// x - y = 159
	leadingZeros := totalBits - distance.BitLen()
	// calculate index based on bits - leadingZeros
	// totalBits - 1 is used because of zero index
	index := totalBits - 1 - leadingZeros
	return index
}

// generic func to get the lowest value
func minGN[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
