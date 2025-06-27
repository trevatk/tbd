package interceptors

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// LoggerUnary
func LoggerUnary(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		logger.Info("request handler", zap.String("method", info.FullMethod))
		return handler(ctx, req)
	}
}
