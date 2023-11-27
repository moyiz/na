package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func validRemoveArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	viper.SetConfigFile(CurrentConfigFile())
	viper.ReadInConfig()
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

func removeRun(cmd *cobra.Command, args []string) {
	viper.SetConfigFile(CurrentConfigFile())
	viper.ReadInConfig()

	settings := viper.AllSettings()
	numArgs := len(args)
	configPart := settings
	for _, part := range args[:numArgs-1] {
		configPart = configPart[part].(map[string]interface{})
	}
	if _, ok := configPart[args[numArgs-1]]; !ok {
		fmt.Println("na:", strings.Join(args, " ")+": not found")
		return
	}

	delete(configPart, args[numArgs-1])
	encodedConfig, _ := json.MarshalIndent(settings, "", " ")
	viper.ReadConfig(bytes.NewReader(encodedConfig))

	viper.WriteConfig()
}
