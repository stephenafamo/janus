package migrator

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
)

// Migrator is a migrator
type Migrator struct {
	m *migrate.Migrate
}

// New creates a new Migrator
func New(m *migrate.Migrate) *Migrator {
	return &Migrator{m: m}
}

// Cmd returns a cobra command for migrations
func (m *Migrator) Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run the database migration",
		Long: `Run the database migration.
	
	By default, running this command will run the "up" migration action.
	Other actiions are "down" and "drop".
	
	To run the down migrations, use:
	
	iteretta migrate down
	
	
	You can supply multiple migrations and they will be run in order
	For example, to run "drop" and then "up" do:
	
	iteretta migrate drop up
		`,
		Args: cobra.MaximumNArgs(1),
		RunE: m.do,
	}
}

func (m *Migrator) do(cmd *cobra.Command, args []string) error {

	action := "up"

	if len(args) > 0 {
		action = args[0]
	}

	err := Migrate(m.m, action)
	if err != nil {
		return err
	}

	return nil
}
