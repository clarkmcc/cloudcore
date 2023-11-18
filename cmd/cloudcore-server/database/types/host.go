package types

import (
	"database/sql"
	"time"
)

type Host struct {
	ID                     string         `db:"id"`
	CreatedAt              time.Time      `db:"created_at"`
	UpdatedAt              time.Time      `db:"updated_at"`
	Status                 Status         `db:"status"`
	Identifier             string         `db:"identifier"`
	LastHeartbeatTimestamp time.Time      `db:"last_heartbeat_timestamp"`
	Online                 bool           `db:"online"`
	Hostname               sql.NullString `db:"hostname"`
	PublicIPAddress        sql.NullString `db:"public_ip_address"`
	PrivateIPAddress       sql.NullString `db:"private_ip_address"`
	OSName                 sql.NullString `db:"os_name"`
	OSFamily               sql.NullString `db:"os_family"`
	OSVersion              sql.NullString `db:"os_version"`
	KernelArchitecture     sql.NullString `db:"kernel_architecture"`
	KernelVersion          sql.NullString `db:"kernel_version"`
	CPUModel               sql.NullString `db:"cpu_model"`
	CPUCores               sql.NullInt64  `db:"cpu_cores"`
}
