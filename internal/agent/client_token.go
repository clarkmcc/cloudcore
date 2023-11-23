// MIT License
//
// Copyright (c) 2024 Clark McCauley
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package agent

import (
	"context"
	"errors"
	"fmt"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"github.com/clarkmcc/cloudcore/internal/token"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type SystemMetadataProvider interface {
	GetSystemMetadata(ctx context.Context) (*rpc.SystemMetadata, error)
}

// tokenManager manages the lifecycle of the JWT-based authentication tokens
// that are used to authenticate the agent with the server.
type tokenManager struct {
	logger           *zap.Logger
	config           *Config
	db               Database
	metadataProvider SystemMetadataProvider

	client *Client
}

// getAuthTokenLocked returns the current auth token, or generates a new one if
// the current token is expired or expiring soon. Callers should be sure to hold
// the Client's lock before calling this method, since this method depends on the
// client and an active gRPC connection in order to acquire a new token.
func (m *tokenManager) getAuthTokenLocked(ctx context.Context) (string, error) {
	tk, err := m.db.AuthToken(ctx)
	if err != nil && !errors.Is(err, ErrAuthTokenNotFound) {
		return "", err
	}
	if tk != nil && !isExpired(tk.Expiration) && !isExpiringSoon(time.Now(), tk.Expiration, tk.Duration) {
		return tk.Token, nil
	}
	tk, err = m.newToken(ctx, tk)
	if err != nil {
		return "", err
	}
	err = m.db.SaveAuthToken(ctx, tk)
	if err != nil {
		return "", err
	}
	return tk.Token, nil
}

func (m *tokenManager) newToken(ctx context.Context, maybeAuthToken *AuthToken) (*AuthToken, error) {
	// Dynamically construct the authentication request based on the type
	// of flow we're performing.
	var req rpc.AuthenticateRequest
	switch {
	case maybeAuthToken != nil:
		m.logger.Debug("using existing token auth flow")
		req = rpc.AuthenticateRequest{
			Flow:  rpc.AuthenticateRequest_TOKEN,
			Token: maybeAuthToken.Token,
		}
	case len(m.config.PreSharedKey) != 0:
		m.logger.Debug("using pre-shared key auth flow")
		req = rpc.AuthenticateRequest{
			Flow:         rpc.AuthenticateRequest_PRE_SHARED_KEY,
			PreSharedKey: m.config.PreSharedKey,
		}
	default:
		return nil, fmt.Errorf("must have a pre-shared key or an existing token to authenticate")
	}

	var err error
	req.SystemMetadata, err = m.metadataProvider.GetSystemMetadata(ctx)
	if err != nil {
		return nil, err
	}

	c, err := m.client.getAuthClientLocked(ctx)
	if err != nil {
		return nil, err
	}
	res, err := c.Authenticate(ctx, &req)
	if err != nil {
		// If we got an unauthenticated error from the server and tried using our
		// existing token, it probably means the existing token is expired, so let's
		// revert back to the pre-shared key.
		//
		// This is a recursive call, but we know we can't get into an infinite loop
		// because the req.Flow will never be TOKEN if there's no auth token provided
		// which we're not providing.
		if s, ok := status.FromError(err); ok && s.Code() == codes.Unauthenticated {
			if req.Flow == rpc.AuthenticateRequest_TOKEN {
				m.logger.Debug("existing token is expired, reverting to pre-shared key")
				return m.newToken(ctx, nil)
			}
		}
		return nil, err
	}
	exp, err := token.GetExpiration(res.Token)
	if err != nil {
		return nil, err
	}
	m.logger.Debug("successfully obtained auth token", zap.Time("exp", exp), zap.String("dur", exp.Sub(time.Now()).String()))
	return &AuthToken{
		Token:      res.Token,
		Expiration: exp,
		Duration:   exp.Sub(time.Now()),
	}, nil
}

func newTokenManager(config *Config, db Database, logger *zap.Logger, client *Client, metadataProvider SystemMetadataProvider) *tokenManager {
	return &tokenManager{
		metadataProvider: metadataProvider,
		logger:           logger.Named("token-manager"),
		config:           config,
		client:           client,
		db:               db,
	}
}

func isExpired(exp time.Time) bool {
	return time.Now().After(exp)
}

// isExpiringSoon returns true if the token expires within 30% of it's total duration
func isExpiringSoon(now time.Time, exp time.Time, dur time.Duration) bool {
	exp = exp.Add(-(dur / 3))
	return exp.Before(now) || exp.Equal(now)
}
