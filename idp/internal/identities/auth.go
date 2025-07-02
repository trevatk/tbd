package identities

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/structx/tbd/lib/protocol/interceptors"
)

type auth struct {
	hmacSecret string
}

// interface compliance
var _ interceptors.Verifier = (*auth)(nil)

// NewAuth return new idenities auth implementation of verifier
func NewAuth(hmacSecret string) interceptors.Verifier {
	return &auth{
		hmacSecret: hmacSecret,
	}
}

// Verify
func (a *auth) Verify(token string) (map[string]interface{}, error) {
	var cc interceptors.CustomClaims
	if t, err := jwt.ParseWithClaims(token, &cc, func(t *jwt.Token) (interface{}, error) {
		return a.hmacSecret, nil
	}); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	} else if claims, ok := t.Claims.(*interceptors.CustomClaims); ok {
		md := make(map[string]interface{})
		md["user_id"] = claims.Issuer
		return md, nil
	}

	return nil, errors.New("unknown claims type")
}

// SignWithClaims
func (a *auth) SignWithClaims(claims map[string]interface{}) (string, error) {
	md := &interceptors.CustomClaims{}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, md)
	s, err := t.SignedString([]byte(a.hmacSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return s, nil
}
