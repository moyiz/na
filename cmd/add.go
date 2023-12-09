package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/moyiz/na/internal/config"
	"github.com/moyiz/na/internal/consts"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:                   "add [name ...] [--] [command ...]",
	DisableFlagsInUseLine: true,
	Short:                 "Adds a nested alias to a configuration",
	Long: consts.Logo + `
Adds a nested alias to configuration.
If the target command consists more than a single word, an optional '--' will
act as a delimiter between the alias and the command.
It will create the config directory if at does not exist.`,
	Aliases: []string{"a"},
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		configFile := ActiveConfigFile()

		// Create directories if not exist
		configDir := path.Dir(configFile)
		if _, err := os.Stat(configDir); err != nil {
			os.MkdirAll(configDir, 0775)
		}

		// Separate alias key and command at `--` or set the command to the last argument
		sep := cmd.ArgsLenAtDash()
		if sep == 0 {
			fmt.Println("na: add: '--' cannot be the first argument.")
			os.Exit(1)
		} else if sep < 0 {
			// No `--`, the last argument is the command
			sep = len(args) - 1
		}

		c := config.GetFromFiles(configFile)
		if err := c.SetAlias(args[:sep], args[sep:]); err != nil {
			fmt.Println("na: add:", err)
			os.Exit(1)
		}
	},
}
