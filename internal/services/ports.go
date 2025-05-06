package services

import (
    "fmt"
    "net"
    "time"
)

// DefaultPorts maps service names (or stack dependency keys) to their common TCP ports.
var DefaultPorts = map[string][]int{
    "apache":      {80, 443},
    "nginx":       {80, 443, 8080},
    "mysql":       {3306},
    "postgresql":  {5432},
    "mongodb":     {27017},
    "redis":       {6379},
    "docker":      {2375, 2376},
    "ssh":         {22},
    "phpfpm":      {9000},
    "memcached":   {11211},
    "elasticsearch": {9200, 9300},
    "logstash":    {5044, 9600},
    "kibana":      {5601},
    "fluentd":     {24224},
    "telegraf":    {8186},
    "influxdb":    {8086},
    "chronograf":  {8888},
    "kapacitor":   {9092},
    // TODO: will add more later
}

// PortStatus represents the open/closed state of a single port for a service.
type PortStatus struct {
    Service string
    Port    int
    Open    bool
}

// CheckPorts probes each port for the given services on localhost.
// It returns a slice of PortStatus entries for every configured port.
// The timeout controls the DialTimeout duration per port.
func CheckPorts(svcs []string, timeout time.Duration) []PortStatus {
    var out []PortStatus
    for _, svc := range svcs {
        ports, ok := DefaultPorts[svc]
        if !ok {
            continue
        }
        for _, p := range ports {
            addr := fmt.Sprintf("127.0.0.1:%d", p)
            conn, err := net.DialTimeout("tcp", addr, timeout)
            open := err == nil
            if conn != nil {
                conn.Close()
            }
            out = append(out, PortStatus{
                Service: svc,
                Port:    p,
                Open:    open,
            })
        }
    }
    return out
}
