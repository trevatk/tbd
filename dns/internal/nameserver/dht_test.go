package nameserver

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	host0 = "127.0.0.1"
	host1 = "127.0.0.2"
	host2 = "127.0.0.3"
	host3 = "127.0.0.4"

	// this is dumb but only happens from
	// dht constructor
	portUint32 uint32 = 53
	port              = fmt.Sprintf("%d", portUint32)

	id0 = newNodeID(host0)
	id1 = newNodeID(host1)
	id2 = newNodeID(host2)
	id3 = newNodeID(host3)
)

func TestFindClosestNodes(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	assert := assert.New(t)

	t.Run("empty", func(t *testing.T) {
		ctlr := gomock.NewController(t)
		mockKv := NewMockkv(ctlr)

		var (
			expected = 0
		)

		dht := NewDHT(mockKv, host0, port)
		ns := dht.findClosestNodes(id1)
		assert.Equal(expected, len(ns))
	})

	t.Run("1", func(t *testing.T) {
		ctlr := gomock.NewController(t)
		mockKv := NewMockkv(ctlr)

		var (
			expected int = 1
		)

		dht := NewDHT(mockKv, host1, port)
		assert.NoError(dht.addNode(ctx, newNode(host2, portUint32, nil)))

		ns := dht.findClosestNodes(id2)
		assert.Equal(expected, len(ns))
	})
}

func TestAddNode(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ctlr := gomock.NewController(t)
	mockKv := NewMockkv(ctlr)

	assert := assert.New(t)

	dht := NewDHT(mockKv, host1, port)

	t.Run("success", func(t *testing.T) {
		var (
			expected error = nil
		)
		err := dht.addNode(ctx, newNode(host2, portUint32, nil))
		assert.Equal(expected, err)
	})
}

func TestFindNode(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ctlr := gomock.NewController(t)
	mockKv := NewMockkv(ctlr)

	assert := assert.New(t)

	dht := NewDHT(mockKv, host0, port)

	t.Run("1", func(t *testing.T) {
		var (
			expectedLen = 1
		)

		nn := newNode(host1, portUint32, nil)
		assert.NoError(dht.addNode(ctx, nn))

		nodes, err := dht.findNode(ctx, id1)
		assert.NoError(err)

		assert.Equal(expectedLen, len(nodes))
	})

	t.Run("2", func(t *testing.T) {
		var (
			expectedLen = 2
		)

		nn := newNode(host2, portUint32, nil)
		assert.NoError(dht.addNode(ctx, nn))

		nodes, err := dht.findNode(ctx, id2)
		assert.NoError(err)

		assert.Equal(expectedLen, len(nodes))
	})

	t.Run("3", func(t *testing.T) {
		var (
			expectedLen = 3
		)

		assert.NoError(dht.addNode(ctx, newNode(host3, portUint32, nil)))
		nodes, err := dht.findNode(ctx, id3)
		assert.NoError(err)
		assert.Equal(expectedLen, len(nodes))
	})
}

func TestFindValue(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var (
		key = newNodeID("structx.io")
	)
	fmt.Println(id0.toString())
	ctlr := gomock.NewController(t)
	defer ctlr.Finish()

	mockKv := NewMockkv(ctlr)
	mockKv.EXPECT().get(key.toString()).Return(r1, nil).Times(1)
	mockKv.EXPECT().get(key.toString()).Return(nil, errKeyNotFound).Times(1)

	assert := assert.New(t)

	dht := NewDHT(mockKv, host0, port)

	t.Run("from_store", func(t *testing.T) {
		var (
			expected error = nil
		)
		value, _, err := dht.findValue(ctx, key)
		assert.Equal(expected, err)
		assert.Equal(r1, value)
	})

	t.Run("not_found", func(t *testing.T) {
		var (
			expected error = errKeyNotFound
		)
		_, _, err := dht.findValue(ctx, key)
		assert.Equal(expected, err)
	})
}
