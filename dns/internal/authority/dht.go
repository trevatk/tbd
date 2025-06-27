package foundations

import (
	"context"
	"crypto/sha1"
	"math/big"
	"net"
	"sort"
	"sync"
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
	h := sha1.New()
	return nodeID(h.Sum([]byte(s)))
}

func (a nodeID) xor(b nodeID) *big.Int {
	distance := make([]byte, nodeLength)
	for i := range a {
		distance[i] = a[i] ^ b[i]
	}
	return new(big.Int).SetBytes(distance)
}

type networkNode struct {
	id   nodeID
	ip   string
	port string
}

type node struct {
	key        nodeID // hashed domain
	domain     string
	recordType string // CNAME, A, MX etc...
	value      string // IP address, CNAME value
	ttl        int

	left, right *node
}

type kBucket struct {
	_    networkNode
	head *node
}

type kademlia struct {
	node networkNode

	mu           *sync.RWMutex
	routingTable []*kBucket
}

// NewDHT return new kademlia implementation of dht
func NewDHT(host, port string) *kademlia {
	totalKBuckets := nodeLength * bitsInBytes
	routingTable := make([]*kBucket, totalKBuckets)
	for i := range totalKBuckets {
		routingTable[i] = &kBucket{
			head: nil,
		}
	}
	hostAndPort := net.JoinHostPort(host, port)
	return &kademlia{
		mu: &sync.RWMutex{},
		node: networkNode{
			ip:   host,
			port: port,
			id:   newNodeID(hostAndPort),
		},
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
		candidates = append(candidates, kb.inOrder()...)
	}

	if len(candidates) == numZero {
		return []*node{}
	}

	// sort all nodes based on XOR value
	sort.Slice(candidates, func(i, j int) bool {
		a := candidates[i].key.xor(nodeID)
		b := candidates[j].key.xor(nodeID)
		return a.Cmp(b) == numNeg1
	})

	// determine the nodes of nodes to return
	// if number of nodes is less than K
	// return min number nodes or K
	numNodes := min(len(candidates), kademliaK)
	return candidates[:numNodes]
}

type newNode struct {
	domain     string
	recordType string
	value      string
	ttl        int
}

func (kb *kBucket) addNode(ctx context.Context, newNode newNode) error {
	nodeID := newNodeID(newNode.domain + newNode.recordType)
	nn := &node{
		key:        nodeID,
		domain:     newNode.domain,
		recordType: newNode.recordType,
		value:      newNode.value,
		ttl:        newNode.ttl,
		left:       nil, right: nil,
	}

	if kb.head == nil {
		kb.head = nn
	}

	n := kb.head
	// for handles insertion into binary search tree
	// when an existing record is found the network_node
	// will need to be pinged to see if it still alive
	//
	// TODO
	// solve eviction policy with kademlia dht when using
	// a bst instead of a linked list to add the last seen
	// at the end of the list
	// this creates a last seen eviction policy
	//
	// this logic could be flawed though as the node
	// is a dns record and logic is when an existing domain
	// and record_type exist
	for {
		distance := n.key.xor(nodeID)

		result := distance.Cmp(big.NewInt(numZero))
		switch result {
		case numZero:

			// if ping fails then insert the new node
			// addr := net.JoinHostPort(kb.nn.ip, kb.nn.port)
			// err := pingNode(ctx, addr)
			// if err != nil {
			// 	n = nn
			// 	return fmt.Errorf("failed to ping node: %w", err)
			// }

			return nil
		case numNeg1:
			if n.left == nil {
				n.left = nn
				return nil
			}
			n = n.left
		default:
			if n.right == nil {
				n.right = nn
				return nil
			}
			n = n.right
		}
	}
}

func (kb *kBucket) inOrder() []*node {
	ns := make([]*node, numZero)
	if kb.head == nil {
		return ns
	}

	var t func(n *node)
	t = func(n *node) {
		if n == nil {
			return
		}
		t(n.left)
		ns = append(ns, n)
		t(n.right)
	}

	t(kb.head)
	return ns
}

// func pingNode(ctx context.Context, target string) error {
// 	conn, err := protocol.NewConn(target)
// 	if err != nil {
// 		return fmt.Errorf("failed to create new client connection: %w", err)
// 	}

// 	resp, err := pb.NewKademliaServiceClient(conn).Ping(ctx, &pb.PingRequest{})
// 	if err != nil {
// 		return fmt.Errorf("failed to execute ping command: %w", err)
// 	}

// 	// TODO
// 	// check if sender exists in nodes
// 	// if not add to node list
// 	_ = resp

// 	return nil
// }
