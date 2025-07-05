package chat

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/trevatk/tbd/lib/protocol"
	v1 "github.com/trevatk/tbd/lib/protocol/chat/v1"
	"github.com/trevatk/tbd/lib/protocol/resolver"
)

const (
	scheme = ""
)

// Client
type Client interface {
	// CreateThread
	CreateThread(context.Context, string) (string, error)

	Close() error
}

type clientV1 struct {
	conn *grpc.ClientConn
}

// interface compliance
var _ Client = (*clientV1)(nil)

// NewClient
func NewClient(target string) (Client, error) {

	opts := []resolver.Option{
		resolver.WithScheme(scheme),
		resolver.WithServiceName(v1.ChatService_ServiceDesc.ServiceName),
	}
	builder := resolver.NewBuilder(opts...)

	conn, err := protocol.NewConnWithResolver(target, builder)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection with resolver: %w", err)
	}

	return &clientV1{
		conn: conn,
	}, nil
}

// CreateThread
func (c *clientV1) CreateThread(ctx context.Context, name string) (string, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*250)
	defer cancel()

	request := &v1.CreateThreadRequest{
		DisplayName: name,
	}

	resp, err := v1.NewChatServiceClient(c.conn).CreateThread(timeout, request)
	if err != nil {
		return "", fmt.Errorf("failed to execute gRPC create thread: %w", err)
	}

	return resp.Id, nil
}

// Close
func (c *clientV1) Close() error {
	return c.conn.Close()
}
