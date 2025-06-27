package foundations

import (
	"context"
	"log/slog"

	"github.com/structx/tbd/lib/gateway"
	pba "github.com/structx/tbd/lib/protocol/dns/authoritative/v1"
	pbk "github.com/structx/tbd/lib/protocol/dns/kademlia/v1"
	pbr "github.com/structx/tbd/lib/protocol/dns/resolver/v1"
)

type dht interface {
	findClosestNodes()
	addNode()
	storeValue()
	getValue(key string) error
}

type grpcTransport struct {
	pbk.UnimplementedKademliaServiceServer
	pba.UnimplementedAuthoritativeServiceServer
	pbr.UnimplementedDNSResolverServiceServer

	logger *slog.Logger
	k      *kademlia
}

// interface compliance
var _ pbk.KademliaServiceServer = (*grpcTransport)(nil)
var _ pba.AuthoritativeServiceServer = (*grpcTransport)(nil)
var _ pbr.DNSResolverServiceServer = (*grpcTransport)(nil)

// NewTransport
func NewTransport(logger *slog.Logger, kademlia *kademlia) []gateway.Transport {
	tr := &grpcTransport{
		logger: logger,
		k:      kademlia,
	}
	trs := []gateway.Transport{
		{
			ServiceDesc: &pba.AuthoritativeService_ServiceDesc,
			Service:     tr,
		},
		{
			ServiceDesc: &pbk.KademliaService_ServiceDesc,
			Service:     tr,
		},
		{
			ServiceDesc: &pbr.DNSResolverService_ServiceDesc,
			Service:     tr,
		},
	}
	return trs
}

// FindNode
func (t *grpcTransport) FindNode(ctx context.Context, in *pbk.FindNodeRequest) (*pbk.FindNodeResponse, error) {
	return newFindNodeResponse(), nil
}

// FindValue
func (t *grpcTransport) FindValue(ctx context.Context, in *pbk.FindValueRequest) (*pbk.FindValueResponse, error) {
	// t.d.getValue()
	return newFindValueResponse(), nil
}

// Join
func (t *grpcTransport) Join(ctx context.Context, in *pbk.JoinRequest) (*pbk.JoinResponse, error) {
	return newJoinResponse(), nil
}

// Ping
func (t *grpcTransport) Ping(ctx context.Context, in *pbk.PingRequest) (*pbk.PingResponse, error) {
	return newPingResponse(), nil
}

// Store
func (t *grpcTransport) Store(ctx context.Context, in *pbk.StoreRequest) (*pbk.StoreResponse, error) {
	return newStoreResponse(), nil
}

// Resolve
func (t *grpcTransport) Resolve(context.Context, *pbr.ResolveRequest) (*pbr.ResolveResponse, error) {
	return newResolveResponse(), nil
}

func newFindNodeResponse() *pbk.FindNodeResponse {
	return &pbk.FindNodeResponse{}
}
func newFindValueResponse() *pbk.FindValueResponse {
	return &pbk.FindValueResponse{}
}

func newJoinResponse() *pbk.JoinResponse {
	return &pbk.JoinResponse{}
}

func newPingResponse() *pbk.PingResponse {
	return &pbk.PingResponse{}
}

func newStoreResponse() *pbk.StoreResponse {
	return &pbk.StoreResponse{}
}

func newResolveResponse() *pbr.ResolveResponse {
	return &pbr.ResolveResponse{}
}
