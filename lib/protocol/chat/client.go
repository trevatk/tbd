package chat

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"google.golang.org/grpc"

	"github.com/trevatk/tbd/lib/protocol"
	v1 "github.com/trevatk/tbd/lib/protocol/chat/v1"
	"github.com/trevatk/tbd/lib/protocol/resolver"
)

const (
	scheme = "did"
)

type (

	// Thread
	Thread struct {
		ID        string
		Name      string
		Members   []string
		CreatedAt time.Time
		UpdatedAt *time.Time
	}

	MessageSend struct {
		ThreadID string
		Contents []byte
	}

	Message struct {
		ThreadID string
		Sender   string
		Contents []byte
		SentAt   time.Time
	}

	Event struct {
		ThreadID string
	}

	// Client
	Client interface {
		// CreateThread
		CreateThread(context.Context, string, []string) (string, error)
		// ListThreads
		ListThreads(context.Context, string) ([]Thread, error)
		// SendMessage
		SendMessage(context.Context, MessageSend) (Message, error)
		// ListMessages
		ListMessages(context.Context, string) ([]Message, error)
		// Subscribe
		Subscribe(context.Context) (chan Event, error)
		// Close
		Close() error
	}

	clientV1 struct {
		conn     *grpc.ClientConn
		userAddr string
	}
)

// interface compliance
var _ Client = (*clientV1)(nil)

// NewClient
func NewClient(target, userAddr string) (Client, error) {
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
		conn:     conn,
		userAddr: userAddr,
	}, nil
}

// CreateThread
func (c *clientV1) CreateThread(ctx context.Context, name string, members []string) (string, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*250)
	defer cancel()

	request := &v1.CreateThreadRequest{
		DisplayName: name,
		Members:     members,
	}

	resp, err := v1.NewChatServiceClient(c.conn).CreateThread(timeout, request)
	if err != nil {
		return "", fmt.Errorf("failed to execute gRPC create thread: %w", err)
	}

	if protocol.Validate(resp); err != nil {
		return "", fmt.Errorf("invalid response message: %w", err)
	}

	return resp.Thread.Id, nil
}

// ListThreads
func (c *clientV1) ListThreads(context.Context, string) ([]Thread, error)

// SendMessage
func (c *clientV1) SendMessage(ctx context.Context, send MessageSend) (Message, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*250)
	defer cancel()

	request := &v1.SendMessageRequest{
		NewMessage: &v1.NewMessage{
			ThreadId: send.ThreadID,
			Contents: send.Contents,
			Sender:   c.userAddr,
		},
	}

	resp, err := v1.NewChatServiceClient(c.conn).SendMessage(timeout, request)
	if err != nil {
		return Message{}, fmt.Errorf("failed to execute gRPC send message: %w", err)
	}

	if err := protocol.Validate(resp); err != nil {
		return Message{}, fmt.Errorf("invalid response message: %w", err)
	}

	return Message{
		ThreadID: resp.Msg.Id,
		Sender:   resp.Msg.Sender,
		Contents: resp.Msg.Contents,
		SentAt:   resp.Msg.SentAt.AsTime(),
	}, nil
}

// ListMessages
func (c *clientV1) ListMessages(context.Context, string) ([]Message, error)

// Subscribe
func (c *clientV1) Subscribe(ctx context.Context) (chan Event, error) {
	request := &v1.SubscribeEventsRequest{}
	stream, err := v1.NewChatServiceClient(c.conn).SubscribeEvents(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to execute gRPC subscribe events: %w", err)
	}
	return listenForEvents(ctx, stream)
}

// Close
func (c *clientV1) Close() error {
	return c.conn.Close()
}

func listenForEvents(
	ctx context.Context,
	stream grpc.ServerStreamingClient[v1.SubscribeEventsResponse],
) (chan Event, error) {
	ch := make(chan Event)
	errCh := make(chan error, 1)

	go func() {
		defer close(ch)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				stream.CloseSend()
				return
			case <-stream.Context().Done():
				return
			default:
				// fallthrough
				msg, err := stream.Recv()
				if err != nil && errors.Is(err, io.EOF) {
					// stream closed with success message
					return
				} else if err != nil {
					errCh <- fmt.Errorf("failed to receive message: %w", err)
					return
				}

				ch <- Event{
					ThreadID: msg.Event.ThreadId,
				}
			}
		}
	}()

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	return ch, nil
}
