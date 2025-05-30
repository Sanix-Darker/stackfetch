package langfetch

import (
	"errors"
	"os/exec"
	"strings"
)

type fakeExec map[string][]byte

func (f fakeExec) Command(name string, args ...string) *exec.Cmd {
	key := name + " " + strings.Join(args, " ")
	fn := func() ([]byte, error) {
		if v, ok := f[key]; ok {
			return v, nil
		}
		return nil, errors.New("not found")
	}
	return helperCmd(fn)
}

func helperCmd(fn func() ([]byte, error)) *exec.Cmd {
	if b, err := fn(); err == nil {
		return exec.Command("echo", string(b))
	}
	return exec.Command("false")
}
