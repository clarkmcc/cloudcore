package database

import (
	"context"
	"fmt"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/config"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database/cockroachdb"
	"github.com/clarkmcc/cloudcore/internal/rpc"
)

type Database interface {
	Migrate() error

	// UpdateMetadata upserts the host and agent with all the associated host metadata
	//UpdateMetadata(context.Context, *rpc.SystemMetadata) (string, error)

	// AuthenticateAgent accepts the agent metadata and the authentication pre-shared key
	// and returns the agent ID if the agent is authenticated. This function will upsert
	// the agent, the host, and add the agent to the appropriate groups.
	AuthenticateAgent(ctx context.Context, psk string, md *rpc.SystemMetadata) (string, error)
}

func New(cfg *config.Config) (Database, error) {
	switch cfg.Database.Type {
	case config.DatabaseTypeCockroachDB:
		return cockroachdb.New(cfg)
	default:
		return nil, fmt.Errorf("unknown database type: %s", cfg.Database.Type)
	}
}
