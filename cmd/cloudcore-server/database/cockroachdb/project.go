package cockroachdb

import (
	"context"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database/types"
)

func (d *Database) getProject(ctx context.Context, id string) (*types.Project, error) {
	var project types.Project
	return &project, d.db.GetContext(ctx, &project, `SELECT * FROM project WHERE id = $1;`, id)
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
