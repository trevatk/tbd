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

// CustomClaims
type CustomClaims struct {
	jwt.RegisteredClaims
}

type contextKey string

const (
	// User
	User contextKey = "user"

	Anon = "anonymous"
)

// Verifier
type Verifier interface {
	Verify(string) (map[string]interface{}, error)
	SignWithClaims(map[string]interface{}) (string, error)
}

// Auth
type Auth struct {
	v Verifier
}

// NewAuth
func NewAuth(verifier Verifier) *Auth {
	return &Auth{
		v: verifier,
	}
}

// ValidToken
func (a *Auth) ValidToken() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		md, _ := metadata.FromIncomingContext(ctx)

		authorization := md.Get("Authorization")
		if len(authorization) > 0 {
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
	if len(authorization) < 1 {
		return "", false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	claims, err := a.v.Verify(token)
	if err != nil {
		return "", false
	}
	return claims["user_id"].(string), true
}
