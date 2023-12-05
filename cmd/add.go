package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func addRun(cmd *cobra.Command, args []string) {
	configFile := CurrentConfigFile()
	configDir := path.Dir(configFile)
	if _, err := os.Stat(configDir); err != nil {
		os.MkdirAll(configDir, 0775)
	}

	viper.SetConfigFile(configFile)
	viper.ReadInConfig()

	commandStart := cmd.ArgsLenAtDash()
	if commandStart == 0 {
		fmt.Println("na: '--' cannot be the first argument.")
		os.Exit(1)
	}
	if commandStart < 0 {
		commandStart = len(args) - 1
	}
	viper.Set(strings.Join(args[:commandStart], "."), strings.Join(args[commandStart:], " "))

	if err := viper.WriteConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
