package resolver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/protobuf/proto"

	"github.com/structx/tbd/lib/protocol"
	pb "github.com/structx/tbd/lib/protocol/dns/resolver/v1"
)

const (
	numZero = 0
	errAttr = "error"
)

type transport struct {
	pb.UnimplementedDNSResolverServiceServer

	logger *slog.Logger

	cache Cache

	ns []string // nameserver addrs
}

// interface compliance
var _ pb.DNSResolverServiceServer = (*transport)(nil)

// NewTransport return new resolver implementation of gateway transport
func NewTransport(logger *slog.Logger, nameservers []string, cache Cache) protocol.Transport {
	tr := &transport{
		logger: logger,
		cache:  cache,
		ns:     nameservers,
	}
	return protocol.Transport{
		ServiceDesc: &pb.DNSResolverService_ServiceDesc,
		Service:     tr,
	}
}

// Resolve
func (t *transport) Resolve(ctx context.Context, in *pb.ResolveRequest) (*pb.ResolveResponse, error) {
	if err := protocol.Validate(in); err != nil {
		return nil, protocol.ErrInvalidArgument()
	}

	// attempt to resolve using local cache
	cacheKey := buildCacheKey(in.Question.Domain, in.Question.RecordType.String())
	value, err := t.cache.Get(cacheKey)
	if err == nil {
		// cache hit
		var rr pb.ResolveResponse
		err = proto.Unmarshal(value, &rr)
		if err != nil {
			t.logger.ErrorContext(ctx, "proto.Unmarshal", slog.String(errAttr, err.Error()))
		}
		return &rr, nil
	} else if !errors.Is(err, ErrKeyNotFound) {
		// uncaught cache error
		t.logger.ErrorContext(ctx, "failed to get cache value", slog.String(errAttr, err.Error()))
	}

	var didJSON []byte
	if len(in.DidToResolve) > numZero {
		didJSON, err = resolveDID(in.DidToResolve)
		if err != nil {
			t.logger.ErrorContext(ctx, "failed to resolve did", slog.String(errAttr, err.Error()))
			return nil, protocol.ErrInternal()
		}
	}

	// dns record resolution

	return newResolveResponse(didJSON), nil
}

func newResolveResponse(didJSON []byte) *pb.ResolveResponse {
	return &pb.ResolveResponse{
		Status:                  pb.ResolveResponse_RESPONSE_STATUS_SUCCESS,
		ResolvedDidDocumentJson: string(didJSON),
	}
}

func buildCacheKey(s1, s2 string) string {
	return fmt.Sprintf("%s:%s", s1, s2)
}
