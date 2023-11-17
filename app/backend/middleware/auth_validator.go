//go:build !dev

package middleware

import (
	"context"
	"fmt"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"go.uber.org/zap"
	"net/http"
)

// getRawClaims extracts the raw claims from the request, however, it is meant
// only to be used in dev mode so that unauthenticated requests can still succeed
// such as requests to load the graphql schema.
//
// Any validation or parser errors are logged rather than returned.
func getRawClaims(
	ctx context.Context,
	v *validator.Validator,
	_ *zap.Logger,
	r *http.Request,
) (*validator.ValidatedClaims, error) {
	token, err := parseToken(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	raw, err := v.ValidateToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}
	if claims, ok := raw.(*validator.ValidatedClaims); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("expected %T, got %T", &validator.ValidatedClaims{}, raw)
}
