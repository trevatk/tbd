package internal

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/trevatk/tbd/lib/protocol"
	"github.com/trevatk/tbd/lib/protocol/did/v1"
	pb "github.com/trevatk/tbd/lib/protocol/wellknown/v1"
)

type transport struct {
	pb.UnimplementedWellKnownServiceServer

	mu sync.Mutex

	logger *slog.Logger

	filePath string
}

// interface compliance
var _ pb.WellKnownServiceServer = (*transport)(nil)

// used for testing
func newGrpcTransport(logger *slog.Logger, filePath string) pb.WellKnownServiceServer {
	return &transport{
		logger:   logger,
		mu:       sync.Mutex{},
		filePath: filePath,
	}
}

// NewTransport return new gRPC transport
func NewTransport(logger *slog.Logger, filePath string) protocol.Transport {
	tr := &transport{
		logger:   logger,
		mu:       sync.Mutex{},
		filePath: filePath,
	}
	return protocol.Transport{
		ServiceDesc: &pb.WellKnownService_ServiceDesc,
		Service:     tr,
	}
}

// GetDIDConfiguration
func (t *transport) GetDIDConfiguration(
	ctx context.Context,
	_ *pb.GetDIDConfigurationRequest,
) (*pb.GetDIDConfigurationResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	fbytes, err := os.ReadFile(filepath.Clean(t.filePath))
	if err != nil {
		t.logger.ErrorContext(ctx, "failed to read file", slog.String("error", err.Error()))
		return nil, protocol.ErrInternal()
	}

	var doc did.Document
	err = json.Unmarshal(fbytes, &doc)
	if err != nil {
		t.logger.ErrorContext(ctx, "failed to unmarshal json", slog.String("error", err.Error()))
		return nil, protocol.ErrInternal()
	}

	return newDidConfigurationResponse(&doc), nil
}

func newDidConfigurationResponse(d *did.Document) *pb.GetDIDConfigurationResponse {
	return &pb.GetDIDConfigurationResponse{
		Doc: d,
	}
}
