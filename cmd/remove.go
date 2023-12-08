package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/moyiz/na/internal/config"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:               "remove [name ...]",
	Short:             "Removes a nested alias from configuration",
	Aliases:           []string{"rm"},
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: validRemoveArgs,
	Run: func(cmd *cobra.Command, args []string) {
		c := config.GetFromFiles(CurrentConfigFile())
		if err := c.UnsetAlias(args...); err != nil {
			fmt.Println("na:", strings.Join(args, " ")+":", err)
			os.Exit(1)
		}
	},
}

func validRemoveArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	c := config.GetFromFiles(CurrentConfigFile())
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
