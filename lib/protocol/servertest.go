package protocol

import (
	"context"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/trevatk/tbd/lib/protocol/interceptors"
)

const (
	bufSize = 1024 * 1024
)

// TestServerOption test server option
type TestServerOption func(*TestServer)

// TestServer used for test purposes only
type TestServer struct {
	logger *slog.Logger

	lis        *bufconn.Listener
	transports []Transport

	gserver *grpc.Server
}

// NewTestServer return new test server with bufnet listener
func NewTestServer(opts ...TestServerOption) *TestServer {
	var (
		lis = bufconn.Listen(bufSize)
	)

	ts := &TestServer{
		lis:        lis,
		transports: make([]Transport, 0),
		logger:     nil,
		gserver:    nil,
	}

	for _, opt := range opts {
		opt(ts)
	}

	return ts
}

// WithTestTransports test server transport option
func WithTestTransports(transports []Transport) TestServerOption {
	return func(ts *TestServer) {
		ts.transports = transports
	}
}

// WithTestLogger test server logger option
func WithTestLogger(logger *slog.Logger) TestServerOption {
	return func(ts *TestServer) {
		ts.logger = logger
	}
}

// BufDialer gRPC bufnet dialer with context
func (ts *TestServer) BufDialer(ctx context.Context, _ string) (net.Conn, error) {
	return ts.lis.DialContext(ctx)
}

// Start server
func (ts *TestServer) Start(ctx context.Context) error {
	ts.logger.DebugContext(ctx, "start gRPC bufnet server")

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptors.LoggerUnary(ts.logger)),
	}

	ts.gserver = grpc.NewServer(opts...)

	for _, t := range ts.transports {
		ts.gserver.RegisterService(t.ServiceDesc, t.Service)
	}
	return ts.gserver.Serve(ts.lis)
}

// Stop stop
func (ts *TestServer) Stop(ctx context.Context) {
	timer := time.AfterFunc(time.Second, func() {
		ts.logger.DebugContext(ctx, "force stop gRPC bufnet server")
		ts.gserver.Stop()
	})
	defer timer.Stop()
	ts.logger.DebugContext(ctx, "graceful stop gRPC bufnet server")
	ts.gserver.GracefulStop()
}
