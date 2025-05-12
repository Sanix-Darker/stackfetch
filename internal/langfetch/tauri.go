package langfetch

import (
	"os"
	"path/filepath"
)

type TauriStack struct{}

func (t TauriStack) Detect(workDir string) (bool, string, error) {
	tauriConf := filepath.Join(workDir, "src-tauri", "tauri.conf.json")
	if _, err := os.Stat(tauriConf); err != nil {
		return false, "", nil
	}
	packageJson := filepath.Join(workDir, "package.json")
	if hasDependency(packageJson, "solid-js") {
		return true, "Tauri + SolidJS", nil
	}
	return true, "Tauri", nil
}
