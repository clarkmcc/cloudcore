package token

import (
	"context"
	"fmt"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/config"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
	"time"
)

const defaultTokenDuration = 5 * time.Minute

type Signer struct {
	secret []byte
}

func (s *Signer) NewToken(agentId string) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(defaultTokenDuration).Unix(),
		"id":  agentId,
	}).SignedString(s.secret)
}

func (s *Signer) ValidateToken(token string) (string, error) {
	tk, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return s.secret, nil
	})
	if err != nil {
		return "", err
	}
	c, ok := tk.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}
	agentId, ok := c["id"]
	if !ok {
		return "", fmt.Errorf("token missing id")
	}
	if _, ok = agentId.(string); !ok {
		return "", fmt.Errorf("token id is not a string")
	}
	return agentId.(string), nil
}

func (s *Signer) ValidateRequest(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing metadata")
	}
	tokens := md.Get("token")
	if len(tokens) == 0 || len(tokens[0]) == 0 {
		return "", fmt.Errorf("missing token")
	}
	return s.ValidateToken(tokens[0])
}

func NewSigner(config *config.Config) *Signer {
	return &Signer{
		secret: []byte(config.Auth.TokenSigningSecret),
	}
}
