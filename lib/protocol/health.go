package protocol

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	system = ""
)

// RegisterHealthCheck register health check server to gRPC server
func RegisterHealthCheck(s *grpc.Server) *health.Server {
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(s, healthcheck)
	return healthcheck
}

// ServeHealthCheck start health check for loop
func ServeHealthCheck(ctx context.Context, healthcheck *health.Server, sleep int64) {
	// asynchronously inspect dependencies and toggle serving status as needed
	next := healthpb.HealthCheckResponse_SERVING

	for {
		select {
		case <-ctx.Done():
			healthcheck.Shutdown()
		default:
			// fallthrough
		}

		healthcheck.SetServingStatus(system, next)

		if next == healthpb.HealthCheckResponse_SERVING {
			next = healthpb.HealthCheckResponse_NOT_SERVING
		} else {
			next = healthpb.HealthCheckResponse_SERVING
		}

		time.Sleep(time.Second * time.Duration(sleep))
	}
}
