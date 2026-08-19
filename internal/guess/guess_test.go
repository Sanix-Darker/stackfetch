package guess

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGuessDetectsKnownStackFile(t *testing.T) {
	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, "level1"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "level1", "package.json"), []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := Guess(tmp)
	if len(got) == 0 {
		t.Fatalf("expected some guesses, got %v", got)
	}

	foundNode := false
	for _, k := range got {
		if k == "node" {
			foundNode = true
		}
	}
	if !foundNode {
		t.Fatalf("expected 'node' in guesses: %v", got)
	}
}

func TestGuessSupportsWildcardBasePatternsAndDepthLimits(t *testing.T) {
	tmp := t.TempDir()
	files := []string{
		"build.sbt",
		"main.csproj",
		"package.json",
		"requirements.txt",
	}
	for _, name := range files {
		if err := os.WriteFile(filepath.Join(tmp, name), []byte(""), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	got := Guess(tmp)
	if len(got) == 0 {
		t.Fatalf("expected guesses from mixed project files, got %v", got)
	}

	if len(got) > maxHits {
		t.Fatalf("expected at most %d unique guesses due hit cap, got %d", maxHits, len(got))
	}

	found := map[string]bool{
		"scala":  false,
		"dotnet": false,
		"node":   false,
		"python": false,
	}
	for _, v := range got {
		if _, ok := found[v]; ok {
			found[v] = true
		}
	}
	for name, ok := range found {
		if !ok {
			t.Fatalf("expected %s to be detected, got %v", name, got)
		}
	}
}
