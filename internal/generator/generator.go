package generator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Directory struct {
	Path  string `yaml:"path"`
	Level string `yaml:"level"`
}

type Definition struct {
	Directories []Directory `yaml:"directories"`
	Exclude     []string    `yaml:"exclude"`
}

func BuildEntries(dirs, production []string) []Directory {
	selected := append([]string{}, production...)
	sort.Strings(selected)

	result := make([]Directory, 0, len(dirs))
	for _, path := range dirs {
		level := "development"

		for _, target := range selected {
			target = filepath.ToSlash(strings.Trim(strings.TrimSpace(target), "/"))
			if target == "." || path == target || strings.HasPrefix(path, target+"/") {
				level = "production"
				break
			}
		}

		result = append(result, Directory{
			Path:  path,
			Level: level,
		})
	}

	return result
}

func Write(path string, dirs []Directory, excludes []string) error {
	var output bytes.Buffer

	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(2)

	err := encoder.Encode(Definition{
		Directories: dirs,
		Exclude:     excludes,
	})
	closeErr := encoder.Close()

	if err != nil {
		return fmt.Errorf("marshal definition: %w", err)
	}

	if closeErr != nil {
		return fmt.Errorf("close definition: %w", closeErr)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	if err := os.WriteFile(path, output.Bytes(), 0o600); err != nil {
		return fmt.Errorf("write output: %w", err)
	}

	return nil
}

func Read(path string) (Definition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Definition{}, fmt.Errorf("read definition: %w", err)
	}

	var definition Definition
	if err := yaml.Unmarshal(data, &definition); err != nil {
		return Definition{}, fmt.Errorf("parse definition: %w", err)
	}

	for _, directory := range definition.Directories {
		if directory.Level != "development" && directory.Level != "production" {
			return Definition{}, fmt.Errorf("invalid requirement level %q for %q", directory.Level, directory.Path)
		}
	}

	return definition, nil
}

func LevelForPath(definition Definition, target string) (string, bool) {
	target = filepath.ToSlash(filepath.Clean(target))
	if isExcluded(target, definition.Exclude) {
		return "", false
	}

	bestPath := ""
	level := ""

	for _, directory := range definition.Directories {
		path := filepath.ToSlash(filepath.Clean(directory.Path))
		if path != "." && target != path && !strings.HasPrefix(target, path+"/") {
			continue
		}

		if path == "." || path == target || len(path) > len(bestPath) {
			bestPath = path
			level = directory.Level
		}
	}

	return level, level != ""
}

func isExcluded(target string, excludes []string) bool {
	for _, excluded := range excludes {
		excluded = filepath.ToSlash(filepath.Clean(excluded))
		if target == excluded || strings.HasPrefix(target, excluded+"/") {
			return true
		}

		for part := range strings.SplitSeq(target, "/") {
			if part == excluded {
				return true
			}
		}
	}

	return false
}
