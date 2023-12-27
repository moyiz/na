package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/moyiz/na/internal/cli"
	"github.com/moyiz/na/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:                   "list [prefix ...]",
	DisableFlagsInUseLine: true,
	Short:                 "Lists aliases and target commands",
	Long: `Lists all aliases under optional given partial prefix.
The output format is a list of full alias names and their target commands,
separated by double-dash (--).

The short forms of 'list' are 'l' and 'ls'.  
By default, all configuration files are merged for this command.`,
	Aliases:           []string{"l", "ls"},
	ValidArgsFunction: validListArgs,
	Run:               listRun,
}

func validListArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if slices.Contains(os.Args, "--") {
		return []string{}, cobra.ShellCompDirectiveDefault
	}
	config.LoadFiles(AllConfigFiles()...)
	return cli.ListNextParts(config.ListAliases(), args), cobra.ShellCompDirectiveNoFileComp
}

func listRun(cmd *cobra.Command, args []string) {
	var aliases []config.Alias
	if byPrefix, _ := cmd.Flags().GetBool("prefix"); byPrefix {
		aliases = config.ListAliasesByPrefix(args...)
	} else {
		aliases = config.ListAliasesByPrefix(args...)
	}
	for _, alias := range aliases {
		fmt.Println(alias.Name, "--", alias.Command)
	}
}

func init() {
	listCmd.Flags().BoolP("prefix", "p", false, "List aliases by prefix instead of exact match")
}
