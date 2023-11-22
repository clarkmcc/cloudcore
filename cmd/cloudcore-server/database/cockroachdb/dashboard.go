package cockroachdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database/types"
)

func (d *Database) GetProjectMetrics(ctx context.Context, projectId string) (*types.ProjectMetrics, error) {
	hosts, err := d.getHostsByPlatform(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("getting hosts by platform: %w", err)
	}

	var totalHosts int
	err = d.db.GetContext(ctx, &totalHosts, `SELECT COUNT(*) FROM host WHERE status = 'active' AND project_id = $1`, projectId)
	if err != nil {
		return nil, fmt.Errorf("getting total hosts: %w", err)
	}

	var totalAgents int
	err = d.db.GetContext(ctx, &totalAgents, `SELECT COUNT(*) FROM agent WHERE status = 'active' AND project_id = $1`, projectId)
	if err != nil {
		return nil, fmt.Errorf("getting total agents: %w", err)
	}

	online, offline, err := d.getOfflineOnlineHostCount(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("getting online/offline host count: %w", err)
	}

	return &types.ProjectMetrics{
		HostsByOsName: hosts,
		TotalAgents:   totalAgents,
		TotalHosts:    totalHosts,
		OnlineHosts:   online,
		OfflineHosts:  offline,
	}, nil
}

// getHostsByPlatform returns a map of os_name to count of hosts with that os_name.
func (d *Database) getHostsByPlatform(ctx context.Context, projectId string) (out []types.OsNameCount, err error) {
	err = d.db.SelectContext(ctx, &out, `
		SELECT os_name, COUNT(*) AS count FROM host
		WHERE status = 'active' AND project_id = $1
		GROUP BY os_name`, projectId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return out, nil
}

func (d *Database) getOfflineOnlineHostCount(ctx context.Context, projectId string) (online, offline int, err error) {
	type result struct {
		OnlineCount  int `db:"online_count"`
		OfflineCount int `db:"offline_count"`
	}

	var results result
	err = d.db.GetContext(ctx, &results, `
		SELECT COUNT(CASE WHEN a.online AND a.last_heartbeat_timestamp > NOW() - INTERVAL '1 minute' THEN 1 END) AS online_count,
			   COUNT(CASE WHEN NOT (a.online AND a.last_heartbeat_timestamp > NOW() - INTERVAL '1 minute') THEN 1 END) AS offline_count
		FROM host h
			 INNER JOIN (SELECT DISTINCT ON (host_id) host_id,
													  last_heartbeat_timestamp,
													  online
						 FROM agent
						 WHERE project_id = $1
						 ORDER BY host_id, last_heartbeat_timestamp DESC
						 ) a ON h.id = a.host_id
		WHERE h.project_id = $1`, projectId)
	if err != nil {
		return 0, 0, err
	}
	return results.OnlineCount, results.OfflineCount, nil
}
