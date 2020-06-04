package cmd

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/stephenafamo/janus/migrator"
	"github.com/stephenafamo/orchestra"
)

// CMD is our multipurpose cli struct
type CMD struct {
	Name, Slug, Version string

	Migrator migrator.Interface
	Workers  map[string]orchestra.Player

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

	if len(c.Workers) > 0 {
		rootCmd.AddCommand(&cobra.Command{
			Use:   "run",
			Short: fmt.Sprintf("Run a %s worker", c.Name),
			Long:  fmt.Sprintf("Run a %s worker", c.Name),
			RunE: func(cmd *cobra.Command, args []string) error {

				conductor := &orchestra.Conductor{
					Timeout: 5 * time.Second,
					Players: make(map[string]orchestra.Player),
				}

				// Start all if no args were given
				if len(args) == 0 {
					conductor.Players = c.Workers
				}

				for _, pl := range args {
					player, ok := c.Workers[pl]
					if ok {
						conductor.Players[pl] = player
					}
				}

				return orchestra.PlayUntilSignal(
					conductor,
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

// GetCobra gets the built *cobra.Command runs the root command.
func (c *CMD) GetCobra() *cobra.Command {
	c.buildCMD()
	return c.cmd
}
