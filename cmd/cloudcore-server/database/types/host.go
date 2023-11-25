package types

import (
	"database/sql"
	"time"
)

type Host struct {
	ID                     string         `db:"id" json:"id,omitempty"`
	CreatedAt              time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt              time.Time      `db:"updated_at" json:"updatedAt"`
	Status                 Status         `db:"status" json:"status,omitempty"`
	ProjectID              string         `db:"project_id" json:"projectId,omitempty"`
	HostID                 string         `db:"host_id" json:"hostId,omitempty"` // The host's system UUID
	Identifier             string         `db:"identifier" json:"identifier,omitempty"`
	LastHeartbeatTimestamp time.Time      `db:"last_heartbeat_timestamp" json:"lastHeartbeatTimestamp"`
	Online                 bool           `db:"online" json:"online,omitempty"`
	Hostname               sql.NullString `db:"hostname" json:"hostname"`
	PublicIPAddress        sql.NullString `db:"public_ip_address" json:"publicIPAddress"`
	PrivateIPAddress       sql.NullString `db:"private_ip_address" json:"privateIPAddress"`
	OSName                 sql.NullString `db:"os_name" json:"osName"`
	OSFamily               sql.NullString `db:"os_family" json:"osFamily"`
	OSVersion              sql.NullString `db:"os_version" json:"osVersion"`
	KernelArchitecture     sql.NullString `db:"kernel_architecture" json:"kernelArchitecture"`
	KernelVersion          sql.NullString `db:"kernel_version" json:"kernelVersion"`
	CPUModel               sql.NullString `db:"cpu_model" json:"cpuModel"`
	CPUCores               sql.NullInt64  `db:"cpu_cores" json:"cpuCores"`
}

type HostGroup struct {
	ID          string         `db:"id"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	Status      Status         `db:"status"`
	ProjectID   string         `db:"project_id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
}
