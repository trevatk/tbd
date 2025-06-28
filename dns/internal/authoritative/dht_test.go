package authoritative

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFindClosestNodes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	assert := assert.New(t)

	t.Run("empty", func(t *testing.T) {
		ctlr := gomock.NewController(t)
		mockKv := NewMockKv(ctlr)

		expected := 0

		dht := NewDHT(mockKv, "127.0.0.1", "53")
		ns := dht.findClosestNodes("127.0.0.1:53")
		assert.Equal(expected, len(ns))
	})

	t.Run("1", func(t *testing.T) {
		ctlr := gomock.NewController(t)
		mockKv := NewMockKv(ctlr)

		dht := NewDHT(mockKv, "127.0.0.1", "53")

		assert.NoError(dht.routingTable[0].addNode(ctx, dht.self, newNode{ipOrHost: "192.168.1.142", port: "53"}))

		expected := 1
		ns := dht.findClosestNodes("127.0.0.1:53")
		assert.Equal(expected, len(ns))
	})
}
