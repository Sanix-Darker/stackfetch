package sysinfo

import (
    "fmt"
    "runtime"
    "strings"

    ps "github.com/shirou/gopsutil/v3/host"
)

type Info struct {
    OS           string `json:"os"`
    Kernel       string `json:"kernel"`
    Arch         string `json:"arch"`
    CPUs         int    `json:"cpus"`
    GoVersion    string `json:"goVersion"`
    Hostname     string `json:"hostname"`
    Uptime       string `json:"uptime"`
    Virtualization string `json:"virtualization,omitempty"`
}

func Collect() Info {
    host, _ := ps.Info()

    virt := ""
    if v, role := host.VirtualizationSystem, host.VirtualizationRole; v != "" {
        virt = fmt.Sprintf("%s (%s)", v, role)
    }

    return Info{
        OS:            fmt.Sprintf("%s %s", host.Platform, host.PlatformVersion),
        Kernel:        host.KernelVersion,
        Arch:          runtime.GOARCH,
        CPUs:          runtime.NumCPU(),
        Hostname:      host.Hostname,
        Uptime:        secondsToHuman(host.Uptime),
        Virtualization: virt,
    }
}

func secondsToHuman(sec uint64) string {
    d := sec / 86400
    h := (sec % 86400) / 3600
    return fmt.Sprintf("%dd %dh", d, h)
}

func (i Info) String() string {
    b := &strings.Builder{}
    fmt.Fprintf(b, "OS: %s\nArch: %s\nCPUs: %d\nKernel: %s\nHostname: %s\nUptime: %s",
        i.OS, i.Arch, i.CPUs, i.Kernel, i.Hostname, i.Uptime)
    if i.Virtualization != "" {
        fmt.Fprintf(b, "\nVM: %s", i.Virtualization)
    }
    return b.String()
}
