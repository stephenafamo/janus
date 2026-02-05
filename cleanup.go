package janus

import (
	"context"
	"log/slog"
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
