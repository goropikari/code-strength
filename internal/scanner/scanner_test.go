package scanner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goropikari/code-strength/internal/scanner"
)

func TestListIncludesRootAndSkipsDefaultsAndExtras(t *testing.T) {
	root := t.TempDir()
	for _, path := range []string{"services/api", ".git/hooks", "node_modules/pkg", "generated/out"} {
		if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(path)), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	dirs, excludes, err := scanner.List(root, []string{"generated"})
	if err != nil {
		t.Fatal(err)
	}

	joined := map[string]bool{}
	for _, dir := range dirs {
		joined[dir] = true
	}

	for _, want := range []string{".", "services", "services/api"} {
		if !joined[want] {
			t.Errorf("missing %q", want)
		}
	}

	for _, unwanted := range []string{".git", ".git/hooks", "node_modules", "generated"} {
		if joined[unwanted] {
			t.Errorf("unexpected %q", unwanted)
		}
	}

	if len(excludes) == 0 {
		t.Fatal("expected excludes")
	}
}

func TestListFollowsDirectorySymlinks(t *testing.T) {
	root := t.TempDir()
	target := t.TempDir()

	if err := os.Mkdir(filepath.Join(target, "windows-dir"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.Symlink(target, filepath.Join(root, "linked")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	dirs, _, err := scanner.List(root, nil)
	if err != nil {
		t.Fatal(err)
	}

	joined := map[string]bool{}
	for _, dir := range dirs {
		joined[dir] = true
	}

	if !joined["linked"] || !joined["linked/windows-dir"] {
		t.Fatalf("expected symlink target directories, got %#v", dirs)
	}
}
