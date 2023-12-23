package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/moyiz/na/internal/cli"
	"github.com/moyiz/na/internal/config"
	"github.com/moyiz/na/internal/utils"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:                   "run [name ...] [--] [args ...]",
	DisableFlagsInUseLine: true,
	Short:                 "Runs a nested alias",
	Long: `Runs a nested alias.
All arguments after an optional double-dash (--) will be passed as arguments
to the target command.  
The target command will be executed in a new sub-shell. It will try to
determine the shell from which it was executed. If the current shell cannot be
determined, it will fallback to 'sh'.

The short form of 'run' is 'r'.  
By default, all configuration files are merged for this command.`,
	Aliases:           []string{"r"},
	ValidArgsFunction: validRunArgs,
	Run:               runRun,
}

func validRunArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if slices.Contains(os.Args, "--") {
		// Potential location to auto complete commands
		return []string{}, cobra.ShellCompDirectiveDefault
	}
	config.LoadFiles(AllConfigFiles()...)
	return cli.ListNextParts(config.ListAliases(), args), cobra.ShellCompDirectiveNoFileComp
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

	config.LoadFiles(AllConfigFiles()...)
	if alias, err := config.GetAlias(aliasParts...); err != nil {
		fmt.Println("na:", strings.Join(aliasParts, " ")+":", err)
	} else {
		command := utils.GenerateCommand(alias.Command, extraArgs)
		utils.RunInCurrentShell(command.Command, command.Args)
	}
}
