package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/goropikari/code-strength/internal/generator"
	"github.com/goropikari/code-strength/internal/scanner"
	"github.com/goropikari/code-strength/internal/selector"
	"github.com/spf13/cobra"
)

type stringFlags []string

func (s *stringFlags) String() string { return fmt.Sprint([]string(*s)) }
func (s *stringFlags) Set(value string) error {
	*s = append(*s, value)
	return nil
}
func (s *stringFlags) Type() string { return "string" }

func newGenerateCommand() *cobra.Command {
	var (
		root, output             string
		production, extraExclude stringFlags
		nonInteractive           bool
	)

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Select production directories and regenerate the YAML definition",
		Example: `  # Interactive search and multi-select.
  code-strength generate --root .

  # Generate without the terminal UI.
  code-strength generate --root . --non-interactive --production services/api

  # Select multiple production directories and add an exclusion.
  code-strength generate --root . --production cmd --production internal --exclude generated`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			absRoot, err := filepath.Abs(root)
			if err != nil {
				return err
			}

			dirs, excludes, err := scanner.List(absRoot, []string(extraExclude))
			if err != nil {
				return err
			}

			selected := []string(production)
			if !nonInteractive {
				selected, err = selector.Select(cmd.InOrStdin(), cmd.OutOrStdout(), dirs, selected)
				if err != nil {
					return err
				}
			}

			entries := generator.BuildEntries(dirs, selected)

			out := output
			if !filepath.IsAbs(out) {
				out = filepath.Join(absRoot, out)
			}

			if err := generator.Write(out, entries, excludes); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "generated %s (%d directories)\n", out, len(entries))

			return nil
		},
	}
	cmd.Flags().StringVarP(&root, "root", "r", ".", "repository root")
	cmd.Flags().StringVarP(&output, "output", "o", ".ai-requirements.yml", "output YAML path, relative to the repository root")
	cmd.Flags().Var(&production, "production", "production directory (repeatable; implies --non-interactive when supplied)")
	cmd.Flags().Var(&extraExclude, "exclude", "additional directory name or path to exclude (repeatable)")
	cmd.Flags().BoolVar(&nonInteractive, "non-interactive", false, "do not open the selector; use --production values")
	cmd.PreRun = func(cmd *cobra.Command, _ []string) {
		if cmd.Flags().Changed("production") {
			nonInteractive = true
		}
	}

	return cmd
}
