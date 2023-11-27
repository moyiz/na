package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func listRun(cmd *cobra.Command, args []string) {
	keys := viper.AllKeys()
	sort.Strings(keys)
	prefix := strings.Join(args, ".")
	for _, k := range keys {
		if strings.HasPrefix(k, prefix) {
			fmt.Println("na add", strings.ReplaceAll(k, ".", " "), "--", viper.GetString(k))
		}
	}
}
