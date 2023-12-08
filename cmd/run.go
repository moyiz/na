package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/moyiz/na/internal/config"
	"github.com/moyiz/na/internal/utils"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:               "run [name ...] [--] [args ...]",
	Short:             "Runs a nested alias",
	Long:              "Runs a nested alias.",
	Aliases:           []string{"r"},
	ValidArgsFunction: validRunArgs,
	Run:               runRun,
}

func validRunArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

func runRun(cmd *cobra.Command, args []string) {
	var aliasParts []string
	var extraArgs []string

	if extraArgsStart := cmd.ArgsLenAtDash(); extraArgsStart == 0 {
		fmt.Println("na: '--' cannot be the first argument.")
		os.Exit(1)
	} else if extraArgsStart > 0 {
		extraArgs = args[extraArgsStart:]
		aliasParts = args[:extraArgsStart]
	} else {
		aliasParts = args
	}

	c := config.GetFromFiles(AllConfigFiles()...)
	if alias, err := c.GetAlias(aliasParts...); err != nil {
		fmt.Println("na:", strings.Join(aliasParts, " ")+":", err)
	} else {
		utils.RunInCurrentShell(alias.Command, extraArgs)
	}
}
