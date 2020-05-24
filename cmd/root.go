package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
	"github.com/stephenafamo/janus/migrator"
	"github.com/stephenafamo/orchestra"
)

// CMD is our multipurpose cli struct
type CMD struct {
	Name, Slug, Version string

	Migrator *migrate.Migrate
	SeedFunc func(cmd *cobra.Command, args []string) error

	Worker orchestra.Player

	cmd *cobra.Command
}

var verbose bool
var debug bool

func (c *CMD) buildCMD() {

	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:     c.Slug,
		Short:   c.Name,
		Long:    c.Name,
		Version: c.Version,
	}

	if c.cmd != nil {
		rootCmd = c.cmd
	}

	if c.Migrator != nil {
		rootCmd.AddCommand(migrator.New(c.Migrator).Cmd())
	}

	if c.SeedFunc != nil {
		rootCmd.AddCommand(&cobra.Command{
			Use:   "seed",
			Short: "Run the database seeder",
			// Long:  "Run the database seeder. Be careful, this will wipe the existing data.",
			Long: fmt.Sprintf(`Run the database seeder`),
			Args: cobra.ArbitraryArgs,
			RunE: c.SeedFunc,
		})
	}

	if c.Worker != nil {
		rootCmd.AddCommand(&cobra.Command{
			Use:   "start",
			Short: fmt.Sprintf("Start %s", c.Name),
			Long:  fmt.Sprintf("Start %s", c.Name),
			RunE: func(cmd *cobra.Command, args []string) error {
				return orchestra.PlayUntilSignal(
					c.Worker,
					os.Interrupt, syscall.SIGTERM,
				)
			},
		})
	}

	c.cmd = rootCmd
}

// Execute runs the root command.
func (c *CMD) Execute() error {
	c.buildCMD()
	return c.cmd.Execute()
}

// GetCobra gets the build *cobra.Command runs the root command.
func (c *CMD) GetCobra() *cobra.Command {
	c.buildCMD()
	return c.cmd
}
