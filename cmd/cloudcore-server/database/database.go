package database

import (
	"context"
	"fmt"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/config"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database/cockroachdb"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database/types"
	"github.com/clarkmcc/cloudcore/internal/rpc"
)

type Database interface {
	Migrate(ctx context.Context) error

	// UpdateMetadata upserts the host and agent with all the associated host metadata
	//UpdateMetadata(context.Context, *rpc.SystemMetadata) (string, error)

	// AuthenticateAgent accepts the agent metadata and the authentication pre-shared key
	// and returns the agent ID if the agent is authenticated. This function will upsert
	// the agent, the host, and add the agent to the appropriate groups.
	AuthenticateAgent(ctx context.Context, psk string, md *rpc.SystemMetadata) (string, error)

	AppDatabase
}

type AppDatabase interface {
	// UpsertUser takes an OIDC subject and upserts the user into the database
	// and returns the users default project.
	UpsertUser(ctx context.Context, subject string) ([]types.Project, error)

	// CreateProject takes an OIDC subject and the project details and creates a new project.
	CreateProject(ctx context.Context, subject, name, description string) (types.Project, error)

	// GetUserProjects takes an OIDC subject returns all the projects the user is a member of.
	GetUserProjects(ctx context.Context, subject string) ([]types.Project, error)

	// CanAccessProject returns true if an OIDC subject and access a specific project ID. This
	// is helpful for determining if a user has the authorization to access a project.
	CanAccessProject(ctx context.Context, subject, projectId string) (bool, error)

	// ListProjectHosts returns all the hosts in a project with the given project ID.
	ListProjectHosts(ctx context.Context, projectId string) ([]types.Host, error)
}

func New(cfg *config.Config, logger *zap.Logger) (Database, error) {
	switch cfg.Database.Type {
	case config.DatabaseTypeCockroachDB:
		return cockroachdb.New(cfg, logger)
	default:
		return nil, fmt.Errorf("unknown database type: %s", cfg.Database.Type)
	}
}
