package main

import (
	"testing"

	"github.com/sanix-darker/stackfetch/internal/services"
)

func TestReadinessFromStatuses(t *testing.T) {
	tests := []struct {
		name      string
		statuses  []services.Status
		wantScore int
	}{
		{
			name:      "independent",
			statuses:  nil,
			wantScore: 0,
		},
		{
			name:      "all running",
			statuses:  []services.Status{{Installed: true, Running: true}, {Installed: true, Running: true}},
			wantScore: 100,
		},
		{
			name:      "installed only",
			statuses:  []services.Status{{Installed: true}, {Installed: true}},
			wantScore: 66,
		},
		{
			name:      "mixed",
			statuses:  []services.Status{{Installed: true, Running: true}, {Installed: true}, {Installed: false, Running: false}},
			wantScore: 55,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := readinessFromStatuses(tt.statuses)
			if r.score != tt.wantScore {
				t.Fatalf("score = %d, want %d", r.score, tt.wantScore)
			}
			if tt.name == "independent" && r.total != 0 {
				t.Fatalf("total = %d, want 0", r.total)
			}
		})
	}
}

func TestReadinessBanner(t *testing.T) {
	got := readinessBanner(readinessFromStatuses([]services.Status{{Installed: true, Running: true}, {Installed: true, Running: true}}))
	if got != "excellent (100%)" {
		t.Fatalf("got=%q", got)
	}

	got = readinessBanner(readinessFromStatuses([]services.Status{{Installed: true}, {Installed: true}}))
	if got != "warning (66%)" {
		t.Fatalf("got=%q", got)
	}

	got = readinessBanner(readinessFromStatuses([]services.Status{{Installed: false}}))
	if got != "critical (0%)" {
		t.Fatalf("got=%q", got)
	}

	got = readinessBanner(readinessFromStatuses(nil))
	if got != "independent" {
		t.Fatalf("got=%q", got)
	}
}
