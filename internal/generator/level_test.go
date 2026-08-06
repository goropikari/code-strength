package generator_test

import (
	"testing"

	"github.com/goropikari/code-strength/internal/generator"
)

func TestLevelForPathUsesNearestParent(t *testing.T) {
	definition := generator.Definition{Directories: []generator.Directory{
		{Path: ".", Level: "development"},
		{Path: "services", Level: "production"},
		{Path: "services/dev", Level: "development"},
	}}

	level, ok := generator.LevelForPath(definition, "services/api/handler.go")
	if !ok || level != "production" {
		t.Fatalf("unexpected level: %q, %v", level, ok)
	}

	level, ok = generator.LevelForPath(definition, "services/dev/main.go")
	if !ok || level != "development" {
		t.Fatalf("unexpected nested level: %q, %v", level, ok)
	}
}

func TestLevelForPathReturnsUnknownForExcludedPath(t *testing.T) {
	definition := generator.Definition{
		Directories: []generator.Directory{{Path: ".", Level: "development"}},
		Exclude:     []string{"node_modules"},
	}

	level, ok := generator.LevelForPath(definition, "node_modules/example/index.js")
	if ok || level != "" {
		t.Fatalf("expected unknown for excluded path, got %q, %v", level, ok)
	}
}
