package middleware

import (
	"context"
	"errors"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

func Authentication(config *config.Config, logger *zap.Logger) gin.HandlerFunc {
	v, err := newValidator(config)
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		token, err := parseToken(c.Request)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, errors.New("invalid auth header"))
			return
		}
		rawClaims, err := v.ValidateToken(c, token)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		// Append the validated claims to the context if we have them
		ctx := c.Request.Context()
		if claims, ok := rawClaims.(*validator.ValidatedClaims); ok {
			c.Request = c.Request.WithContext(withClaimsContext(ctx, claims))
		}
		c.Next()
	}
}

func parseToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}
	return parts[1], nil
}

func newValidator(config *config.Config) (*validator.Validator, error) {
	iss, err := url.Parse("https://" + config.Auth0.Domain + "/")
	if err != nil {
		return nil, err
	}
	provider := jwks.NewCachingProvider(iss, time.Minute)
	return validator.New(
		provider.KeyFunc,
		validator.RS256,
		iss.String(),
		[]string{config.Auth0.Audience},
	)
}

var claimsContextKey = reflect.TypeOf(&validator.ValidatedClaims{}).String()

func withClaimsContext(ctx context.Context, claims *validator.ValidatedClaims) context.Context {
	return context.WithValue(ctx, claimsContextKey, claims)
}

func ClaimsFromContext(ctx context.Context) *validator.ValidatedClaims {
	claims, ok := ctx.Value(claimsContextKey).(*validator.ValidatedClaims)
	if !ok {
		panic("claims not found in context")
	}
	return claims
}

func SubjectFromContext(ctx context.Context) string {
	return ClaimsFromContext(ctx).RegisteredClaims.Subject
}
