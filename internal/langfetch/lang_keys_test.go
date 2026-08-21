package langfetch

import "testing"

func TestIsSupported(t *testing.T) {
	if !IsSupported("go") {
		t.Fatalf("go should be supported")
	}
	if IsSupported("probably-not-a-key") {
		t.Fatalf("unexpected supported key")
	}
}

func TestSuggestKeys(t *testing.T) {
	suggestions := SuggestKeys("gom", 3)
	if len(suggestions) == 0 {
		t.Fatalf("expected at least one suggestion")
	}
	if suggestions[0] != "go" {
		t.Fatalf("expected go first suggestion, got %q", suggestions[0])
	}
	if len(suggestions) > 3 {
		t.Fatalf("max suggestions exceeded: %d", len(suggestions))
	}

	none := SuggestKeys("", 3)
	if len(none) != 0 {
		t.Fatalf("expected empty suggestions for empty query")
	}
}
