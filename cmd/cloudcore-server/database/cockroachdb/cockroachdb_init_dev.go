//go:build dev

package cockroachdb

import (
	"context"
	_ "github.com/lib/pq"
)

func (d *Database) loadTestData(ctx context.Context) error {
	d.logger.Info("loading dev mode data")
	_, err := d.db.ExecContext(ctx, `UPDATE global_state SET is_dev_mode = true WHERE true;`)
	if err != nil {
		return err
	}
	_, err = d.db.ExecContext(ctx, `INSERT INTO agent_psk (project_id, name, key, uses_remaining) VALUES ((SELECT project.id FROM project INNER JOIN tenant ON tenant.id = project.tenant_id WHERE tenant.name = 'Default'), 'Dev Mode PSK', '00000000-0000-0000-0000-000000000000', NULL) ON CONFLICT DO NOTHING;`)
	return err
}
