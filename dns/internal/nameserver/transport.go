package nameserver

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/trevatk/tbd/lib/protocol"

	pba "github.com/trevatk/tbd/lib/protocol/dns/authoritative/v1"
	pbk "github.com/trevatk/tbd/lib/protocol/dns/kademlia/v1"
	pbr "github.com/trevatk/tbd/lib/protocol/dns/resolver/v1"
)

//go:generate mockgen -destination mock_dht_test.go -package nameserver . dht
type dht interface {
	findClosestNodes(nodeID) []*node
	addNode(context.Context, *node) error
	findNode(context.Context, nodeID) ([]*node, error)
	findValue(context.Context, nodeID) (*record, []*node, error)
	getSelf() *node

	getValue(string) (*record, error)
	setValue(string, *record) error
}

type grpcTransport struct {
	pbk.UnimplementedKademliaServiceServer
	pba.UnimplementedAuthoritativeServiceServer
	pbr.UnimplementedDNSResolverServiceServer

	logger *slog.Logger
	dht    dht
}

// interface compliance
var _ pbk.KademliaServiceServer = (*grpcTransport)(nil)
var _ pba.AuthoritativeServiceServer = (*grpcTransport)(nil)
var _ pbr.DNSResolverServiceServer = (*grpcTransport)(nil)

// NewTransport return new authoritative transport implementation
func NewTransport(logger *slog.Logger, dht dht) []protocol.Transport {
	tr := &grpcTransport{
		logger: logger,
		dht:    dht,
	}
	trs := []protocol.Transport{
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
	err := protocol.Validate(in)
	if err != nil {
		return nil, protocol.ErrInvalidArgument()
	}

	targetID, err := nodeIDFromStr(in.TargetNodeId)
	if err != nil {
		return nil, protocol.ErrInvalidArgument()
	}

	ns, err := t.dht.findNode(ctx, targetID)
	if err != nil {
		t.logger.ErrorContext(ctx, "find_node", slog.String("error", err.Error()))
		return nil, protocol.ErrInternal()
	}

	return newFindNodeResponse(ns, t.dht.getSelf(), in.RequestId), nil
}

// FindValue
func (t *grpcTransport) FindValue(ctx context.Context, in *pbk.FindValueRequest) (*pbk.FindValueResponse, error) {
	err := protocol.Validate(in)
	if err != nil {
		return nil, protocol.ErrInvalidArgument()
	}

	targetID, err := nodeIDFromStr(in.Key)
	if err != nil {
		t.logger.ErrorContext(ctx, "node id from string", slog.String("error", err.Error()))
	}

	record, closestNodes, err := t.dht.findValue(ctx, targetID)
	if err != nil {
		t.logger.ErrorContext(ctx, "find_value", slog.String("error", err.Error()))
		return nil, protocol.ErrInternal()
	}

	if record != nil {
		return newFindValueResponseWithRecord(t.dht.getSelf(), record, in.RequestId), nil
	}

	return newFindValueResponseWithClosestNodes(t.dht.getSelf(), closestNodes, in.RequestId), nil
}

// Ping
func (t *grpcTransport) Ping(ctx context.Context, in *pbk.PingRequest) (*pbk.PingResponse, error) {
	err := protocol.Validate(in)
	if err != nil {
		return nil, protocol.ErrInvalidArgument()
	}

	n, err := senderToNode(in.Sender)
	if err != nil {
		return nil, protocol.ErrInvalidArgument()
	}

	err = t.dht.addNode(ctx, &n)
	if err != nil {
		t.logger.ErrorContext(ctx, "add node", slog.String("error", err.Error()))
		return nil, protocol.ErrInternal()
	}

	return newPingResponse(t.dht.getSelf(), in.RequestId), nil
}

// Store
func (t *grpcTransport) Store(ctx context.Context, in *pbk.StoreRequest) (*pbk.StoreResponse, error) {
	err := protocol.Validate(in)
	if err != nil {
		return nil, protocol.ErrInvalidArgument()
	}

	t.logger.DebugContext(ctx, "store_value", slog.Any("request", in))

	if err := t.dht.setValue(in.Record.Id, &record{
		domain:     in.Record.Domain,
		recordType: in.Record.RecordType.String(),
		value:      in.Record.Value,
		ttl:        in.Record.Ttl,
	}); err != nil {
		t.logger.ErrorContext(ctx, "kv set value", slog.String("error", err.Error()))
		return nil, protocol.ErrInternal()
	}

	return newStoreResponse(t.dht.getSelf(), in.RequestId), nil
}

// Resolve
func (t *grpcTransport) Resolve(ctx context.Context, in *pbr.ResolveRequest) (*pbr.ResolveResponse, error) {
	err := protocol.Validate(in)
	if err != nil {
		return nil, protocol.ErrInvalidArgument()
	}

	t.logger.DebugContext(ctx, "resolve", slog.Any("request", in))

	return newResolveResponse(), nil
}

// CreateRecord
func (t *grpcTransport) CreateRecord(context.Context, *pba.CreateRecordRequest) (*pba.CreateRecordResponse, error) {
	return nil, nil
}

// CreateZone
func (t *grpcTransport) CreateZone(context.Context, *pba.CreateZoneRequest) (*pba.CreateZoneResponse, error) {
	return nil, nil
}

func newFindNodeResponse(ns []*node, sender *node, requestID string) *pbk.FindNodeResponse {
	closestNodes := make([]*pbk.Node, 0, len(ns))
	for _, n := range ns {
		closestNodes = append(closestNodes, nodeToSender(n))
	}
	return &pbk.FindNodeResponse{
		Sender:       nodeToSender(sender),
		RequestId:    requestID,
		ClosestNodes: closestNodes,
	}
}
func newFindValueResponseWithRecord(n *node, r *record, requestID string) *pbk.FindValueResponse {
	return &pbk.FindValueResponse{
		Sender: nodeToSender(n),
		Result: &pbk.FindValueResponse_Record{
			Record: &pbk.Record{
				Domain:     r.domain,
				RecordType: recordTypeToPb(r.recordType),
				Value:      r.value,
				Ttl:        r.ttl,
			},
		},
		RequestId: requestID,
	}
}

func newFindValueResponseWithClosestNodes(n *node, ns []*node, requestID string) *pbk.FindValueResponse {
	closestNodes := make([]*pbk.Node, 0, len(ns))
	for _, n := range ns {
		closestNodes = append(closestNodes, &pbk.Node{
			NodeId:     n.id.toString(),
			IpOrDomain: n.ipOrHost,
			Port:       n.port,
			LastSeen:   timestamppb.New(n.lastSeen),
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

func newPingResponse(n *node, requestID string) *pbk.PingResponse {
	return &pbk.PingResponse{
		Sender:    nodeToSender(n),
		RequestId: requestID,
	}
}

func newStoreResponse(n *node, requestID string) *pbk.StoreResponse {
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

func nodeToSender(n *node) *pbk.Node {
	return &pbk.Node{
		NodeId:     n.id.toString(),
		IpOrDomain: n.ipOrHost,
		Port:       n.port,
		LastSeen:   timestamppb.New(n.lastSeen),
	}
}

func senderToNode(n *pbk.Node) (node, error) {
	nodeID, err := nodeIDFromStr(n.NodeId)
	if err != nil {
		return node{}, fmt.Errorf("node is from string: %w", err)
	}
	return node{
		id:       nodeID,
		ipOrHost: n.IpOrDomain,
		port:     n.Port,
		lastSeen: n.LastSeen.AsTime(),
	}, nil
}

func pbToRecordType(rt pbk.Record_RECORDTYPE) string {
	switch rt {
	case pbk.Record_RECORDTYPE_A:
		return "A"
	case pbk.Record_RECORDTYPE_CNAME:
		return "CNAME"
	case pbk.Record_RECORDTYPE_DID:
		return "DID"
	default:
		return "unspecified"
	}
}

func recordTypeToPb(s string) pbk.Record_RECORDTYPE {
	switch strings.ToLower(s) {
	case "a":
		return pbk.Record_RECORDTYPE_A
	case "cname":
		return pbk.Record_RECORDTYPE_CNAME
	case "did":
		return pbk.Record_RECORDTYPE_DID
	default:
		return pbk.Record_RECORDTYPE_UNSPECIFIED
	}
}
