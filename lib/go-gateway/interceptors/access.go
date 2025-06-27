package interceptors

import (
	"context"

	"google.golang.org/grpc"
)

// Checker
type Checker interface {
	HasResourceAccess(userHash, resourceHash string) bool
}

// AccessControl
type AccessControl struct {
	c Checker
}

// NewAccessControl
func NewAccessControl(checker Checker) *AccessControl {
	return &AccessControl{
		c: checker,
	}
}

// EnsureResourceAccess
func (ac *AccessControl) EnsureResourceAccess() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		userHash := ctx.Value(User).(string)
		granted := ac.c.HasResourceAccess(userHash, "")
		if !granted {
		}
		return handler(ctx, req)
	}
}
