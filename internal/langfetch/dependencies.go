package langfetch

import (
	"os"
	"strings"
)

// depMap defines inter-item dependencies: if item X requires Y,Z before it.
var depMap = map[string][]string{
	// Web stacks
	"lamp":  {"apache", "mysql", "php"},
	"lemp":  {"nginx", "mysql", "php"},
	"lapp":  {"apache", "postgresql", "php"},
	"lnmp":  {"nginx", "mysql", "php"},
	"mean":  {"mongodb", "node", "npm", "angular"},
	"mern":  {"mongodb", "node", "npm", "react"},
	"mevn":  {"mongodb", "node", "npm", "vue"},
	"pern":  {"postgresql", "node", "npm", "react"},
	"elk":   {"elasticsearch", "logstash", "kibana"},
	"efk":   {"elasticsearch", "fluentd", "kibana"},
	"jam":   {"node", "npm", "gatsby"},
	"wamp":  {"apache", "mysql", "php"},
	"xampp": {"apache", "mariadb", "php", "perl"},

	// New modern stacks
	"bun":         {"bun", "elysia", "sqlite"},
	"deno":        {"deno", "fresh"},
	"next-trpc":   {"node", "next", "trpc", "typescript", "tailwindcss"},
	"phoenix":     {"elixir", "erlang", "phoenix", "postgresql"},
	"tauri":       {"rust", "node", "webkitgtk"},
	"tauri-solid": {"tauri", "solid-js", "typescript"},

	// Other stacks
	"tall":   {"tailwindcss", "alpinejs", "laravel", "livewire"},
	"smack":  {"spark", "mesos", "akka", "cassandra", "kafka"},
	"tick":   {"telegraf", "influxdb", "chronograf", "kapacitor"},
	"famp":   {"apachectl", "mysql", "php"},
	"slad":   {"aws", "sam"},
	"kstack": {"kubectl", "kustomize", "helm", "k3d"},
	"dcp":    {"docker", "compose", "podman"},

	// Language toolchains
	"golang": {"go"},
	"python": {"python", "pip"},
	"rust":   {"rustc", "cargo"},
	"ruby":   {"ruby", "gem", "bundler"},
	"java":   {"java", "javac", "maven"},
	"dotnet": {"dotnet"},
	"dart":   {"dart", "flutter"},
	"swift":  {"swift", "swiftc"},
	"r":      {"r", "rscript"},
	"node":   {"node", "npm"},
	"php":    {"php", "composer"},

	// Database components
	"mongodb":    {"mongod"},
	"mysql":      {"mysqld"},
	"postgresql": {"postgres"},
	"sqlite":     {"sqlite3"},
	"cassandra":  {"cassandra"},
	"kafka":      {"kafka-server-start"},

	// Web servers
	"apache": {"apache2", "httpd"},
	"nginx":  {"nginx"},
	"caddy":  {"caddy"},
}

// Dependencies returns the list of items a given key depends on.
// It returns nil if there are no registered dependencies.
func Dependencies(key string) []string {
	key = strings.ToLower(key)
	if deps, ok := depMap[key]; ok {
		return deps
	}
	return nil
}

// hasDependency checks if a package.json has a specific dependency
func hasDependency(packageJsonPath, dep string) bool {
	content, err := os.ReadFile(packageJsonPath)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), `"`+dep+`"`)
}

// fileContains checks if a file contains specific text
func fileContains(filePath, text string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), text)
}
