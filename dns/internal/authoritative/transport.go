package authoritative

import (
	"context"
	"encoding/hex"
	"errors"
	"log/slog"

	"github.com/structx/tbd/lib/gateway"
	"github.com/structx/tbd/lib/protocol"

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

// temp fix unused
var _ dht = (nil)

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

// NewTransport return new authoritative transport implementation
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
	err := protocol.Validate(in)
	if err != nil {
		return nil, gateway.ErrInvalidArgument()
	}

	t.logger.DebugContext(ctx, "find_value", slog.Any("request", in))

	value, err := t.k.kv.get(in.Key)
	if err != nil && errors.Is(err, errKeyNotFound) {
		// node does not have value
		// respond with closest nodes
		closestNodes := t.k.findClosestNodes(in.Key)
		return newFindValueResponseWithClosestNodes(t.k.self, closestNodes, in.RequestId), nil
	} else if err != nil {
		t.logger.ErrorContext(ctx, "kv get value", slog.String("error", err.Error()))
		return nil, gateway.ErrInternal()
	}

	return newFindValueResponseWithRecord(t.k.self, value, in.RequestId), nil
}

// Join
func (t *grpcTransport) Join(ctx context.Context, in *pbk.JoinRequest) (*pbk.JoinResponse, error) {
	err := protocol.Validate(in)
	if err != nil {
		return nil, gateway.ErrInvalidArgument()
	}

	t.logger.DebugContext(ctx, "join", slog.Any("request", in))

	return newJoinResponse(), nil
}

// Ping
func (t *grpcTransport) Ping(ctx context.Context, in *pbk.PingRequest) (*pbk.PingResponse, error) {
	err := protocol.Validate(in)
	if err != nil {
		return nil, gateway.ErrInvalidArgument()
	}
	return newPingResponse(t.k.self, in.RequestId), nil
}

// Store
func (t *grpcTransport) Store(ctx context.Context, in *pbk.StoreRequest) (*pbk.StoreResponse, error) {
	err := protocol.Validate(in)
	if err != nil {
		return nil, gateway.ErrInvalidArgument()
	}

	t.logger.DebugContext(ctx, "store_value", slog.Any("request", in))

	if err := t.k.kv.set(in.Record.Id, record{
		domain:     in.Record.Domain,
		recordType: in.Record.RecordType,
		value:      in.Record.Value,
		ttl:        in.Record.Ttl,
	}); err != nil {
		t.logger.ErrorContext(ctx, "kv set value", slog.String("error", err.Error()))
		return nil, gateway.ErrInternal()
	}

	return newStoreResponse(t.k.self, in.RequestId), nil
}

// Resolve
func (t *grpcTransport) Resolve(ctx context.Context, in *pbr.ResolveRequest) (*pbr.ResolveResponse, error) {
	err := protocol.Validate(in)
	if err != nil {
		return nil, gateway.ErrInvalidArgument()
	}

	t.logger.DebugContext(ctx, "resolve", slog.Any("request", in))

	return newResolveResponse(), nil
}

func newFindNodeResponse() *pbk.FindNodeResponse {
	return &pbk.FindNodeResponse{}
}
func newFindValueResponseWithRecord(n node, r record, requestID string) *pbk.FindValueResponse {
	return &pbk.FindValueResponse{
		Sender: nodeToSender(n),
		Result: &pbk.FindValueResponse_Record{
			Record: &pbk.Record{
				Domain:     r.domain,
				RecordType: r.recordType,
				Value:      r.value,
				Ttl:        r.ttl,
			},
		},
		RequestId: requestID,
	}
}

func newFindValueResponseWithClosestNodes(n node, ns []*node, requestID string) *pbk.FindValueResponse {
	closestNodes := make([]*pbk.Node, 0, len(ns))
	for _, n := range ns {
		closestNodes = append(closestNodes, &pbk.Node{
			NodeId:     hex.EncodeToString(n.id[:]),
			IpOrDomain: n.ipOrHost,
			// Port:       n.port,
		})
	}

	return &pbk.FindValueResponse{
		Sender: nodeToSender(n),
		Result: &pbk.FindValueResponse_ClosestNodes{
			ClosestNodes: &pbk.ClosestNodes{
				Nodes: closestNodes,
			},
		},
		RequestId: requestID,
	}
}

func newJoinResponse() *pbk.JoinResponse {
	return &pbk.JoinResponse{}
}

func newPingResponse(n node, requestID string) *pbk.PingResponse {
	return &pbk.PingResponse{
		Sender:    nodeToSender(n),
		RequestId: requestID,
	}
}

func newStoreResponse(n node, requestID string) *pbk.StoreResponse {
	return &pbk.StoreResponse{
		Sender:    nodeToSender(n),
		Success:   true,
		RequestId: requestID,
	}
}

func newResolveResponse() *pbr.ResolveResponse {
	return &pbr.ResolveResponse{
		Status:              pbr.ResolveResponse_RESPONSE_STATUS_SUCCESS,
		AuthoritativeAnswer: true,
	}
}

func nodeToSender(n node) *pbk.Node {
	return &pbk.Node{
		NodeId:     hex.EncodeToString(n.id[:]),
		IpOrDomain: n.ipOrHost,
	}
}
