package selector_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goropikari/code-strength/internal/selector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelectSearchesAndTogglesMultipleDirectories(t *testing.T) {
	got, err := selector.Select(strings.NewReader("services\n1,2\ndone\n"), &strings.Builder{}, []string{"services/api", "services/web", "docs"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 || got[0] != "services/api" || got[1] != "services/web" {
		t.Fatalf("unexpected selection: %#v", got)
	}
}

func TestSelectParentSelectsDescendants(t *testing.T) {
	input := strings.NewReader("services\n1\ndone\n")
	dirs := []string{"services", "services/api", "services/web", "docs"}

	got, err := selector.Select(input, &strings.Builder{}, dirs, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 3 || got[0] != "services" || got[1] != "services/api" || got[2] != "services/web" {
		t.Fatalf("unexpected recursive selection: %#v", got)
	}
}

func TestSelectUsesLineModeForNonTerminalFiles(t *testing.T) {
	t.Run("regular files do not enter interactive mode", func(t *testing.T) {
		// Arrange
		inputPath := filepath.Join(t.TempDir(), "input")
		input, err := os.Create(inputPath)
		require.NoError(t, err)
		_, err = input.WriteString("done\n")
		require.NoError(t, err)
		_, err = input.Seek(0, 0)
		require.NoError(t, err)
		output, err := os.Create(filepath.Join(t.TempDir(), "output"))
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, input.Close())
			require.NoError(t, output.Close())
		})

		// Act
		got, err := selector.Select(input, output, []string{"docs"}, nil)

		// Assert
		require.NoError(t, err)
		assert.Empty(t, got)
	})
}

func TestSelectIgnoresInvalidIndexes(t *testing.T) {
	t.Run("invalid indexes do not panic or select a directory", func(t *testing.T) {
		// Arrange
		input := strings.NewReader("\nbad,0,4,2\ndone\n")
		dirs := []string{"docs", "services"}

		// Act
		got, err := selector.Select(input, &strings.Builder{}, dirs, nil)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, []string{"services"}, got)
	})
}
