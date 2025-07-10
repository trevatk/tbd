package nameserver

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/trevatk/tbd/lib/logging"
	"github.com/trevatk/tbd/lib/protocol"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	pba "github.com/trevatk/tbd/lib/protocol/dns/authoritative/v1"
	pbk "github.com/trevatk/tbd/lib/protocol/dns/kademlia/v1"
)

var (
	n1 = &node{
		id:       id0,
		ipOrHost: host0,
		port:     portUint32,
		lastSeen: time.Now(),
	}
)

func newKademliaTransport(logger *slog.Logger, dht dht) pbk.KademliaServiceServer {
	return &grpcTransport{
		logger: logger,
		dht:    dht,
	}
}

func newFoundationsTransport(logger *slog.Logger, dht dht) pba.AuthoritativeServiceServer {
	return &grpcTransport{
		logger: logger,
		dht:    dht,
	}
}

func TestPing(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	mockDht := NewMockdht(ctrl)
	mockDht.EXPECT().addNode(gomock.Any(), gomock.AssignableToTypeOf(&node{})).Return(nil).AnyTimes()
	mockDht.EXPECT().getSelf().Return(n1).AnyTimes()

	logger := logging.New("DEBUG")

	g := newKademliaTransport(logger, mockDht)

	assert := assert.New(t)

	t.Run("success", func(t *testing.T) {
		var (
			expected error = nil
			r              = &pbk.PingRequest{
				Sender:    nodeToSender(n1),
				RequestId: uuid.New().String(),
			}
		)

		resp, err := g.Ping(ctx, r)
		assert.Equal(expected, err)

		assert.Equal(n1.id.toString(), resp.Sender.NodeId)
		assert.Equal(r.RequestId, resp.RequestId)
		assert.Equal(n1.ipOrHost, resp.Sender.IpOrDomain)
	})
}

// func TestCreateRecordRPC(t *testing.T) {

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	ctrl, ctx := gomock.WithContext(ctx, t)
// 	defer ctrl.Finish()

// 	mockDht := NewMockdht(ctrl)

// 	logger := logging.New("DEBUG")

// 	g := newFoundationsTransport(logger, mockDht)

// 	assert := assert.New(t)

// 	t.Run("success", func(t *testing.T) {

// 	})
// }

func TestCreateZoneRPC(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl, ctx := gomock.WithContext(ctx, t)
	defer ctrl.Finish()

	mockDht := NewMockdht(ctrl)
	mockDht.EXPECT().setValue(gomock.AssignableToTypeOf(""), gomock.AssignableToTypeOf(&record{})).Return(nil).Times(1)
	mockDht.EXPECT().
		setValue(gomock.AssignableToTypeOf(""), gomock.AssignableToTypeOf(&record{})).
		Return(errKeyExists).
		Times(1)

	logger := logging.New("DEBUG")

	g := newFoundationsTransport(logger, mockDht)

	assert := assert.New(t)

	t.Run("success", func(t *testing.T) {
		var (
			expected error = nil
			request        = &pba.CreateZoneRequest{
				DomainOrNamespace: "structx.local",
			}
		)
		resp, err := g.CreateZone(ctx, request)
		assert.Equal(expected, err)
		assert.Equal(request.DomainOrNamespace, resp.Zone.DomainOrNamespace)
		assert.NotEmpty(resp.Zone.CreatedAt)
	})

	t.Run("already_exists", func(t *testing.T) {
		var (
			expected error = protocol.ErrAlreadyExists()
			request        = &pba.CreateZoneRequest{
				DomainOrNamespace: "structx.local",
			}
		)
		_, err := g.CreateZone(ctx, request)
		assert.Equal(expected, err)
	})
}
