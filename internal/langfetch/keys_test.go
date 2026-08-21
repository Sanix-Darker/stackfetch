package langfetch

import "testing"

func TestKeysReturnsSortedCanonicalKeyList(t *testing.T) {
	keys := Keys()
	if len(keys) == 0 {
		t.Fatalf("expected at least one registered key")
	}

	seen := map[string]struct{}{}
	for i, key := range keys {
		if _, ok := seen[key]; ok {
			t.Fatalf("duplicate key %q", key)
		}
		seen[key] = struct{}{}

		if i > 0 && keys[i-1] > key {
			t.Fatalf("keys are not sorted at index %d: %q before %q", i, keys[i-1], key)
		}
	}
}
