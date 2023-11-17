package cockroachdb

import (
	"context"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database/types"
)

func (d *Database) getProject(ctx context.Context, id string) (*types.Project, error) {
	var project types.Project
	return &project, d.db.GetContext(ctx, &project, `SELECT * FROM project WHERE id = $1;`, id)
}

func (d *Database) GetUserProjects(ctx context.Context, subject string) ([]types.Project, error) {
	var projects []types.Project
	return projects, d.db.SelectContext(ctx, &projects, `
		SELECT p.* FROM project p
		INNER JOIN user_project up ON p.id = up.project_id AND up.status = 'active'
	   	INNER JOIN "user" u ON up.user_id = u.id AND u.status = 'active'
		WHERE u.subject = $1
		ORDER BY p.updated_at DESC;
	`, subject)
}

func (d *Database) getUserProjects(ctx context.Context, userId string) ([]types.Project, error) {
	var projects []types.Project
	return projects, d.db.SelectContext(ctx, &projects, `
		SELECT p.* FROM project p
		INNER JOIN user_project up ON p.id = up.project_id AND up.status = 'active'
	   	INNER JOIN "user" u ON up.user_id = u.id AND u.status = 'active'
		WHERE up.user_id = $1
		ORDER BY p.updated_at DESC;
	`, userId)
}

// CreateProject creates a new project with the given name and description.
func (d *Database) CreateProject(ctx context.Context, subject, name, description string) (types.Project, error) {
	// Find the user ID and tenant ID of the subject
	var user types.User
	err := d.db.GetContext(ctx, &user, `SELECT id, tenant_id FROM "user" WHERE subject = $1;`, subject)
	if err != nil {
		return types.Project{}, err
	}

	var project types.Project
	err = d.db.GetContext(ctx, &project, `
		INSERT INTO project (tenant_id, name, description)
		VALUES ($1, $2, $3)
		RETURNING *;
	`, user.TenantID, name, description)
	if err != nil {
		return types.Project{}, err
	}

	// Add the user to the project
	_, err = d.db.ExecContext(ctx, `
		INSERT INTO user_project (user_id, project_id, status)
		VALUES ($1, $2, 'active');
	`, user.ID, project.ID)
	if err != nil {
		return types.Project{}, err
	}
	return project, err
}
