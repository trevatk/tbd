// Package internal dapp application and controller layers
package internal

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/google/uuid"

	"github.com/structx/tbd/lib/protocol"
	pb "github.com/structx/tbd/lib/protocol/chat/v1"
)

type thread struct {
	id        string
	name      string
	createdAt time.Time
}

type transport struct {
	pb.UnimplementedChatServiceServer

	mu sync.RWMutex

	logger  *slog.Logger
	threads map[string]thread
}

// NewTransport return dapp specific gateway transport implementation
func NewTransport(logger *slog.Logger) (*grpc.ServiceDesc, pb.ChatServiceServer) {
	return &pb.ChatService_ServiceDesc, &transport{
		logger:  logger,
		threads: make(map[string]thread),
		mu:      sync.RWMutex{},
	}
}

// CreateThread
func (t *transport) CreateThread(ctx context.Context, in *pb.CreateThreadRequest) (*pb.CreateThreadResponse, error) {
	t.logger.DebugContext(ctx, "CreateThread", "request", in)

	if err := protocol.Validate(in); err != nil {
		return nil, protocol.ErrInvalidArgument()
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	var th thread

	for {
		// create in for loop to check if thread id already exists
		// continue loop to generate new id
		id, err := uuid.NewV7()
		if err != nil {
			t.logger.ErrorContext(ctx, "failed to create thread id", "error", err)
			return nil, protocol.ErrInternal()
		}

		if _, ok := t.threads[id.String()]; ok {
			// id exists
			// continue to regenerate
			continue
		}

		th = thread{
			id:        id.String(),
			name:      in.DisplayName,
			createdAt: time.Now(),
		}
		break
	}

	// add thread to inmemory store
	t.threads[th.id] = th

	return newCreateResponse(th), nil
}

func newCreateResponse(t thread) *pb.CreateThreadResponse {
	return &pb.CreateThreadResponse{
		Id:          t.id,
		DisplayName: t.name,
	}
}
