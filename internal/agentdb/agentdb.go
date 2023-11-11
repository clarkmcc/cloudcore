package agentdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/clarkmcc/cloudcore/internal/config"
)

var (
	// ErrAuthTokenNotFound is returned when the auth token is not found
	// in the database indicating that the agent does not have a token.
	ErrAuthTokenNotFound = errors.New("auth token not found")
	ErrNoAgentID         = errors.New("no agent id")
)

type AgentDB interface {
	AuthToken(ctx context.Context) (*AuthToken, error)
	SaveAuthToken(ctx context.Context, token *AuthToken) error

	AgentID(ctx context.Context) (string, error)
	SaveAgentID(ctx context.Context, agentID string) error
}

func New(cfg *config.AgentConfig) (AgentDB, error) {
	switch cfg.Database.Flavor {
	case config.AgentDatabaseFlavorMemory:
		return newMemoryDB()
	default:
		return nil, fmt.Errorf("unknown database flavor: %s", cfg.Database.Flavor)
	}
}
