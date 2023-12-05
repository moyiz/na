package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/moyiz/na/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func validRunArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	currentKey := strings.Join(args, ".")
	commands := make([]string, 0)
	containsDoubleDash := slices.Contains(os.Args, "--")
	for _, k := range viper.AllKeys() {
		if strings.HasPrefix(k, currentKey) {
			trail, _ := strings.CutPrefix(k, currentKey)
			if suggestion := strings.Split(strings.TrimLeft(trail, "."), ".")[0]; suggestion != "" && !containsDoubleDash {
				commands = append(commands, suggestion)
			}
		}
	}
	if currentKey != "" && len(commands) == 0 && containsDoubleDash {
		// ToDo: Add completions of the actual command here.
		return commands, cobra.ShellCompDirectiveDefault
	}
	return commands, cobra.ShellCompDirectiveNoFileComp
}

func runRun(cmd *cobra.Command, args []string) {
	var alias string
	var extraArgs []string

	if extraArgsStart := cmd.ArgsLenAtDash(); extraArgsStart == 0 {
		fmt.Println("na: '--' cannot be the first argument.")
		os.Exit(1)
	} else if extraArgsStart > 0 {
		extraArgs = args[extraArgsStart:]
		alias = strings.Join(args[:extraArgsStart], ".")
	} else {
		alias = strings.Join(args, ".")
	}

	if command := viper.GetString(alias); command == "" {
		fmt.Println("na:", strings.ReplaceAll(alias, ".", " ")+": not found")
	} else {
		utils.RunInCurrentShell(command, extraArgs)
	}
}
