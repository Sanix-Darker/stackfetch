package guess

import (
    "os"
    "path/filepath"
    "sort"
    "strings"
)

// MaxDepth defines how deep the directory walker descends from the root.
const MaxDepth = 3

// Guess inspects the project tree up to MaxDepth and infers langfetch keys
// using the filename + extension maps generated in patterns.go. The union of
// matches is returned in stable (sorted) order.
func Guess(root string) []string {
    base := baseNameRules()
    ext := extRules()
    seen := map[string]struct{}{}

    add := func(items []string) {
        for _, it := range items {
            seen[it] = struct{}{}
        }
    }

    _ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return nil // ignore I/O errors for robustness
        }
        if d.IsDir() {
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
        return nil
    })

    out := make([]string, 0, len(seen))
    for k := range seen { out = append(out, k) }
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
