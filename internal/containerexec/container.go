package containerexec

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Runtime string

const (
	Docker Runtime = "docker"
	Podman Runtime = "podman"
)

func DetectRuntime() Runtime {
	if _, err := exec.LookPath("docker"); err == nil {
		return Docker
	}
	if _, err := exec.LookPath("podman"); err == nil {
		return Podman
	}
	return ""
}

func Exec(runtime Runtime, id, bin string, args ...string) (string, error) {
	if runtime == "" {
		runtime = DetectRuntime()
	}
	if runtime == "" {
		return "", fmt.Errorf("no container runtime found")
	}
	cmdArgs := append([]string{"exec", "-i", id, bin}, args...)
	cmd := exec.Command(string(runtime), cmdArgs...)
	var buf bytes.Buffer
	cmd.Stdout, cmd.Stderr = &buf, &buf

	done := make(chan error, 1)
	go func() { done <- cmd.Run() }()
	select {
	case err := <-done:
		return strings.TrimSpace(buf.String()), err
	case <-time.After(15 * time.Second):
		_ = cmd.Process.Kill()
		return "", fmt.Errorf("container exec timeout")
	}
}

// Wrapper to be compatible with langfetch.execCmd
func Runner(id string, rt Runtime) func(string, ...string) *exec.Cmd {
	return func(name string, args ...string) *exec.Cmd {
		all := append([]string{name}, args...)
		return exec.Command(os.Args[0],
			append(
				[]string{
					"__stackfetch_container_proxy", id, string(rt),
				}, all...,
			)...,
		)
	}
}
