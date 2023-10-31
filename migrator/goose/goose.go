package goose

import (
	"context"
	"database/sql"
	"errors"
	"io/fs"
	"math"

	"github.com/pressly/goose/v3"
	"github.com/stephenafamo/janus/migrator"
)

func Get(db *sql.DB, dialect string, source fs.FS, opts ...goose.OptionsFunc) migrator.Migrator {
	return impl{
		db:      db,
		dialect: dialect,
		source:  source,
		options: opts,
	}
}

type impl struct {
	db      *sql.DB
	dialect string
	source  fs.FS
	options []goose.OptionsFunc
}

func (g impl) setup() error {
	goose.SetBaseFS(g.source)

	if err := goose.SetDialect(g.dialect); err != nil {
		return err
	}

	return nil
}

func (g impl) Up(ctx context.Context, limit int) (int, error) {
	if err := g.setup(); err != nil {
		return 0, err
	}

	if limit == 0 {
		limit = math.MaxInt64
	}

	for i := 0; i < limit; i++ {
		err := goose.UpByOneContext(ctx, g.db, ".", g.options...)
		if err != nil {
			if errors.Is(err, goose.ErrNoNextVersion) {
				return i, nil
			}
			return i, err
		}
	}

	return limit, nil
}

func (g impl) Down(ctx context.Context, limit int) (int, error) {
	if err := g.setup(); err != nil {
		return 0, err
	}

	if limit == 0 {
		limit = math.MaxInt64
	}

	for i := 0; i < limit; i++ {
		err := goose.DownContext(ctx, g.db, ".", g.options...)
		if err != nil {
			if errors.Is(err, goose.ErrNoCurrentVersion) {
				return i, nil
			}
			return i, err
		}
	}

	return limit, nil
}
