package audit

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/structx/tbd/lib/protocol"
	pb "github.com/structx/tbd/lib/protocol/audit/v1"
)

type transport struct {
	pb.UnimplementedAuditServiceServer

	logger *slog.Logger
	svc    *serviceImpl
}

// NewTransport return idp audit gateway transport implementation
func NewTransport(logger *slog.Logger, svc *serviceImpl) (*grpc.ServiceDesc, pb.AuditServiceServer) {
	return &pb.AuditService_ServiceDesc,
		&transport{
			logger: logger,
			svc:    svc,
		}
}

func (t *transport) Decision(ctx context.Context, in *pb.CreateDecisionRequest) (*emptypb.Empty, error) {
	return protocol.NewEmptyResponse(), nil
}

func (t *transport) ListDecisions(ctx context.Context, in *pb.ListDecisionsRequest) (*pb.ListDecisionsResponse, error) {
	txs, err := t.svc.listTxs(in.Limit, in.Offset)
	if err != nil {
		t.logger.ErrorContext(ctx, "failed to list transactions", "error", err)
		return nil, protocol.ErrInternal()
	}

	return newListDecisionsResponse(txs), nil
}

func newListDecisionsResponse(txs []*tx) *pb.ListDecisionsResponse {
	return &pb.ListDecisionsResponse{
		Txs: []*pb.Tx{transformTx(txs[0])},
	}
}

func transformTx(tx *tx) *pb.Tx {
	return &pb.Tx{
		Hash: tx.Hash,
		From: tx.From,
		To:   tx.To,
		Sig:  tx.Sig,
	}
}
