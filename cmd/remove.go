package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/moyiz/na/internal/cli"
	"github.com/moyiz/na/internal/config"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:                   "remove [prefix ...]",
	DisableFlagsInUseLine: true,
	Short:                 "Removes nested aliases from a configuration",
	Long: `Removes a nested alias from a configuration.
It accepts partial aliases paths to remove multiple aliases at once.
By default, the global (home directory config) configuration is used.

The short form of 'remove' is 'rm'.  
By default, the global configuration file is used for this command.  

Example na.yaml:  
my:
    alias1: ...  
    alias2: ...  
    another:  
        alias: ...  
	   
To delete 'my alias1':  
$ na rm my alias1  
To delete 'my alias1' and 'my alias2':  
$ na rm my al  
To delete 'my another alias':  
$ na rm my another  
To delete all:  
$ na rm my`,
	Aliases:           []string{"rm"},
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: validRemoveArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadFiles(ActiveConfigFile())
		if err := config.UnsetAlias(args...); err != nil {
			fmt.Println("na:", strings.Join(args, " ")+":", err)
			os.Exit(1)
		}
	},
}

func validRemoveArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if slices.Contains(os.Args, "--") {
		return []string{}, cobra.ShellCompDirectiveDefault
	}
	config.LoadFiles(ActiveConfigFile())
	return cli.ListNextParts(config.ListAliases(), args), cobra.ShellCompDirectiveNoFileComp
}
