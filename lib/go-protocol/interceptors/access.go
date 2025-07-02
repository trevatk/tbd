package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Checker interceptor access control
type Checker interface {
	HasResourceAccess(userHash, resourceHash string) bool
}

// AccessControl access control wrapper
type AccessControl struct {
	c Checker
}

// NewAccessControl return new access control wrapper with gRPC unary interceptor
func NewAccessControl(checker Checker) *AccessControl {
	return &AccessControl{
		c: checker,
	}
}

// EnsureResourceAccess gRPC unary interceptor func
func (ac *AccessControl) EnsureResourceAccess() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// interceptor func check if user has access to requested resource
		userHash := ctx.Value(User).(string)

		granted := ac.c.HasResourceAccess(userHash, "")
		if !granted {
			return nil, status.Error(codes.PermissionDenied, codes.PermissionDenied.String())
		}
		return handler(ctx, req)
	}
}
