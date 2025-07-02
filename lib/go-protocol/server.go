package protocol

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
)

// Server protcol server lifecycle
type Server interface {
	StartAndStop(ctx context.Context) error
}

type gateway struct {
	logger *slog.Logger

	host, port string
	tranports  []Transport
}

// ServerOption server option pattern
type ServerOption func(*gateway)

// NewServer returns new server with options
func NewServer(opts ...ServerOption) *gateway {
	g := &gateway{}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// WithHost server host
func WithHost(host string) ServerOption {
	return func(g *gateway) {
		g.host = host
	}
}

// WithPort server port
func WithPort(port string) ServerOption {
	return func(g *gateway) {
		g.port = port
	}
}

// WithTransports server transports to be served
func WithTransports(transports []Transport) ServerOption {
	return func(g *gateway) {
		g.tranports = transports
	}
}

// WithLogger server logger
func WithLogger(logger *slog.Logger) ServerOption {
	return func(g *gateway) {
		g.logger = logger
	}
}

// StartAndStop starts and stops server
func (g *gateway) StartAndStop(ctx context.Context) error {
	lis, err := net.Listen("tcp", net.JoinHostPort(g.host, g.port))
	if err != nil {
		return fmt.Errorf("network listener %w", err)
	}
	defer func() { _ = lis.Close() }()

	s := grpc.NewServer()
	for _, tr := range g.tranports {
		s.RegisterService(tr.ServiceDesc, tr.Service)
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		timer := time.AfterFunc(time.Second, func() {
			g.logger.InfoContext(shutdownCtx, "gRPC server force shutdown")
			s.Stop()
		})
		defer timer.Stop()
		g.logger.InfoContext(shutdownCtx, "gRPC server graceful shutdown")
		s.GracefulStop()
	}()

	g.logger.InfoContext(ctx, "start gRPC server", slog.String("listener_addr", lis.Addr().String()))
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("start gRPC server: %w", err)
	}

	return nil
}
