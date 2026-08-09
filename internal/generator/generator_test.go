package generator_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goropikari/code-strength/internal/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildEntriesMarksSelectedSubtreeProduction(t *testing.T) {
	t.Run("selected directory and descendants are production", func(t *testing.T) {
		// Arrange
		dirs := []string{".", "services", "services/api", "docs"}

		// Act
		got := generator.BuildEntries(dirs, []string{" services/ "})

		// Assert
		assert.Equal(t, []generator.Directory{
			{
				Path:  ".",
				Level: "development",
			},
			{
				Path:  "services",
				Level: "production",
			},
			{
				Path:  "services/api",
				Level: "production",
			},
			{
				Path:  "docs",
				Level: "development",
			},
		}, got)
	})

	t.Run("dot selects every directory", func(t *testing.T) {
		// Arrange
		dirs := []string{"docs", "services"}

		// Act
		got := generator.BuildEntries(dirs, []string{"."})

		// Assert
		assert.Equal(t, []generator.Directory{
			{
				Path:  "docs",
				Level: "production",
			},
			{
				Path:  "services",
				Level: "production",
			},
		}, got)
	})
}

func TestWriteRegeneratesDefinition(t *testing.T) {
	t.Run("writes YAML into a new parent directory", func(t *testing.T) {
		// Arrange
		path := filepath.Join(t.TempDir(), "nested", "requirements.yml")

		// Act
		require.NoError(t, generator.Write(path, []generator.Directory{{
			Path:  ".",
			Level: "development",
		}}, []string{".git"}))
		data, err := os.ReadFile(path)

		// Assert
		require.NoError(t, err)

		text := string(data)
		assert.Contains(t, text, "path: .")
		assert.Contains(t, text, "level: development")
		assert.Contains(t, text, "- .git")
	})
}
