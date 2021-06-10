package gm

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/stephenafamo/janus/migrator"
)

func Get(m *migrate.Migrate) migrator.Interface {
	return gm{m}
}

type gm struct {
	m *migrate.Migrate
}

func (g gm) Up(limit int) error {
	if limit == 0 {
		return g.m.Up()
	}
	return g.m.Steps(limit)
}

func (g gm) Down(limit int) error {
	if limit == 0 {
		return g.m.Down()
	}
	return g.m.Steps(-1 * limit)
}
