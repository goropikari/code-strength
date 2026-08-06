package generator_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goropikari/code-strength/internal/generator"
)

func TestBuildEntriesMarksSelectedSubtreeProduction(t *testing.T) {
	got := generator.BuildEntries([]string{".", "services", "services/api", "docs"}, []string{"services"})
	if got[0].Level != "development" || got[2].Level != "production" || got[3].Level != "development" {
		t.Fatalf("unexpected entries: %#v", got)
	}
}

func TestWriteRegeneratesDefinition(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "requirements.yml")
	if err := generator.Write(path, []generator.Directory{{Path: ".", Level: "development"}}, []string{".git"}); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	text := string(data)
	if !strings.Contains(text, "path: .") || !strings.Contains(text, "level: development") || !strings.Contains(text, "- .git") {
		t.Fatalf("unexpected YAML: %s", text)
	}
}
