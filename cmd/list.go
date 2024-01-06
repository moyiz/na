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
	Use:                   "list [-o (simple|cmd|yaml)] [prefix ...]",
	DisableFlagsInUseLine: true,
	Short:                 "Lists aliases and target commands",
	Long: `Lists all aliases under optional partial alias name.
To match partial path name, use '-p|--prefix'.

Supports 3 output formats, selectable with '-o|--output':
- Simple (default): <alias> -- <command>
- Command: na add <alias> -- <command>
- YAML: Same as configuration format.

The short forms of 'list' are 'l' and 'ls'.  
By default, all configuration files are merged for this command.`,
	Aliases:           []string{"l", "ls"},
	ValidArgsFunction: validListArgs,
	Run:               listRun,
}

var flagOutputFormat = cli.OutputFormatSimple

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
		aliases = config.ListAliases(args...)
	}
	if out, err := cli.OutputFormatToFunc[flagOutputFormat](aliases); err != nil {
		fmt.Println("na: list:", err)
		os.Exit(1)
	} else {
		fmt.Print(out)
	}
}

func init() {
	listCmd.Flags().BoolP("prefix", "p", false, "List aliases by prefix instead of exact match")
	listCmd.Flags().VarP(&flagOutputFormat, "output", "o", "Output format (simple|cmd|yaml)")
	listCmd.RegisterFlagCompletionFunc("output", func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		formats := make([]string, 0)
		for _, f := range cli.AllOutputFormats() {
			formats = append(formats, string(f))
		}
		return formats, cobra.ShellCompDirectiveDefault
	})
}
