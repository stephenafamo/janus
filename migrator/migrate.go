package migrator

import (
	"context"
	"fmt"
	"strings"
	"testing"
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
		err = fmt.Errorf("unknown migration action specified")
	}

	if err != nil {
		return count, fmt.Errorf("could not carry out the %q action: %w", action, err)
	}

	return count, nil
}

type migAction struct {
	Direction string
	Count     int
}

// Test tests that the migrations are working properly
// It starts by running the down migrations, then the up migrations
// Then it runs the down migrations in a stepped pattern to make sure that
// each migration cleans up after itself properly
//
// NOTE: Do not use this in production, it is only for testing
// It is recommended to make sure you are connected to the test database
func Test(ctx context.Context, tb testing.TB, mig Migrator) {
	tb.Helper()
	// Run the down migrations
	_, err := Migrate(ctx, mig, "down", 0)
	if err != nil {
		tb.Fatalf("error running %s-%d migration: %v", "down", 0, err)
	}

	// Run the up migrations
	count, err := Migrate(ctx, mig, "up", 0)
	if err != nil {
		tb.Fatalf("error running %s-%d migration: %v", "up", 0, err)
	}

	// run the down migrations in a stepped pattern to make sure that it is done
	for i := count; i > 0; i-- {
		// for each step, to make sure it cleans up properly after itself
		stepActions := []migAction{
			{
				Direction: "down",
				Count:     1,
			},
			{
				Direction: "up",
				Count:     1,
			},
			{
				Direction: "down",
				Count:     1,
			},
		}

		for _, action := range stepActions {
			_, err = Migrate(ctx, mig, action.Direction, action.Count)
			if err != nil {
				tb.Fatalf("error running %s-%d migration: %v", action.Direction, action.Count, err)
			}
		}
	}
}
