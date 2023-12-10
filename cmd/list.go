package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/moyiz/na/internal/config"
	"github.com/moyiz/na/internal/consts"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:                   "list [name ...]",
	DisableFlagsInUseLine: true,
	Short:                 "Lists aliases and target commands",
	Long: consts.Logo + `
Lists all aliases under optional given partial prefix.
The output format is a list of full alias names and their target commands,
separated by double-dash (--).
By default, all configuration files are merged for this command.`,
	Aliases:           []string{"l", "ls"},
	ValidArgsFunction: validListArgs,
	Run:               listRun,
}

func validListArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if slices.Contains(os.Args, "--") {
		return []string{}, cobra.ShellCompDirectiveDefault
	}
	config.GetFromFiles(AllConfigFiles()...)
	return config.ListNextParts(args), cobra.ShellCompDirectiveNoFileComp
}

func listRun(cmd *cobra.Command, args []string) {
	for _, alias := range config.ListAliases(args...) {
		fmt.Println(alias.Name, "--", alias.Command)
	}
}
