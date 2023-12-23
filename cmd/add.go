package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/moyiz/na/internal/config"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:                   "add [name ...] [--] [command ...]",
	DisableFlagsInUseLine: true,
	Short:                 "Adds a nested alias to a configuration",
	Long: `Adds or sets a new nested alias into configuration.
Without double-dash (--), the last argument will be treated as the target
command.
Otherwise, double-dash will act as the delimiter between the alias and the
target command.

Target commands may contain argument placeholders formatted as %n%. These
placeholders will be substituted with the given arguments respectively
(starting from 1). Placeholders can be negative to indicate the direction of
the argument list, e.g. %-1% will be substituted with the last argument.
The remaining arguments (those who were not candidates for substitution) will
be passed to the target command.

The short form of 'add' is 'a'.  
By default, the global configuration file is used for this command.  
The target configuration file and its parent directories will be created if not
exist.

Example na.yaml:  
my:  
    alias: echo %-1% %1%  

Example outputs:  
$ na run my alias -- a b  
b a  
$ na run my alias -- a b c d  
d a b c`,
	Aliases: []string{"a"},
	Args:    cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		configFile := ActiveConfigFile()

		// Create directories if not exist
		configDir := path.Dir(configFile)
		if _, err := os.Stat(configDir); err != nil {
			os.MkdirAll(configDir, 0o775)
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

		config.LoadFiles(configFile)
		if err := config.SetAlias(args[:sep], args[sep:]); err != nil {
			fmt.Println("na: add:", err)
			os.Exit(1)
		}
	},
}
