package interceptors

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	// rename emptyAuthorizations
	minAuthorization = 0
)

// CustomClaims token claims
type CustomClaims struct {
	jwt.RegisteredClaims
}

type contextKey string

const (
	// User context key
	User contextKey = "user"
	// Anon anonymous user
	Anon = "anonymous"
)

// Verifier auth token validator
type Verifier interface {
	Verify(string) (map[string]interface{}, error)
	SignWithClaims(map[string]interface{}) (string, error)
}

// Auth gRPC interceptor wrapper to include service verifier
type Auth struct {
	v Verifier
}

// NewAuth return new auth verifier
func NewAuth(verifier Verifier) *Auth {
	return &Auth{
		v: verifier,
	}
}

// ValidToken gRPC unary auth interceptor
func (a *Auth) ValidToken() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// auth interceptor
		// check for bearer token to validate
		// if token is not included anonymous user is assumed
		md, _ := metadata.FromIncomingContext(ctx)

		authorization := md.Get("Authorization")
		if len(authorization) > minAuthorization {
			userHash, ok := a.valid(authorization)
			if !ok {
				return nil, status.Error(codes.Unauthenticated, codes.Unauthenticated.String())
			}
			ctx = context.WithValue(ctx, User, userHash)
		} else {
			ctx = context.WithValue(ctx, User, Anon)
		}

		return handler(ctx, req)
	}
}

func (a *Auth) valid(authorization []string) (string, bool) {
	if len(authorization) == minAuthorization {
		return "", false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	claims, err := a.v.Verify(token)
	if err != nil {
		return "", false
	}
	return claims["user_id"].(string), true
}
