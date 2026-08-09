package scanner_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/goropikari/code-strength/internal/scanner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListIncludesRootAndSkipsDefaultsAndExtras(t *testing.T) {
	root := t.TempDir()
	for _, path := range []string{"services/api", "services/web", ".git/hooks", "node_modules/pkg", "generated/out", "reports/archive/deep"} {
		require.NoError(t, os.MkdirAll(filepath.Join(root, filepath.FromSlash(path)), 0o755))
	}

	dirs, excludes, err := scanner.List(root, []string{" generated ", "/generated/", "reports/archive", " "})
	require.NoError(t, err)

	joined := map[string]bool{}
	for _, dir := range dirs {
		joined[dir] = true
	}

	assert.True(t, joined["."])
	assert.True(t, joined["services"])
	assert.True(t, joined["services/api"])
	assert.True(t, joined["services/web"])
	assert.False(t, joined[".git"])
	assert.False(t, joined[".git/hooks"])
	assert.False(t, joined["node_modules"])
	assert.False(t, joined["generated"])
	assert.False(t, joined["generated/out"])
	assert.False(t, joined["reports/archive"])
	assert.False(t, joined["reports/archive/deep"])

	assert.NotContains(t, excludes, "")
	assert.Equal(t, 1, countString(excludes, "generated"))

	if len(excludes) == 0 {
		t.Fatal("expected excludes")
	}
}

func countString(values []string, want string) int {
	count := 0

	for _, value := range values {
		if value == want {
			count++
		}
	}

	return count
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
