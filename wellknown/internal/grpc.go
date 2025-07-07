package internal

import (
	"context"
	"log/slog"

	"github.com/trevatk/tbd/lib/protocol"
	pb "github.com/trevatk/tbd/lib/protocol/wellknown/v1"
)

type transport struct {
	pb.UnimplementedWellKnownServiceServer

	logger *slog.Logger
}

// interface compliance
var _ pb.WellKnownServiceServer = (*transport)(nil)

// used for testing
func newGrpcTransport(logger *slog.Logger) pb.WellKnownServiceServer {
	return &transport{
		logger: logger,
	}
}

// NewTransport return new gRPC transport
func NewTransport(logger *slog.Logger) protocol.Transport {
	tr := &transport{
		logger: logger,
	}
	return protocol.Transport{
		ServiceDesc: &pb.WellKnownService_ServiceDesc,
		Service:     tr,
	}
}

// GetDIDConfiguration
func (t *transport) GetDIDConfiguration(
	context.Context,
	*pb.GetDIDConfigurationRequest,
) (*pb.GetDIDConfigurationResponse, error) {
	return newDidConfigurationResponse(), nil
}

func newDidConfigurationResponse() *pb.GetDIDConfigurationResponse {
	return &pb.GetDIDConfigurationResponse{}
}
