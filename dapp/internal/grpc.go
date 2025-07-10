package internal

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/trevatk/tbd/lib/protocol"
	pb "github.com/trevatk/tbd/lib/protocol/chat/v1"
	wellknown "github.com/trevatk/tbd/lib/protocol/wellknown/v1"
)

type (
	threadCreate struct {
		name    string
		members []string
	}

	thread struct {
		id        string
		name      string
		members   []string
		createdAt time.Time
		UpdatedAt *time.Time
	}

	service interface {
		createThread(context.Context, threadCreate) (thread, error)
	}

	transport struct {
		pb.UnimplementedChatServiceServer
		wellknown.UnimplementedWellKnownServiceServer

		logger *slog.Logger
		svc    service
	}
)

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

	create := threadCreate{
		name:    in.DisplayName,
		members: in.Members,
	}

	thread, err := t.svc.createThread(ctx, create)
	if err != nil {
		t.logger.ErrorContext(ctx, "create thread", slog.String("error", err.Error()))
		return nil, protocol.ErrInternal()
	}

	return newCreateThreadResponse(thread), nil
}

// ListMessages
func (t *transport) ListMessages(context.Context, *pb.ListMessagesRequest) (*pb.ListMessagesResponse, error) {
	return nil, nil
}

// ListThreads
func (t *transport) ListThreads(context.Context, *pb.ListThreadsRequest) (*pb.ListThreadsResponse, error) {
	return nil, nil
}

// SendMessage
func (t *transport) SendMessage(context.Context, *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	return nil, nil
}

// SubscribeEvents
func (t *transport) SubscribeEvents(
	*pb.SubscribeEventsRequest,
	grpc.ServerStreamingServer[pb.SubscribeEventsResponse],
) error {
	return nil
}

func newCreateThreadResponse(t thread) *pb.CreateThreadResponse {
	var updatedAt *timestamppb.Timestamp
	if t.UpdatedAt != nil {
		updatedAt = timestamppb.New(*t.UpdatedAt)
	}
	return &pb.CreateThreadResponse{
		Thread: &pb.Thread{
			Id:          t.id,
			DisplayName: t.name,
			Members:     t.members,
			CreatedAt:   timestamppb.New(t.createdAt),
			UpdatedAt:   updatedAt,
		},
	}
}
