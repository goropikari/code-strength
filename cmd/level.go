package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goropikari/code-strength/internal/generator"
	"github.com/spf13/cobra"
)

//nolint:cyclop // Command setup keeps path validation and output behavior together.
func newLevelCommand() *cobra.Command {
	var root, output string

	cmd := &cobra.Command{
		Use:     "level PATH",
		Aliases: []string{"get"},
		Short:   "Print the requirement level for a directory or file",
		Example: `  # Query a path relative to the repository root.
  code-strength level --root . internal/scanner/scanner.go

  # Query an absolute path.
  code-strength level --root . /work/project/cmd/main.go

  # Print "unknown" when no existing path or definition matches.
  code-strength level --root . does/not/exist`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			absRoot, err := filepath.Abs(root)
			if err != nil {
				return err
			}

			configPath := output
			if !filepath.IsAbs(configPath) {
				configPath = filepath.Join(absRoot, configPath)
			}

			definition, err := generator.Read(configPath)
			if err != nil {
				return err
			}

			target := args[0]
			if !filepath.IsAbs(target) {
				target = filepath.Join(absRoot, target)
			}

			target = filepath.Clean(target)

			if _, err := os.Stat(target); err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), "unknown")

				return nil
			}

			relative, err := filepath.Rel(absRoot, target)
			if err != nil {
				return err
			}

			if relative == ".." || len(relative) > 3 && relative[:3] == ".."+string(filepath.Separator) {
				fmt.Fprintln(cmd.OutOrStdout(), "unknown")

				return nil
			}

			level, ok := generator.LevelForPath(definition, filepath.ToSlash(relative))
			if !ok {
				fmt.Fprintln(cmd.OutOrStdout(), "unknown")

				return nil
			}

			fmt.Fprintln(cmd.OutOrStdout(), level)

			return nil
		},
	}
	cmd.Flags().StringVarP(&root, "root", "r", ".", "repository root")
	cmd.Flags().StringVarP(&output, "output", "o", ".ai-requirements.yml", "requirement definition YAML path")

	return cmd
}
