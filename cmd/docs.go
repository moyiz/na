package cmd

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCmd = &cobra.Command{
	Use:   "docs [-m mode] [-o path]",
	Short: "Generate docs",
	Run:   docsRun,
}

func docsRun(cmd *cobra.Command, args []string) {
	output, _ := cmd.Flags().GetString("output")
	mode, _ := cmd.Flags().GetString("mode")
	if !slices.Contains([]string{"man", "md", "yaml"}, mode) {
		fmt.Println("Invalid mode:", mode)
		os.Exit(1)
	}
	// Create directories if not exist
	if _, err := os.Stat(output); err != nil {
		if err := os.MkdirAll(output, 0o775); err != nil {
			fmt.Println("Failed to create output directory:", err)
			os.Exit(1)
		}
	}
	if err := genDocs(mode, output); err != nil {
		fmt.Println("Failed to generate docs", err)
		os.Exit(1)
	}
}

func genDocs(mode string, output string) error {
	switch mode {
	case "man":
		header := &doc.GenManHeader{
			Title:   "NA",
			Section: "1",
			Source:  "moyiz/na " + Version,
		}
		return doc.GenManTree(rootCmd, header, output)
	case "md":
		return doc.GenMarkdownTree(rootCmd, output)
	case "yaml":
		return doc.GenYamlTree(rootCmd, output)
	default:
		return errors.New("invalid mode")
	}
}

func init() {
	docsCmd.Hidden = true
	docsCmd.Flags().StringP("mode", "m", "man", "Output format (man|md)")
	docsCmd.Flags().StringP("output", "o", "./docs", "Output directory for generated docs")
}
