package services

import (
	"os/exec"
	"runtime"
	"strings"
)

// Status represents the installation and runtime state of a service or CLI tool.
type Status struct {
	Name      string // service or tool name
	Installed bool   // executable found in PATH
	Running   bool   // process/service is currently active
}

// Check inspects only the provided service/tool names.
// It returns a slice of Status, one per requested name.
func Check(names []string) []Status {
	statuses := make([]Status, 0, len(names))
	seen := map[string]struct{}{}

	for _, raw := range names {
		name := strings.ToLower(raw)
		if _, dup := seen[name]; dup {
			continue
		}
		seen[name] = struct{}{}
		statuses = append(statuses, StatusByName(name))
	}

	return statuses
}

// StatusByName returns the Status for a single service/tool.
func StatusByName(name string) Status {
	s := Status{Name: name}

	// Check installation
	if _, err := exec.LookPath(name); err == nil {
		s.Installed = true
	}

	// Check running state
	switch runtime.GOOS {
	case "linux", "darwin":
		// pgrep -x matches exact process name
		if err := exec.Command("pgrep", "-x", name).Run(); err == nil {
			s.Running = true
		}
	case "windows":
		// tasklist /FI filter
		out, err := exec.Command("tasklist", "/FI", "IMAGENAME eq "+name+".exe").Output()
		if err == nil && strings.Contains(string(out), name+".exe") {
			s.Running = true
		}
	}

	return s
}
