package langfetch

import "strings"

// depMap defines inter-item dependencies: if item X requires Y,Z before it.
var depMap = map[string][]string{
    // Web stacks
    "lamp":  {"apache", "mysql", "php"},
    "lemp":  {"nginx", "mysql", "php"},
    "lapp":  {"apache", "postgresql", "php"},
    "lnmp":  {"nginx", "mysql", "php"},
    "mean":  {"mongodb", "node", "npm", "angular"},
    "mern":  {"mongodb", "node", "npm", "yarn"},
    "mevn":  {"mongodb", "node", "npm", "vue"},
    "pern":  {"postgresql", "node", "npm", "react"},
    "elk":   {"elasticsearch", "logstash", "kibana"},
    "efk":   {"elasticsearch", "fluentd", "kibana"},
    "jam":   {"node", "npm", "gatsby"},
    "wamp":  {"apache", "mysql", "php"},
    "xampp": {"apache", "mariadb", "php", "perl"},

    // Other stacks
    "tall":  {"tailwindcss", "alpinejs", "laravel", "livewire"},
    "smack": {"spark", "mesos", "akka", "cassandra", "kafka"},
    "tick":  {"telegraf", "influxdb", "chronograf", "kapacitor"},
    "famp":  {"apachectl", "mysql", "php"},
    "slad":  {"aws", "sam"},
    "kstack": {"kubectl", "kustomize", "helm", "k3d"},
    "dcp":   {"docker", "compose", "podman"},
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
