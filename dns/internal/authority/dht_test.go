package foundations

import (
	"context"
	"testing"
)

func TestFindClosestNodes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	dht := NewDHT("127.0.0.1", "53")
	t.Run("empty", func(t *testing.T) {
		expected := 0
		ns := dht.findClosestNodes("127.0.0.1:53")
		if len(ns) != expected {
			t.Fatalf("unexpected node length %d expected %d", len(ns), expected)
		}
	})

	n1 := newNode{
		domain:     "ns1.structx.io",
		recordType: "A",
		value:      "127.0.0.1",
		ttl:        -1,
	}
	dht.routingTable[0].addNode(ctx, n1)
	t.Run("1", func(t *testing.T) {
		expected := 1
		ns := dht.findClosestNodes("127.0.0.1:53")
		if len(ns) != expected {
			t.Fatalf("unexpected node length %d expected %d", len(ns), expected)
		}
	})
}
