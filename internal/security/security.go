package security

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Report struct {
	Root      bool   `json:"root"`
	SELinux   string `json:"selinux"`
	SSHOpen   bool   `json:"sshOpen"`
	KernelEOL bool   `json:"kernelEOL"`
}

func Collect() Report {
	return Report{
		Root:      os.Geteuid() == 0,
		SELinux:   selinuxMode(),
		SSHOpen:   portOpen(22),
		KernelEOL: kernelEOL(),
	}
}

func portOpen(p int) bool {
	c, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", p), 250*time.Millisecond)
	if err == nil {
		_ = c.Close()
		return true
	}
	return false
}

func selinuxMode() string {
	raw, err := os.ReadFile("/sys/fs/selinux/enforce")
	if err != nil {
		return "disabled"
	}
	if len(raw) > 0 && raw[0] == '1' {
		return "enforcing"
	}
	return "permissive"
}

func kernelEOL() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	parts := strings.Split(strings.TrimSpace(string(out)), ".")
	if len(parts) < 2 {
		return false
	}
	v, _ := strconv.ParseFloat(parts[0]+"."+parts[1], 64)
	return v < 5.15
}
