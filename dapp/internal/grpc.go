package internal

import (
	"context"
	"log/slog"

	"github.com/trevatk/tbd/lib/protocol"
	pb "github.com/trevatk/tbd/lib/protocol/chat/v1"
)

type thread struct {
	id   string
	name string
}

type service interface {
	createThread(context.Context, string) (thread, error)
}

type transport struct {
	pb.UnimplementedChatServiceServer

	logger *slog.Logger
	svc    service
}

// NewTransport
func NewTransport(logger *slog.Logger, svc service) protocol.Transport {
	tr := &transport{
		logger: logger,
		svc:    svc,
	}
	return protocol.Transport{
		ServiceDesc: &pb.ChatService_ServiceDesc,
		Service:     tr,
	}
}

// CreateThread
func (t *transport) CreateThread(ctx context.Context, in *pb.CreateThreadRequest) (*pb.CreateThreadResponse, error) {
	if err := protocol.Validate(in); err != nil {
		return nil, protocol.ErrInvalidArgument()
	}

	thread, err := t.svc.createThread(ctx, in.DisplayName)
	if err != nil {
		t.logger.ErrorContext(ctx, "create thread", slog.String("error", err.Error()))
		return nil, protocol.ErrInternal()
	}

	return newCreateThreadResponse(thread), nil
}

func newCreateThreadResponse(t thread) *pb.CreateThreadResponse {
	return &pb.CreateThreadResponse{
		Id:          t.id,
		DisplayName: t.name,
	}
}
