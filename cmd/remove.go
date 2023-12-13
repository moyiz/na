package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/moyiz/na/internal/cli"
	"github.com/moyiz/na/internal/config"
	"github.com/moyiz/na/internal/consts"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:                   "remove [name ...]",
	DisableFlagsInUseLine: true,
	Short:                 "Removes nested aliases from a configuration",
	Long: consts.Logo + `
Removes a nested alias from a configuration.
It accepts partial aliases paths to remove multiple aliases at once.
By default, the global (home directory config) configuration is used.`,
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
