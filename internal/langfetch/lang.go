package langfetch

import (
    "errors"
    "os/exec"
    "strings"
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

// Helper to grab first line of command output
func runFirst(cmd string, args ...string) (string, error) {
    out, err := execCmd(cmd, args...).Output()
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
