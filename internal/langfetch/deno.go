package langfetch

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type DenoStack struct{}

func (d DenoStack) Detect(workDir string) (bool, string, error) {
	if _, err := os.Stat(filepath.Join(workDir, "deno.json")); err == nil {
		versionCmd := exec.Command("deno", "--version")
		version, err := versionCmd.Output()
		if err != nil {
			return false, "", err
		}
		return true, fmt.Sprintf("Deno %s", string(bytes.Split(version, []byte("\n"))[0])), nil
	}
	return false, "", nil
}
