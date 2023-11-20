package cockroachdb

import (
	"context"
	"database/sql"
	"errors"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database/types"
)

// ListProjectHosts returns all the hosts in a project.
// todo: pagination, sorting, etc...
func (d *Database) ListProjectHosts(ctx context.Context, projectId string) (out []types.Host, err error) {
	err = d.db.SelectContext(ctx, &out, `
		SELECT a.last_heartbeat_timestamp, a.online, h.id, h.status, h.created_at, h.updated_at, identifier, hostname, public_ip_address, private_ip_address, os_name, os_family, os_version, kernel_architecture, kernel_version, cpu_model, cpu_cores FROM host h
		INNER JOIN (
			SELECT
				DISTINCT ON (host_id) host_id,
				last_heartbeat_timestamp,
				online
			FROM
				agent
			ORDER BY
				host_id, last_heartbeat_timestamp DESC
		) a ON h.id = a.host_id
		WHERE h.project_id = $1`, projectId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	return out, nil
}

func (d *Database) GetHost(ctx context.Context, hostId, projectId string) (out types.Host, err error) {
	return out, d.db.GetContext(ctx, &out, `SELECT * FROM host WHERE id = $1 AND project_id = $2`, hostId, projectId)
}

func (d *Database) GetEventLogsByHost(ctx context.Context, hostId string, limit int) (out []types.AgentEventLog, err error) {
	return out, d.db.SelectContext(ctx, &out, `SELECT * FROM agent_event WHERE host_id = $1 ORDER BY created_at DESC LIMIT $2`, hostId, limit)
}
