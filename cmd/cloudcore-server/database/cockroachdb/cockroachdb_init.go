//go:build !dev

package cockroachdb

import (
	"context"
	_ "github.com/lib/pq"
)

func (d *Database) loadTestData(ctx context.Context) error {
	_, err := d.db.ExecContext(ctx, `UPDATE global_state SET is_dev_mode = false WHERE true;`)
	return err
}
