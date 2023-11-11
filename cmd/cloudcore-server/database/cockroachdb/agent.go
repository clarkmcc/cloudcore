package cockroachdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database/types"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"
	"time"
)

var (
	ErrAgentNotFound        = errors.New("agent not found")
	ErrPreSharedKeyNotFound = errors.New("pre-shared key not found")
)

func (d *Database) AuthenticateAgent(ctx context.Context, key string, metadata *rpc.SystemMetadata) (agentID string, err error) {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer handleRollback(&err, tx)

	// Make sure the psk exists and is usable
	var psk types.PreSharedKey
	err = tx.GetContext(ctx, &psk, `
		SELECT id, project_id, created_at, updated_at, name, key, status, uses_remaining, expiration FROM agent_psk 
		WHERE key = $1
		  	AND status = 'active'
			AND uses_remaining > 0
			AND expiration > NOW()
		LIMIT 1
	`, key)
	if err != nil {
		return "", fmt.Errorf("getting pre-shared key: %w", err)
	}

	// Decrement the uses remaining
	_, err = tx.ExecContext(ctx, `
		UPDATE agent_psk SET uses_remaining = uses_remaining - 1 WHERE id = $1;
	`, psk.ID)
	if err != nil {
		return "", fmt.Errorf("updating pre-shared key uses remaining: %w", err)
	}

	// Check for any PSK groups
	var groups []types.AgentGroup
	err = tx.SelectContext(ctx, &groups, `
		SELECT g.* FROM agent_group g
		INNER JOIN agent_group_psk gp ON g.id = gp.agent_group_id AND gp.status = 'active'
		INNER JOIN agent_psk p ON gp.agent_psk_id = p.id AND p.status = 'active'
		WHERE p.id = $1 AND g.status = 'active';
	`, psk.ID)
	if err != nil {
		return "", fmt.Errorf("finding agent groups: %w", err)
	}

	// First, we upsert the host on the identifier field
	rows, err := tx.NamedQuery(`
		INSERT INTO host (project_id, identifier, hostname, host_id, public_ip_address, os_name, os_family, os_version, kernel_architecture, kernel_version, cpu_model, cpu_cores)
		VALUES(:project_id, :identifier, :hostname, :host_id, :public_ip_address, :os_name, :os_family, :os_version, :kernel_architecture, :kernel_version, :cpu_model, :cpu_cores)
		ON CONFLICT (identifier) DO UPDATE SET hostname = :hostname, host_id = :host_id, public_ip_address = :public_ip_address, os_name = :os_name, os_family = :os_family, os_version = :os_version, kernel_architecture = :kernel_architecture, kernel_version = :kernel_version, cpu_model = :cpu_model, cpu_cores = :cpu_cores
		RETURNING id
	`, map[string]any{
		"project_id":          psk.ProjectID,
		"identifier":          metadata.GetIdentifiers().GetHostIdentifier(),
		"hostname":            metadata.GetIdentifiers().GetHostname(),
		"host_id":             metadata.GetIdentifiers().GetHostId(),
		"public_ip_address":   metadata.GetIdentifiers().GetPublicIpAddress(),
		"os_name":             metadata.GetOs().GetName(),
		"os_family":           metadata.GetOs().GetFamily(),
		"os_version":          metadata.GetOs().GetVersion(),
		"kernel_architecture": metadata.GetKernel().GetArch(),
		"kernel_version":      metadata.GetKernel().GetVersion(),
		"cpu_model":           metadata.GetCpu().GetModel(),
		"cpu_cores":           metadata.GetCpu().GetCores(),
	})
	if err != nil {
		return "", err
	}
	hostID, err := getReturningID(rows)
	if err != nil {
		return "", err
	}

	// Upsert the agent now
	agentID = metadata.GetIdentifiers().GetAgentIdentifier()
	if len(agentID) == 0 {
		rows, err = tx.NamedQuery(`
		INSERT INTO agent (project_id, host_id, online, last_heartbeat_timestamp)
		VALUES (:project_id, :host_id, :online, :last_heartbeat_timestamp)
		RETURNING id
	`, map[string]any{
			"project_id":               psk.ProjectID,
			"host_id":                  hostID,
			"online":                   true,
			"last_heartbeat_timestamp": time.Now(),
		})
		if err != nil {
			return "", err
		}
		agentID, err = getReturningID(rows)
		if err != nil {
			return "", err
		}
	} else {
		// Agent already exists, just update the last heartbeat
		_, err = tx.ExecContext(ctx, `
			UPDATE agent SET online = true, last_heartbeat_timestamp = $1 WHERE id = $2;
		`, time.Now(), agentID)
		if err != nil {
			return "", fmt.Errorf("updating agent heartbeat: %w", err)
		}
	}

	// Add the agent to the groups
	for _, g := range groups {
		_, err = tx.NamedExecContext(ctx, `
			INSERT INTO agent_group_member (project_id, agent_id, agent_group_id)
			VALUES (:project_id, :agent_id, :agent_group_id);
		`, map[string]any{
			"project_id":     psk.ProjectID,
			"agent_id":       agentID,
			"agent_group_id": g.ID,
		})
		if err != nil {
			return "", fmt.Errorf("adding agent to group %s (%s): %w", g.Name, g.ID, err)
		}
	}
	return agentID, tx.Commit()
}

func handleRollback(err *error, tx *sqlx.Tx) {
	if *err != nil {
		multierr.AppendFunc(err, tx.Rollback)
	}
}
