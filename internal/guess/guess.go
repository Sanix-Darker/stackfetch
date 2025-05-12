package guess

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// MaxDepth defines how deep the directory walker descends from the root.
const MaxDepth = 4 // Increased from 3 to better detect nested projects

// Guess inspects the project tree up to MaxDepth and infers langfetch keys
// using the filename + extension maps generated in patterns.go. The union of
// matches is returned in stable (sorted) order.
func Guess(root string) []string {
	base := baseNameRules()
	ext := extRules()
	dir := dirRules()
	seen := map[string]struct{}{}

	add := func(items []string) {
		for _, it := range items {
			if it != "" { // Skip empty entries
				seen[it] = struct{}{}
			}
		}
	}

	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // ignore I/O errors for robustness
		}

		relPath, _ := filepath.Rel(root, path)
		if isIgnoredPath(relPath) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			// Check for directory-based patterns (like src-tauri for Tauri)
			if items, ok := dir[d.Name()]; ok {
				add(items)
			}

			if depth(root, path) >= MaxDepth {
				return filepath.SkipDir
			}
			return nil
		}

		name := strings.ToLower(d.Name())
		if items, ok := base[name]; ok {
			add(items)
		}

		if e := strings.ToLower(filepath.Ext(name)); e != "" {
			if items, ok := ext[e]; ok {
				add(items)
			}
		}

		// Check file content patterns for some stacks
		if depth(root, path) <= 2 { // Only check top-level files
			if stack := detectFromFileContent(path); stack != "" {
				add([]string{stack})
			}
		}

		return nil
	})

	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// depth returns the directory depth of path relative to root.
func depth(root, path string) int {
	rel, _ := filepath.Rel(root, path)
	if rel == "." {
		return 0
	}
	return strings.Count(rel, string(os.PathSeparator)) + 1
}

// isIgnoredPath checks if a path should be ignored (node_modules, etc.)
func isIgnoredPath(relPath string) bool {
	parts := strings.Split(relPath, string(os.PathSeparator))
	for _, part := range parts {
		switch part {
		case "node_modules", "vendor", "target", "dist", "build":
			return true
		}
	}
	return false
}

// detectFromFileContent checks specific files for stack signatures
func detectFromFileContent(path string) string {
	name := filepath.Base(path)
	switch name {
	case "deno.json", "deno.lock":
		return "deno"
	case "bun.lockb":
		return "bun"
	case "mix.exs":
		if content, err := os.ReadFile(path); err == nil {
			if strings.Contains(string(content), `:phoenix`) {
				return "phoenix"
			}
		}
	case "tauri.conf.json":
		return "tauri"
	}
	return ""
}
