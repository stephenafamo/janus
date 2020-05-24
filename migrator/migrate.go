package migrator

import (
	"fmt"
	"log"
	"strings"
)

// Migrate runs the migrations from files in the migrations folder (relative)
func Migrate(m Interface, action string) error {
	var err error

	switch strings.ToLower(action) {
	case "drop":
		err = m.Drop()
	case "down":
		err = m.Down()
	case "up":
		err = m.Up()
	default:
		err = fmt.Errorf("Unknown migration action specified")
	}

	if err != nil {
		return fmt.Errorf("Could not carry out the %q action... %v", action, err)
	}

	log.Printf("Successfully carried out the %q action", action)

	return nil
}
