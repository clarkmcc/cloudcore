package cockroachdb

import (
	"errors"
	"github.com/jmoiron/sqlx"
)

func getReturningID(rows *sqlx.Rows) (id string, err error) {
	if !rows.Next() {
		return "", errors.New("expected to get an ID")
	}
	err = rows.Scan(&id)
	if err != nil {
		return "", err
	}
	return id, rows.Close()
}
