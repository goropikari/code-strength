package scanner

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var DefaultExcludes = []string{".git", "node_modules", "vendor", ".idea", ".vscode", "dist", "build", "tmp"}

// List returns repository-relative directories, including the root as ".".
// An exclude matches either a complete relative path or any directory name.
func List(root string, extra []string) ([]string, []string, error) {
	excludes := append([]string{}, DefaultExcludes...)
	excludes = append(excludes, extra...)
	seen := map[string]bool{}

	cleanExcludes := make([]string, 0, len(excludes))
	for _, item := range excludes {
		item = filepath.ToSlash(strings.Trim(strings.TrimSpace(item), "/"))
		if item == "" || seen[item] {
			continue
		}

		seen[item] = true
		cleanExcludes = append(cleanExcludes, item)
	}

	sort.Strings(cleanExcludes)

	if _, err := os.Stat(root); err != nil {
		return nil, nil, err
	}

	dirs := []string{"."}
	visited := make(map[string]bool)

	err := walkDirectories(root, ".", cleanExcludes, visited, &dirs)
	if err != nil {
		return nil, nil, err
	}

	sort.Strings(dirs)

	return dirs, cleanExcludes, nil
}

func walkDirectories(path, rel string, excludes []string, visited map[string]bool, dirs *[]string) error {
	realPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}

	if visited[realPath] {
		return nil
	}

	visited[realPath] = true

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		childRel := filepath.ToSlash(filepath.Join(rel, entry.Name()))
		if matchesExclude(childRel, excludes) {
			continue
		}

		childPath := filepath.Join(path, entry.Name())

		info, err := os.Stat(childPath)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			continue
		}

		*dirs = append(*dirs, childRel)

		if err := walkDirectories(childPath, childRel, excludes, visited, dirs); err != nil {
			return err
		}
	}

	return nil
}

func matchesExclude(path string, excludes []string) bool {
	for _, item := range excludes {
		if path == item || strings.HasPrefix(path, item+"/") {
			return true
		}

		for part := range strings.SplitSeq(path, "/") {
			if part == item {
				return true
			}
		}
	}

	return false
}
