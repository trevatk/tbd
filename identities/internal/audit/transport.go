package audit

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/structx/tbd/lib/gateway"
	pb "github.com/structx/tbd/lib/protocol/audit/v1"
)

type transport struct {
	pb.UnimplementedAuditServiceServer

	logger *slog.Logger
	svc *serviceImpl
}

// NewTransport
func NewTransport(logger *slog.Logger, svc *serviceImpl) (*grpc.ServiceDesc, pb.AuditServiceServer) {
	return &pb.AuditService_ServiceDesc,
	&transport{
		logger: logger,
		svc: svc,
	}
}

func (t *transport) Decision(ctx context.Context, in *pb.CreateDecisionRequest) (*emptypb.Empty, error) {
	return gateway.NewEmptyResponse(), nil
}

func (t *transport) ListDecisions(ctx context.Context, in *pb.ListDecisionsRequest) (*pb.ListDecisionsResponse, error) {
	txs, err := t.svc.listTxs(in.Limit, in.Offset)
	if err != nil {
		t.logger.Error("failed to list transactions", zap.Error(err))
		return nil, gateway.ErrInternal()
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
