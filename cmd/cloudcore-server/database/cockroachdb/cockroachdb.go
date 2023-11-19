package cockroachdb

import (
	"context"
	"embed"
	"errors"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cockroachdb"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Database struct {
	db     *sqlx.DB
	dbName string
	logger *zap.Logger
}

func (d *Database) Migrate(ctx context.Context) error {
	fs, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}
	driver, err := cockroachdb.WithInstance(d.db.DB, &cockroachdb.Config{
		DatabaseName: d.dbName,
	})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", fs, d.dbName, driver)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return d.loadTestData(ctx)
}

func New(config *config.Config, logger *zap.Logger) (*Database, error) {
	db, err := sqlx.Open("postgres", config.Database.ConnectionString)
	if err != nil {
		return nil, err
	}
	return &Database{
		db:     db,
		dbName: config.Database.Name,
		logger: logger.Named("db"),
	}, db.Ping()
}
