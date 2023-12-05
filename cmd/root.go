package cmd

import (
	"os"
	"path"

	"github.com/adrg/xdg"
	"github.com/moyiz/na/internal/consts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "na",
	Short: "A simple tool to create nested alias-like commands with dynamic completions",
	Long:  consts.Logo,
	Run:   listRun,
}

var addCmd = &cobra.Command{
	Use:   "add [name ...] [--] [command ...]",
	Short: "Adds a nested alias to configuration",
	Long: `Adds a nested alias to configuration.
If the command consists more than a single word, an optional '--' will act as
a delimiter between the alias and the command.
It will create the config directory if not exists`,
	Aliases: []string{"a"},
	Args:    cobra.MinimumNArgs(2),
	Run:     addRun,
}

var listCmd = &cobra.Command{
	Use:     "list [name ...]",
	Short:   "Shows a list of aliases",
	Long:    "Show a list of aliases under an optional given prefix.",
	Aliases: []string{"l", "ls"},
	Run:     listRun,
}

var runCmd = &cobra.Command{
	Use:               "run [name ...] [--] [args ...]",
	Short:             "Runs a nested alias",
	Long:              "Runs a nested alias.",
	Aliases:           []string{"r"},
	ValidArgsFunction: validRunArgs,
	Run:               runRun,
}

var removeCmd = &cobra.Command{
	Use:               "remove [name ...]",
	Short:             "Removes a nested alias from configuration",
	Aliases:           []string{"rm"},
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: validRemoveArgs,
	Run:               removeRun,
}

// Returns the active configuration file according to flags.
// Defaults to $XDG_CONFIG_HOME/na/config.yaml
func CurrentConfigFile() string {
	if local, _ := rootCmd.Flags().GetBool("local"); local {
		return ".na.yaml"
	} else if config, _ := rootCmd.Flags().GetString("config"); config != "" {
		return config
	}
	return path.Join(xdg.ConfigHome, "na", "na.yaml")
}

func init() {
	cobra.EnableCommandSorting = true
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("config", "c", "", "Path of the config file to use")
	rootCmd.PersistentFlags().BoolP("local", "l", false, "Use local config (.na.yaml)")
	rootCmd.MarkFlagsMutuallyExclusive("local", "config")
	rootCmd.AddCommand(runCmd, addCmd, removeCmd, listCmd)
}

func initConfig() {
	if configFile, _ := rootCmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(path.Clean(configFile))
		viper.ReadInConfig()
	} else if local, _ := rootCmd.Flags().GetBool("local"); local {
		viper.SetConfigFile(".na.yaml")
		viper.ReadInConfig()
	} else {
		// Select and merge configs according to this order:
		// .na.yaml
		// ~/.config/na/na.yaml ($XDG_CONFIG_HOME/na/config.yaml)
		// /etc/xdg/na/na.yaml ($XDG_CONFIG_DIRS/na/config.yaml)
		viper.SetConfigFile(path.Join(xdg.ConfigDirs[0], "na", "na.yaml"))
		viper.ReadInConfig()
		viper.SetConfigFile(path.Join(xdg.ConfigHome, "na", "na.yaml"))
		viper.MergeInConfig()
		viper.SetConfigFile(".na.yaml")
		viper.MergeInConfig()
	}
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
