package authoritative

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"math/big"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/structx/tbd/lib/protocol"
	pb "github.com/structx/tbd/lib/protocol/dns/kademlia/v1"
)

var (
	errInvalidRequestID = errors.New("invalid request id")
)

const (
	kademliaK   = 3 // replication factor
	nodeLength  = sha1.Size
	bitsInBytes = 8

	numZero = 0
	numNeg1 = -1
)

type nodeID [nodeLength]byte

func newNodeID(s string) nodeID {
	return nodeID(sha1.Sum([]byte(s)))
}

func (a nodeID) xor(b nodeID) *big.Int {
	distance := make([]byte, nodeLength)
	for i := range a {
		distance[i] = a[i] ^ b[i]
	}
	return new(big.Int).SetBytes(distance)
}

type node struct {
	id             nodeID
	ipOrHost, port string
	lastSeen       time.Time
}

type kBucket struct {
	mu *sync.RWMutex

	contacts []*node
}

type kademlia struct {
	self node

	mu           *sync.RWMutex
	routingTable []*kBucket

	kv Kv
}

// NewDHT return new kademlia implementation of dht
func NewDHT(kv Kv, ipOrHost, port string) *kademlia {
	totalKBuckets := nodeLength * bitsInBytes
	routingTable := make([]*kBucket, totalKBuckets)
	for i := range totalKBuckets {
		routingTable[i] = &kBucket{
			mu:       &sync.RWMutex{},
			contacts: make([]*node, numZero),
		}
	}
	hostAndPort := net.JoinHostPort(ipOrHost, port)
	return &kademlia{
		mu: &sync.RWMutex{},
		self: node{
			ipOrHost: ipOrHost,
			port:     port,
			id:       newNodeID(hostAndPort),
		},
		kv:           kv,
		routingTable: routingTable,
	}
}

func (ka *kademlia) findClosestNodes(key string) []*node {
	ka.mu.RLock()
	defer ka.mu.RUnlock()

	var (
		nodeID     = newNodeID(key)
		candidates = make([]*node, numZero)
	)

	// iterate over all nodes in every bucket
	// add them all as candidates to be sorted later
	for _, kb := range ka.routingTable {
		kb.mu.RLock()
		candidates = append(candidates, kb.contacts...)
		kb.mu.RUnlock()
	}

	if len(candidates) == numZero {
		return []*node{}
	}

	// sort all nodes based on XOR value
	sort.Slice(candidates, func(i, j int) bool {
		a := candidates[i].id.xor(nodeID)
		b := candidates[j].id.xor(nodeID)
		return a.Cmp(b) == numNeg1
	})

	// determine the nodes of nodes to return
	// if number of nodes is less than K
	// return min number nodes or K
	numNodes := min(len(candidates), kademliaK)
	return candidates[:numNodes]
}

type newNode struct {
	ipOrHost, port string
}

func (kb *kBucket) addNode(ctx context.Context, n node, newNode newNode) error {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	hostAndPort := net.JoinHostPort(newNode.ipOrHost, newNode.port)
	nodeID := newNodeID(hostAndPort)

	nn := &node{
		id:       nodeID,
		ipOrHost: newNode.ipOrHost,
		port:     newNode.port,
		lastSeen: time.Now(),
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
		target := net.JoinHostPort(oldestNode.ipOrHost, oldestNode.port)

		if err := pingNode(ctx, n, target); err != nil {
			// failed to ping oldest node
			// replace with new node
			kb.contacts[oldestIndex] = nn
			return nil
		}

		// oldest node is still active updated last seen
		kb.contacts[oldestIndex].lastSeen = time.Now()

		return nil
	}

	// kbucket is not full
	// new nodes can be appended to list
	kb.contacts = append(kb.contacts, nn)

	return nil
}

func pingNode(ctx context.Context, sender node, target string) error {
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
