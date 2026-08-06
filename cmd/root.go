package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "code-strength",
	Short: "Generate directory-specific AI requirement levels",
	Long:  "Generate and query directory-specific AI requirement levels for a repository.",
}

func init() {
	rootCmd.AddCommand(newGenerateCommand())
	rootCmd.AddCommand(newLevelCommand())
}

// Execute runs the CLI and exits with a non-zero status when a command fails.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
