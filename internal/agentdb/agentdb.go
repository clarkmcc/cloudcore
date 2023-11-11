package agentdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/clarkmcc/cloudcore/internal/config"
	"github.com/google/uuid"
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

func New(ctx context.Context, cfg *config.AgentConfig) (AgentDB, error) {
	db, err := newDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("new db: %w", err)
	}

	// If we don't have an Agent ID, then generate one
	id, err := db.AgentID(ctx)
	if err != nil && !errors.Is(err, ErrNoAgentID) {
		return nil, err
	}
	if len(id) > 0 {
		return db, nil
	}
	return db, db.SaveAgentID(ctx, uuid.New().String())
}

func newDB(cfg *config.AgentConfig) (AgentDB, error) {
	switch cfg.Database.Flavor {
	case config.AgentDatabaseFlavorMemory:
		return newMemoryDB()
	default:
		return nil, fmt.Errorf("unknown database flavor: %s", cfg.Database.Flavor)
	}
}
