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

var serviceAliases = map[string][]string{
	"apache":        {"apache2", "apachectl", "httpd"},
	"mysql":         {"mysql", "mysqld"},
	"mariadb":       {"mariadb", "mariadbd"},
	"postgresql":    {"postgresql", "postgres"},
	"mongodb":       {"mongod", "mongos", "mongodb"},
	"cassandra":     {"cassandra", "cassandra-server"},
	"kafka":         {"kafka", "kafka-server-start"},
	"rabbitmq":      {"rabbitmq", "rabbitmq-server", "rabbitmqctl"},
	"elasticsearch": {"elasticsearch", "elasticsearch-8"},
	"kubernetes":    {"kubernetes", "kubelet", "kube-apiserver", "kubectl"},
	"nginx":         {"nginx"},
	"redis":         {"redis-server"},
	"java":          {"java"},
	"node":          {"node"},
	"php":           {"php"},
	"docker":        {"docker", "dockerd"},
	"postgres":      {"postgres"},
	"k3d":           {"k3d"},
	"helm":          {"helm"},
	"kubectl":       {"kubectl"},
}

func aliasesFor(name string) []string {
	if aliases, ok := serviceAliases[name]; ok {
		return aliases
	}
	return []string{name}
}

// StatusByName returns the Status for a single service/tool.
func StatusByName(name string) Status {
	s := Status{Name: name}
	candidates := aliasesFor(name)

	// Check installation
	for _, candidate := range candidates {
		if _, err := exec.LookPath(candidate); err == nil {
			s.Installed = true
			break
		}
	}

	// Check running state
	switch runtime.GOOS {
	case "linux", "darwin":
		// pgrep -x matches exact process name
		for _, candidate := range candidates {
			if err := exec.Command("pgrep", "-x", candidate).Run(); err == nil {
				s.Running = true
				break
			}
		}
	case "windows":
		// tasklist /FI filter
		for _, candidate := range candidates {
			out, err := exec.Command("tasklist", "/FI", "IMAGENAME eq "+candidate+".exe").Output()
			if err == nil && strings.Contains(string(out), candidate+".exe") {
				s.Running = true
				break
			}
		}
	}

	return s
}
