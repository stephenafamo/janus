package migrator

import (
	"fmt"
	"strings"
)

// Migrate runs the migrations with the given implementation
func Migrate(m Interface, action string, limit int) (int, error) {
	var count int
	var err error

	switch strings.ToLower(action) {
	case "down":
		count, err = m.Down(limit)
	case "up":
		count, err = m.Up(limit)
	default:
		err = fmt.Errorf("Unknown migration action specified")
	}

	if err != nil {
		return count, fmt.Errorf("Could not carry out the %q action: %w", action, err)
	}

	return count, nil
}
