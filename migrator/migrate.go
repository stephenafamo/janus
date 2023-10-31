package migrator

import (
	"context"
	"fmt"
	"strings"
)

// Migrator in the main Migration tool.
type Migrator interface {
	Up(ctx context.Context, limit int) (int, error)
	Down(ctx context.Context, limit int) (int, error)
}

// Interface is a migrator
//
// Deprecated: Use Migrator instead
type Interface = Migrator

// Migrate runs the migrations with the given implementation
func Migrate(ctx context.Context, m Migrator, action string, limit int) (int, error) {
	var count int
	var err error

	switch strings.ToLower(action) {
	case "down":
		count, err = m.Down(ctx, limit)
	case "up":
		count, err = m.Up(ctx, limit)
	default:
		err = fmt.Errorf("Unknown migration action specified")
	}

	if err != nil {
		return count, fmt.Errorf("Could not carry out the %q action: %w", action, err)
	}

	return count, nil
}
