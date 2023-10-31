package sm

import (
	"context"
	"database/sql"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/stephenafamo/janus/migrator"
)

func Get(db *sql.DB, dialect string, source migrate.MigrationSource) migrator.Interface {
	return sm{
		db:      db,
		dialect: dialect,
		source:  source,
	}
}

type sm struct {
	db      *sql.DB
	dialect string
	source  migrate.MigrationSource
}

func (s sm) Up(ctx context.Context, limit int) (int, error) {
	return migrate.ExecMaxContext(ctx, s.db, s.dialect, s.source, migrate.Up, limit)
}

func (s sm) Down(ctx context.Context, limit int) (int, error) {
	return migrate.ExecMaxContext(ctx, s.db, s.dialect, s.source, migrate.Down, limit)
}
