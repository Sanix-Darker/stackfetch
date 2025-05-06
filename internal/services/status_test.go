package services

import (
    "testing"
)

func TestCheck(t *testing.T) {
    // Assuming 'go' is always installed; test Check() picks it up
    statuses := Check([]string{"go", "go"}) // duplicate should be skipped
    if len(statuses) != 1 {
        t.Fatalf("expected 1 status entry, got %d", len(statuses))
    }
    if !statuses[0].Installed {
        t.Errorf("expected 'go' to be reported as installed")
    }
}

func TestStatusByName(t *testing.T) {
    s := StatusByName("this-does-not-exist-xyz")
    if s.Installed || s.Running {
        t.Errorf("expected non-existent binary to report Installed=false, Running=false")
    }
}
