package langfetch

import (
	"fmt"
	"os/exec"
)

type BunStack struct{}

func (b BunStack) Detect(workDir string) (bool, string, error) {
	if _, err := exec.LookPath("bun"); err != nil {
		return false, "", nil
	}
	versionCmd := exec.Command("bun", "--version")
	version, err := versionCmd.Output()
	if err != nil {
		return false, "", err
	}
	return true, fmt.Sprintf("Bun %s", string(version)), nil
}
