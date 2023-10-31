package migrator

import "context"

// Interface in the main Migration tool.
type Interface interface {
	Up(ctx context.Context, limit int) (int, error)
	Down(ctx context.Context, limit int) (int, error)
}

// Migrator is a migrator
type Migrator = Interface
