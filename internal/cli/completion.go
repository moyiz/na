package cli

import (
	"slices"
	"strings"

	"github.com/moyiz/na/internal/config"
)

// Given a list of alias name parts, return a list of valid next parts.
// Example:
//
//	my:
//	  aliases:
//	    one: cmd1
//	    two: cmd2
//
// ListNextParts([]string{"my"}) -> []string{"aliases"}
// ListNextParts([]string{"my", "aliases"}) -> []string{"cmd1", "cmd2"}
func ListNextParts(aliases []config.Alias, parts []string) []string {
	currentPrefix := strings.Join(parts, " ")
	suggestions := make([]string, 0)
	for _, a := range aliases {
		trail, found := strings.CutPrefix(a.Name, currentPrefix)
		if tf := strings.Fields(trail); found && len(tf) > 0 && !slices.Contains(suggestions, tf[0]) {
			suggestions = append(suggestions, tf[0])
		}
	}

	return suggestions
}
