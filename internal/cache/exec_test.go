package cache

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunCachesOutput(t *testing.T) {
	tmp := t.TempDir()
	scriptPath := filepath.Join(tmp, "version-print")
	if err := os.WriteFile(scriptPath, []byte("#!/bin/sh\necho first-run\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	origDir := dir
	origPath := os.Getenv("PATH")
	defer func() {
		dir = origDir
		_ = os.Setenv("PATH", origPath)
	}()

	dir = t.TempDir()
	_ = os.MkdirAll(filepath.Dir(scriptPath), 0o755)
	if err := os.Symlink(scriptPath, filepath.Join(dir, "stackfetch-version-print")); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}
	_ = os.Setenv("PATH", dir+string(os.PathListSeparator)+origPath)

	out1, err := Run("stackfetch-version-print")
	if err != nil {
		t.Fatalf("first run failed: %v", err)
	}
	if string(out1) != "first-run\n" {
		t.Fatalf("first output mismatch: %s", out1)
	}

	// Change command behavior, but cache should still serve the original output.
	if err := os.WriteFile(scriptPath, []byte("#!/bin/sh\necho second-run\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	out2, err := Run("stackfetch-version-print")
	if err != nil {
		t.Fatalf("second run failed: %v", err)
	}
	if string(out2) != "first-run\n" {
		t.Fatalf("cached output was not used: %s", out2)
	}
}
