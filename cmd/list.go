package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/moyiz/na/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:               "list [name ...]",
	Short:             "Shows a list of aliases",
	Long:              "Show a list of aliases under an optional given prefix.",
	Aliases:           []string{"l", "ls"},
	ValidArgsFunction: validListArgs,
	Run:               listRun,
}

func validListArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if slices.Contains(os.Args, "--") {
		// Potential location to auto complete commands
		return []string{}, cobra.ShellCompDirectiveDefault
	}

	c := config.GetFromFiles(AllConfigFiles()...)
	currentPrefix := strings.Join(args, " ")
	suggestions := make([]string, 0)
	for _, a := range c.ListAliases(args...) {
		trail, _ := strings.CutPrefix(a.Name, currentPrefix)
		if trailFields := strings.Fields(trail); len(trailFields) > 0 {
			suggestions = append(suggestions, trailFields[0])
		} else {
			break
		}
	}
	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

func listRun(cmd *cobra.Command, args []string) {
	for _, alias := range config.ListAliases(args...) {
		fmt.Println(alias.Name, "--", alias.Command)
	}
}
