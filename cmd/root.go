package cmd

import (
	"errors"
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/moyiz/na/internal/config"
	"github.com/moyiz/na/internal/consts"
	"github.com/spf13/cobra"
)

var Version string

var rootCmd = &cobra.Command{
	Use:               "na",
	Short:             "CLI tool to effortlessly manage context aware nested shortcuts for shell commands.",
	Long:              consts.Logo,
	Run:               listRun,
	Version:           Version,
	DisableAutoGenTag: true,
}

func getConfigFromArgs() (string, error) {
	if local, _ := rootCmd.Flags().GetBool("local"); local {
		return ".na.yaml", nil
	} else if config, _ := rootCmd.Flags().GetString("config"); config != "" {
		return config, nil
	} else {
		return "", errors.New("no config set")
	}
}

// Returns the active configuration file path according to flags.
// Defaults to $XDG_CONFIG_HOME/na/na.yaml
func ActiveConfigFile() string {
	if configFile, err := getConfigFromArgs(); err == nil {
		return configFile
	} else {
		return path.Join(xdg.ConfigHome, "na", "na.yaml")
	}
}

// Returns a list of all configs if none was set via args.
func AllConfigFiles() []string {
	if configFile, err := getConfigFromArgs(); err == nil {
		return []string{configFile}
	} else {
		return []string{
			path.Join(xdg.ConfigDirs[0], "na", "na.yaml"),
			path.Join(xdg.ConfigHome, "na", "na.yaml"),
			".na.yaml",
		}
	}
}

func init() {
	cobra.EnableCommandSorting = true
	cobra.OnInitialize(func() { config.LoadFiles(AllConfigFiles()...) })
	rootCmd.PersistentFlags().StringP("config", "c", "", "Path of the config file to use")
	rootCmd.PersistentFlags().BoolP("local", "l", false, "Use local config (.na.yaml)")
	rootCmd.MarkFlagsMutuallyExclusive("local", "config")
	rootCmd.AddCommand(runCmd, addCmd, removeCmd, listCmd, docsCmd)
}

// Quiet feature (TBA):
// Basename aware. Any basename other than `na` or `go` will be
// treated as `na run BASENAME`.
func Execute() {
	basename := path.Base(os.Args[0])
	switch basename {
	case "na", "go":
		if err := rootCmd.Execute(); err != nil {
			os.Exit(1)
		}
	default:
		args := make([]string, 0)
		args = append(args, "run", basename)
		args = append(args, os.Args[1:]...)
		rootCmd.SetArgs(args)
		rootCmd.Execute()
	}
}
