package generator_test

import (
	"testing"

	"github.com/goropikari/code-strength/internal/generator"
	"github.com/stretchr/testify/assert"
)

func TestLevelForPathUsesNearestParent(t *testing.T) {
	definition := generator.Definition{Directories: []generator.Directory{
		{
			Path:  ".",
			Level: "development",
		},
		{
			Path:  "services",
			Level: "production",
		},
		{
			Path:  "services/dev",
			Level: "development",
		},
	}}

	t.Run("exact directory uses its level", func(t *testing.T) {
		// Arrange

		// Act
		level, ok := generator.LevelForPath(definition, "services")

		// Assert
		assert.Equal(t, "production", level)
		assert.True(t, ok)
	})

	t.Run("child path uses the nearest parent level", func(t *testing.T) {
		// Arrange

		// Act
		level, ok := generator.LevelForPath(definition, "services/api/handler.go")

		// Assert
		assert.Equal(t, "production", level)
		assert.True(t, ok)
	})

	t.Run("nested child overrides its parent", func(t *testing.T) {
		// Arrange

		// Act
		level, ok := generator.LevelForPath(definition, "services/dev/main.go")

		// Assert
		assert.Equal(t, "development", level)
		assert.True(t, ok)
	})

	t.Run("same-length matching directory wins", func(t *testing.T) {
		// Arrange
		definition := generator.Definition{Directories: []generator.Directory{
			{
				Path:  "bar",
				Level: "production",
			},
			{
				Path:  "bar",
				Level: "development",
			},
		}}

		// Act
		level, ok := generator.LevelForPath(definition, "bar/file.go")

		// Assert
		assert.Equal(t, "production", level)
		assert.True(t, ok)
	})

	t.Run("a non-excluded path remains known", func(t *testing.T) {
		// Arrange

		// Act
		level, ok := generator.LevelForPath(generator.Definition{
			Directories: []generator.Directory{{
				Path:  ".",
				Level: "development",
			}},
			Exclude: []string{"node_modules"},
		}, "src/app")

		// Assert
		assert.Equal(t, "development", level)
		assert.True(t, ok)
	})
}

func TestLevelForPathReturnsUnknownForExcludedPath(t *testing.T) {
	definition := generator.Definition{
		Directories: []generator.Directory{{
			Path:  ".",
			Level: "development",
		}},
		Exclude: []string{"node_modules"},
	}

	t.Run("exact excluded path is unknown", func(t *testing.T) {
		// Arrange

		// Act
		level, ok := generator.LevelForPath(definition, "node_modules")

		// Assert
		assert.Empty(t, level)
		assert.False(t, ok)
	})

	t.Run("excluded child path is unknown", func(t *testing.T) {
		// Arrange

		// Act
		level, ok := generator.LevelForPath(definition, "node_modules/example/index.js")

		// Assert
		assert.Empty(t, level)
		assert.False(t, ok)
	})

	t.Run("excluded nested path segment is unknown", func(t *testing.T) {
		// Arrange

		// Act
		level, ok := generator.LevelForPath(definition, "src/node_modules")

		// Assert
		assert.Empty(t, level)
		assert.False(t, ok)
	})

	t.Run("excluded composite path matches its child", func(t *testing.T) {
		// Arrange
		definition := generator.Definition{
			Directories: []generator.Directory{{
				Path:  ".",
				Level: "development",
			}},
			Exclude: []string{"src/node_modules"},
		}

		// Act
		level, ok := generator.LevelForPath(definition, "src/node_modules/example.js")

		// Assert
		assert.Empty(t, level)
		assert.False(t, ok)
	})
}
