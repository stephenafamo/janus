package goose

import (
	"context"
	"errors"
	"math"

	"github.com/pressly/goose/v3"
	"github.com/stephenafamo/janus/migrator"
)

func Get(provider *goose.Provider) migrator.Migrator {
	return impl{
		provider: provider,
	}
}

type impl struct {
	provider *goose.Provider
}

func (g impl) Up(ctx context.Context, limit int) (int, error) {
	if limit == 0 {
		limit = math.MaxInt64
	}

	for i := 0; i < limit; i++ {
		_, err := g.provider.UpByOne(ctx)
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
	if limit == 0 {
		limit = math.MaxInt64
	}

	for i := 0; i < limit; i++ {
		_, err := g.provider.Down(ctx)
		if err != nil {
			if errors.Is(err, goose.ErrNoNextVersion) {
				return i, nil
			}
			return i, err
		}
	}

	return limit, nil
}
