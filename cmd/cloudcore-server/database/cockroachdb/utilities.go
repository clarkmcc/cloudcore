package cockroachdb

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

func getReturningID(rows *sqlx.Rows) (id string, err error) {
	if !rows.Next() {
		return "", errors.New("expected to get an ID")
	}
	return id, rows.Scan(&id)
}
