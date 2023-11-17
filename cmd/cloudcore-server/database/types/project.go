package types

import (
	"database/sql"
	"time"
)

type Project struct {
	ID          string         `db:"id"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	Status      Status         `db:"status"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	TenantID    string         `db:"tenant_id"`
}
