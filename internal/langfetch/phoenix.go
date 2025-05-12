package langfetch

import (
	"os"
	"path/filepath"
)

type PhoenixStack struct{}

func (p PhoenixStack) Detect(workDir string) (bool, string, error) {
	mixFile := filepath.Join(workDir, "mix.exs")
	if _, err := os.Stat(mixFile); err != nil {
		return false, "", nil
	}
	if fileContains(mixFile, `:phoenix`) && fileContains(mixFile, `:phoenix_live_view`) {
		return true, "Phoenix + LiveView", nil
	}
	return false, "", nil
}
