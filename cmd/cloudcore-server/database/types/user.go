package types

import "time"

type User struct {
	ID        string    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Status    Status    `db:"status"`
	Subject   string    `db:"subject"`
	TenantID  string    `db:"tenant_id"`
}
