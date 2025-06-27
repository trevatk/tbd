package gateway

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	system = ""

	timeoutSeconds = 15

	sleep = 3
)

type server struct {
	logger *slog.Logger

	host, port string
	transports []Transport
}

// Option server option pattern
type Option func(*server)

// Transport gRPC service description and implementation
type Transport struct {
	ServiceDesc interface{}
	Service     any
}

// New returns new server with options
func New(opts ...Option) *server {
	s := &server{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// WithHost server host
func WithHost(host string) Option {
	return func(s *server) {
		s.host = host
	}
}

// WithPort server port
func WithPort(port string) Option {
	return func(s *server) {
		s.port = port
	}
}

// WithTransports server transports to be served
func WithTransports(transports []Transport) Option {
	return func(s *server) {
		s.transports = transports
	}
}

// WithLogger server logger
func WithLogger(logger *slog.Logger) Option {
	return func(s *server) {
		s.logger = logger
	}
}

func (s *server) StartAndStop(ctx context.Context) error {
	lis, err := net.Listen("tcp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		return fmt.Errorf("network listener %w", err)
	}
	defer func() { _ = lis.Close() }()

	ss := grpc.NewServer()
	for _, t := range s.transports {
		desc, ok := t.ServiceDesc.(*grpc.ServiceDesc)
		if !ok {
			return errors.New("invalid service description provided")
		}
		ss.RegisterService(desc, t.Service)
	}

	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(ss, healthcheck)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(ctx, time.Second*timeoutSeconds)
		defer cancel()
		timer := time.AfterFunc(time.Second*timeoutSeconds, func() {
			s.logger.InfoContext(shutdownCtx, "gRPC server force shutdown")
			ss.Stop()
		})
		defer timer.Stop()
		s.logger.InfoContext(shutdownCtx, "gRPC server graceful shutdown")
		ss.GracefulStop()
	}()

	go func() {
		// asynchronously inspect dependencies and toggle serving status as needed
		next := healthpb.HealthCheckResponse_SERVING

		for {
			healthcheck.SetServingStatus(system, next)

			if next == healthpb.HealthCheckResponse_SERVING {
				next = healthpb.HealthCheckResponse_NOT_SERVING
			} else {
				next = healthpb.HealthCheckResponse_SERVING
			}

			time.Sleep(time.Second * sleep)
		}
	}()

	s.logger.InfoContext(ctx, "start gRPC server", slog.String("listener_addr", lis.Addr().String()))
	if err := ss.Serve(lis); err != nil {
		return fmt.Errorf("start gRPC server: %w", err)
	}

	return nil
}
