package interceptors

import (
	"context"

	"google.golang.org/grpc"
)

// UnaryInterceptor
type UnaryInterceptor func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
