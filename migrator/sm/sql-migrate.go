package sm

import (
	"database/sql"

	migrate "github.com/rubenv/sql-migrate"
	"github.com/stephenafamo/janus/migrator"
)

func Get(db *sql.DB, dialect string, source migrate.MigrationSource) migrator.Interface {
	return sm{}
}

type sm struct {
	db      *sql.DB
	dialect string
	source  migrate.MigrationSource
}

func (s sm) Up(limit int) error {
	_, err := migrate.ExecMax(s.db, s.dialect, s.source, migrate.Up, limit)
	return err
}

func (s sm) Down(limit int) error {
	_, err := migrate.ExecMax(s.db, s.dialect, s.source, migrate.Down, limit)
	return err
}
