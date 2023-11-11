package database

import (
	"context"
	"fmt"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database/cockroachdb"
	"github.com/clarkmcc/cloudcore/internal/config"
	"github.com/clarkmcc/cloudcore/internal/rpc"
)

type Database interface {
	Migrate() error

	// UpdateMetadata upserts the host and agent with all the associated host metadata
	UpdateMetadata(context.Context, *rpc.SystemMetadata) (string, error)
}

func New(cfg *config.ServerConfig) (Database, error) {
	switch cfg.Database.Type {
	case config.DatabaseTypeCockroachDB:
		return cockroachdb.New(cfg)
	default:
		return nil, fmt.Errorf("unknown database type: %s", cfg.Database.Type)
	}
}
