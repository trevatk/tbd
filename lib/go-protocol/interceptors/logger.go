package interceptors

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// LoggerUnary gRPC unary logger interceptor
func LoggerUnary(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		attrs := []slog.Attr{}
		if p, ok := peer.FromContext(ctx); ok {
			attrs = append(attrs, slog.Any("peer", p))
		}

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			attrs = append(attrs, slog.Any("metadata", md))
		}

		attrs = append(attrs, slog.String("method", info.FullMethod), slog.Any("request", req))
		logger.LogAttrs(ctx, slog.LevelInfo, "request", attrs...)

		return handler(ctx, req)
	}
}
