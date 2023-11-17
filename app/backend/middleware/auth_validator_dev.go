//go:build dev

package middleware

import (
	"context"
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
	logger *zap.Logger,
	r *http.Request,
) (*validator.ValidatedClaims, error) {
	token, err := parseToken(r)
	if err != nil {
		logger.Error("auth failed in dev mode: failed to parse token", zap.Error(err))
		return nil, nil
	}
	claims, err := v.ValidateToken(ctx, token)
	if err != nil {
		logger.Error("auth failed in dev mode: failed to validate token", zap.Error(err))
		return nil, nil
	}
	return claims.(*validator.ValidatedClaims), nil
}
