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
)

const (
	system = ""
)

type Server struct {
	logger *slog.Logger

	host, port string
	transports []Transport
}

type Option func(*Server)

// Transport
type Transport struct {
	ServiceDesc interface{}
	Service     any
}

// New
func New(opts ...Option) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// WithHost
func WithHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

// WithPort
func WithPort(port string) Option {
	return func(s *Server) {
		s.port = port
	}
}

// WithTransports
func WithTransports(transports []Transport) Option {
	return func(s *Server) {
		s.transports = transports
	}
}

// WithLogger
func WithLogger(logger *slog.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

func (s *Server) StartAndStop(ctx context.Context) error {
	lis, err := net.Listen("tcp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		return fmt.Errorf("network listener %w", err)
	}
	defer lis.Close()

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
		shutdownCtx, cancel := context.WithTimeout(ctx, time.Second*15)
		defer cancel()
		timer := time.AfterFunc(time.Second*15, func() {
			s.logger.InfoContext(shutdownCtx, "gRPC server force shutdown")
			ss.Stop()
		})
		defer timer.Stop()
		s.logger.InfoContext(shutdownCtx, "gRPC server graceful shutdown")
		ss.GracefulStop()
	}()

	s.logger.InfoContext(ctx, "start gRPC server", slog.String("listener_addr", lis.Addr().String()))
	if err := ss.Serve(lis); err != nil {
		return fmt.Errorf("start gRPC server: %w", err)
	}

	return nil
}
