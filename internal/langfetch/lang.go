package langfetch

import (
	"errors"
	"os/exec"
	"sort"
	"strings"

	"github.com/sanix-darker/stackfetch/internal/cache"
)

type LangInfo struct {
	Name    string   `json:"name"`
	Version string   `json:"version,omitempty"`
	Details []string `json:"details,omitempty"`
}

func (li LangInfo) String() string {
	s := li.Name
	if li.Version != "" {
		s += " Version: " + li.Version
	}
	if len(li.Details) > 0 {
		s += "\n" + strings.Join(li.Details, "\n")
	}
	return s
}

// executor abstraction for test stubbing
var execCmd = exec.Command

// ExecRunner can be swapped at runtime (e.g. container mode).
// Default implementation = cached local exec.  Override from other packages.
var ExecRunner = func(bin string, args ...string) ([]byte, error) {
	return cache.Run(bin, args...)
}

type fetchFn func() (LangInfo, error)

// Registry of language names -> fetcher
var registry = map[string]fetchFn{}

// Register called in each lang file's init
func register(key string, fn fetchFn) { registry[strings.ToLower(key)] = fn }

// Fetch entrypoint
func Fetch(lang string) (LangInfo, error) {
	key := strings.ToLower(lang)
	if fn, ok := registry[key]; ok {
		return fn()
	}
	msg := ">> Unsupported language/stack !\n\n>> Please create a feature request on https://github.com/sanix-darker/stackfetch to add it): " + lang
	return LangInfo{}, errors.New(msg)
}

// Keys returns the list of registered fetch keys in deterministic order.
func Keys() []string {
	keys := make([]string, 0, len(registry))
	for key := range registry {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// IsSupported reports whether a language/stack key is registered.
func IsSupported(key string) bool {
	_, ok := registry[strings.ToLower(strings.TrimSpace(key))]
	return ok
}

// SuggestKeys returns up to max similar registered keys for query.
func SuggestKeys(query string, max int) []string {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" || max <= 0 {
		return nil
	}
	candidates := Keys()
	if len(candidates) == 0 {
		return nil
	}
	if max > len(candidates) {
		max = len(candidates)
	}

	type suggestion struct {
		key   string
		score int
	}
	out := make([]suggestion, 0, len(candidates))

	for _, key := range candidates {
		if strings.EqualFold(key, query) {
			out = append(out, suggestion{key: key, score: 0})
			continue
		}
		if strings.Contains(key, query) {
			out = append(out, suggestion{key: key, score: 1})
			continue
		}
		if strings.HasPrefix(key, query) {
			out = append(out, suggestion{key: key, score: 2})
			continue
		}
		distance := levenshteinDistance(key, query)
		if distance <= 3 {
			out = append(out, suggestion{key: key, score: distance + 2})
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].score == out[j].score {
			return out[i].key < out[j].key
		}
		return out[i].score < out[j].score
	})

	if len(out) > max {
		out = out[:max]
	}

	result := make([]string, 0, len(out))
	for _, suggestion := range out {
		result = append(result, suggestion.key)
	}
	return result
}

func levenshteinDistance(a, b string) int {
	if a == b {
		return 0
	}
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	aRunes := []rune(a)
	bRunes := []rune(b)
	rows := len(aRunes) + 1
	cols := len(bRunes) + 1

	prev := make([]int, cols)
	curr := make([]int, cols)
	for j := 0; j < cols; j++ {
		prev[j] = j
	}
	for i := 1; i < rows; i++ {
		curr[0] = i
		for j := 1; j < cols; j++ {
			cost := 0
			if aRunes[i-1] != bRunes[j-1] {
				cost = 1
			}
			insert := prev[j] + 1
			delete := curr[j-1] + 1
			substitute := prev[j-1] + cost

			curr[j] = min(insert, min(delete, substitute))
		}
		prev, curr = curr, prev
	}
	return prev[cols-1]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper to grab first line of command output
func runFirst(cmd string, args ...string) (string, error) {
	out, err := ExecRunner(cmd, args...)
	if err != nil {
		return "", err
	}
	s := strings.TrimSpace(string(out))
	if i := strings.IndexByte(s, '\n'); i > 0 {
		s = s[:i]
	}
	return s, nil
}

// Helper to run command; errors ignored.
func runSilently(cmd string, args ...string) string {
	out, _ := execCmd(cmd, args...).Output()
	return strings.TrimSpace(string(out))
}
