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
	Use:                   "remove [-p|--prefix] <path ...>",
	DisableFlagsInUseLine: true,
	Short:                 "Removes nested aliases from a configuration",
	Long: `Removes nested aliases from a configuration.  
It supports two modes for selecting aliases to remove:  
- Exact matching [default] - will select all aliases associated with the  
  given path.  
- Prefix matching (-p|--prefix) - will select any alias that is prefixed by  
  the given path.  

The short form of 'remove' is 'rm'.  
By default, the global configuration file is used for this command.  

Given the following configuration:  
my:  
    alias1: ...  
    alias2: ...  
    another:  
        alias: ...  
  
$ na rm my alias1 # Will remove 'my alias1'  
$ na rm my al # Will not match with any alias  
$ na rm -p my al # Will match 'my alias1' and 'my alias2'`,
	Aliases:           []string{"rm"},
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: validRemoveArgs,
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadFiles(ActiveConfigFile())
		var unsetFunc func(...string) error
		if byPrefix, _ := cmd.Flags().GetBool("prefix"); byPrefix {
			unsetFunc = config.UnsetAliasByPrefix
		} else {
			unsetFunc = config.UnsetAlias
		}
		if err := unsetFunc(args...); err != nil {
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

func init() {
	removeCmd.Flags().BoolP("prefix", "p", false, "Remove aliases by prefix instead of exact match")
}
