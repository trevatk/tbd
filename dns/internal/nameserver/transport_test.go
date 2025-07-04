package nameserver

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	"github.com/trevatk/tbd/lib/logging"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	pb "github.com/trevatk/tbd/lib/protocol/dns/kademlia/v1"
)

var (
	n1 = &node{
		id:       id0,
		ipOrHost: host0,
		port:     portUint32,
		lastSeen: time.Now(),
	}
)

func newGrpcTransport(logger *slog.Logger, dht dht) pb.KademliaServiceServer {
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

	g := newGrpcTransport(logger, mockDht)

	assert := assert.New(t)

	t.Run("success", func(t *testing.T) {
		var (
			expected error = nil
			r              = &pb.PingRequest{
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
