package token

import (
	"github.com/clarkmcc/cloudcore/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const defaultTokenDuration = time.Minute

type Signer struct {
	secret []byte
}

func (s *Signer) NewToken() (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(defaultTokenDuration).Unix(),
	}).SignedString(s.secret)
}

func (s *Signer) ValidateToken(token string) error {
	_, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return s.secret, nil
	})
	return err
}

func NewSigner(config *config.ServerConfig) *Signer {
	return &Signer{
		secret: []byte(config.Auth.TokenSigningSecret),
	}
}
