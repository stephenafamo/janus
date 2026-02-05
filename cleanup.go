package janus

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"
)

type StopFunc = func(context.Context) error

func NoopStopFunc(_ context.Context) error {
	return nil
}

func NoContextStopFunc(f func() error) StopFunc {
	return func(_ context.Context) error {
		return f()
	}
}

func Stop(ctx context.Context, timeout time.Duration, stop StopFunc) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), timeout)
	defer cancel()

	if err := stop(ctx); err != nil {
		slog.Default().LogAttrs(
			ctx, slog.LevelError, "error during cleanup",
			slog.String("error", err.Error()),
		)
	}
}

func CombineStopFuncs(stops ...StopFunc) StopFunc {
	return func(ctx context.Context) error {
		var wg sync.WaitGroup
		wg.Add(len(stops))

		errs := make([]error, len(stops))
		for i, stop := range stops {
			go func(i int, stop StopFunc) {
				defer wg.Done()
				errs[i] = stop(ctx)
			}(i, stop)
		}

		return errors.Join(errs...)
	}
}
