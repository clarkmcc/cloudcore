package cockroachdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database/types"
)

// UpsertUser creates a new user if one doesn't exist and returns all the user's projects.
func (d *Database) UpsertUser(ctx context.Context, subject string) ([]types.Project, error) {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer handleRollback(&err, tx)

	// Check if the user exists first
	var user types.User
	err = tx.Get(&user, `SELECT * FROM "user" WHERE subject = $1;`, subject)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == nil {
		return d.getUserProjects(ctx, user.ID)
	}

	// At this point, we're creating a new tenant, project, and user
	var tenantId string
	err = tx.QueryRowx(`INSERT INTO tenant (name) VALUES ($1) RETURNING id`, subject).Scan(&tenantId)
	if err != nil {
		return nil, fmt.Errorf("creating tenant: %w", err)
	}

	var project types.Project
	err = tx.QueryRowx(`INSERT INTO project (name, tenant_id) VALUES ('Default', $1) RETURNING *`, tenantId).StructScan(&project)
	if err != nil {
		return nil, fmt.Errorf("creating project: %w", err)
	}

	var userId string
	err = tx.QueryRowx(`INSERT INTO "user" (subject, tenant_id) VALUES ($1, $2) RETURNING id;`, subject, tenantId).Scan(&userId)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	// Associated user with project and make this project the user's default
	_, err = tx.Exec(`INSERT INTO user_project (user_id, project_id) VALUES ($1, $2)`, userId, project.ID)
	if err != nil {
		return nil, fmt.Errorf("adding user to project: %w", err)
	}

	return []types.Project{project}, tx.Commit()
}

// CanAccessProject returns true if a user can access a given project ID.
func (d *Database) CanAccessProject(ctx context.Context, subject, projectId string) (bool, error) {
	err := d.db.QueryRowxContext(ctx, `
		SELECT * FROM "user" u 
		    INNER JOIN user_project p ON p.user_id = u.id 
		WHERE u.subject = $1 AND p.project_id = $2;`,
		subject, projectId).Err()
	if err == nil {
		return true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return false, err
}
