package gm

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/stephenafamo/janus/migrator"
)

func Get(m *migrate.Migrate) migrator.Interface {
	return gm{m}
}

type gm struct {
	m *migrate.Migrate
}

func (g gm) Up(limit int) (int, error) {
	prev, _, err := g.m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return 0, err
	}

	if limit == 0 {
		err = g.m.Up()
	} else {
		err = g.m.Steps(limit)
	}
	if err != nil {
		return 0, err
	}

	current, _, err := g.m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return 0, err
	}

	return int(current - prev), nil
}

func (g gm) Down(limit int) (int, error) {
	prev, _, err := g.m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return 0, err
	}

	if limit == 0 {
		err = g.m.Down()
	} else {
		err = g.m.Steps(-1 * limit)
	}
	if err != nil {
		return 0, err
	}

	current, _, err := g.m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return 0, err
	}

	return int(prev - current), nil
}
