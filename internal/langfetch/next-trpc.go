package langfetch

import (
	"os"
	"path/filepath"
)

type NextTRPCStack struct{}

func (n NextTRPCStack) Detect(workDir string) (bool, string, error) {
	packageJson := filepath.Join(workDir, "package.json")
	if _, err := os.Stat(packageJson); err != nil {
		return false, "", nil
	}
	// Check for Next.js + tRPC deps
	if hasDependency(packageJson, "next") && hasDependency(packageJson, "@trpc/server") {
		return true, "Next.js + tRPC Stack", nil
	}
	return false, "", nil
}
