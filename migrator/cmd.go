package migrator

// Interface in the main Migration tool.
type Interface interface {
	Up(limit int) (int, error)
	Down(limit int) (int, error)
}

// Migrator is a migrator
type Migrator = Interface
