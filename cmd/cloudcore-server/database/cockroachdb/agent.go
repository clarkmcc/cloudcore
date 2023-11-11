package cockroachdb

import (
	"context"
	"errors"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"go.uber.org/multierr"
	"time"
)

var (
	ErrAgentNotFound = errors.New("agent not found")
)

// UpdateMetadata updates the metadata for the host and agent with the given ID. If there is no agent ID, then
// we create the agent and return the ID.
func (d *Database) UpdateMetadata(ctx context.Context, metadata *rpc.SystemMetadata) (string, error) {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer multierr.AppendFunc(&err, tx.Rollback)

	// First, we upsert the host on the identifier field
	rows, err := d.db.NamedQueryContext(ctx, `
		INSERT INTO hosts (identifier, hostname, host_id, public_ip_address, os_name, os_family, os_version, kernel_architecture, kernel_version, cpu_model, cpu_cores)
		VALUES(:identifier, :hostname, :host_id, :public_ip_address, :os_name, :os_family, :os_version, :kernel_architecture, :kernel_version, :cpu_model, :cpu_cores)
		ON CONFLICT (identifier) DO UPDATE SET hostname = :hostname, host_id = :host_id, public_ip_address = :public_ip_address, os_name = :os_name, os_family = :os_family, os_version = :os_version, kernel_architecture = :kernel_architecture, kernel_version = :kernel_version, cpu_model = :cpu_model, cpu_cores = :cpu_cores
		RETURNING id
	`, map[string]any{
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

	// Next, we upsert the agent on the agent_id field
	rows, err = d.db.NamedQueryContext(ctx, `
		UPSERT INTO agents (id, host_id, online, last_heartbeat_timestamp)
		VALUES (:id, :host_id, :online, :last_heartbeat_timestamp)
		RETURNING id
	`, map[string]any{
		"id":                       metadata.GetIdentifiers().GetAgentIdentifier(),
		"host_id":                  hostID,
		"online":                   true,
		"last_heartbeat_timestamp": time.Now(),
	})
	if err != nil {
		return "", err
	}
	agentID, err := getReturningID(rows)
	if err != nil {
		return "", err
	}
	return agentID, tx.Commit()
}
