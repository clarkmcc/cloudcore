package agent

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrAuthTokenNotFound is returned when the auth token is not found
	// in the database indicating that the agent does not have a token.
	ErrAuthTokenNotFound = errors.New("auth token not found")
	ErrNoAgentID         = errors.New("no agent id")
)

type Database interface {
	AuthToken(ctx context.Context) (*AuthToken, error)
	SaveAuthToken(ctx context.Context, token *AuthToken) error

	AgentID(ctx context.Context) (string, error)
	SaveAgentID(ctx context.Context, agentID string) error
}

func NewDatabase(cfg *Config) (Database, error) {
	switch cfg.Database.Flavor {
	case databaseFlavorMemory:
		return newMemoryDB()
	default:
		return nil, fmt.Errorf("unknown database flavor: %s", cfg.Database.Flavor)
	}
}
