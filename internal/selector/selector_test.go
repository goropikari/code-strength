package selector_test

import (
	"strings"
	"testing"

	"github.com/goropikari/code-strength/internal/selector"
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
