package guess

import (
    "os"
    "testing"
)

func TestGuess(t *testing.T) {
    os.RemoveAll("tmp")
    os.MkdirAll("tmp/level1/package.json", 0o755)
    // NOTE: For now we drop support for 3 level (for better performances)
    //os.WriteFile("tmp/level1/level2/package.json", []byte("{}"), 0o644)

    got := Guess("tmp")
    if len(got) == 0 {
        t.Fatalf("expected some guesses, got %v", got)
    }
    found := false
    for _, k := range got {
        if k == "node" {
            found = true
        }
    }
    if !found {
        t.Fatalf("expected 'node' in guesses: %v", got)
    }
}
