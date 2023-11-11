package types

import (
	"database/sql"
	"time"
)

type PreSharedKey struct {
	ID            string       `db:"id"`
	CreatedAt     time.Time    `db:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at"`
	Status        Status       `db:"status"`
	ProjectID     string       `db:"project_id"`
	Name          string       `db:"name"`
	Key           string       `db:"key"`
	UsesRemaining int32        `db:"uses_remaining"`
	Expiration    sql.NullTime `db:"expiration"`
}

type AgentGroup struct {
	ID          string         `db:"id"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	Status      Status         `db:"status"`
	ProjectID   string         `db:"project_id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
}
