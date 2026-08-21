package main

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/sanix-darker/stackfetch/internal/langfetch"
	"github.com/spf13/cobra"
)

type keysJSONResponse struct {
	Count int      `json:"count"`
	Keys  []string `json:"keys"`
}

func TestFilterStringSlice(t *testing.T) {
	values := []string{"alpha", "charlie", "bravo", "bravo"}

	got := filterStringSlice(values, "", "desc")
	want := []string{"charlie", "bravo", "bravo", "alpha"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v want=%v", got, want)
	}

	got = filterStringSlice(values, "BR", "asc")
	want = []string{"bravo", "bravo"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestFilterStringSliceRejectsInvalidSort(t *testing.T) {
	_, err := runKeysCommand(t, "keys", "--sort", "zigzag")
	if err == nil {
		t.Fatalf("expected invalid sort error")
	}
}

func TestKeysCommand(t *testing.T) {
	out, err := runKeysCommand(t, "keys")
	if err != nil {
		t.Fatalf("run keys: %v", err)
	}

	got := nonEmptyLines(strings.TrimSpace(out))
	want := langfetchKeysForTest()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("full key list mismatch\n\ngot=%v\n\nwant=%v", got, want)
	}
}

func TestKeysCommandWithFilterAndSort(t *testing.T) {
	out, err := runKeysCommand(t, "keys", "--contains", "go")
	if err != nil {
		t.Fatalf("run filtered keys: %v", err)
	}

	got := nonEmptyLines(strings.TrimSpace(out))
	want := filterStringSlice(langfetchKeysForTest(), "go", "asc")
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("filtered key list mismatch\n\ngot=%v\n\nwant=%v", got, want)
	}

	for _, key := range got {
		if !strings.Contains(key, "go") {
			t.Fatalf("key %q does not match filter", key)
		}
	}
}

func TestKeysCommandWithDescSort(t *testing.T) {
	out, err := runKeysCommand(t, "keys", "--sort", "desc")
	if err != nil {
		t.Fatalf("run desc keys: %v", err)
	}

	got := nonEmptyLines(strings.TrimSpace(out))
	want := filterStringSlice(langfetchKeysForTest(), "", "desc")
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("desc key list mismatch\n\ngot=%v\n\nwant=%v", got, want)
	}
}

func TestKeysCommandJSON(t *testing.T) {
	out, err := runKeysCommand(t, "--json", "keys", "--contains", "go")
	if err != nil {
		t.Fatalf("run json keys: %v", err)
	}

	var payload keysJSONResponse
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("decode json payload: %v\nraw=%q", err, out)
	}

	want := filterStringSlice(langfetchKeysForTest(), "go", "asc")
	if payload.Count != len(want) {
		t.Fatalf("count=%d want=%d", payload.Count, len(want))
	}
	if !reflect.DeepEqual(payload.Keys, want) {
		t.Fatalf("json keys mismatch\n\ngot=%v\n\nwant=%v", payload.Keys, want)
	}
}

func runKeysCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var jsonOut bool
	var markdownOut bool
	var out bytes.Buffer
	root := &cobra.Command{
		Use: "stackfetch",
	}
	root.PersistentFlags().BoolVarP(&jsonOut, "json", "j", false, "For JSON output")
	root.PersistentFlags().BoolVar(&markdownOut, "markdown", false, "For Markdown output")
	root.AddCommand(newKeysCmd())
	root.SetOut(&out)
	root.SetErr(&out)
	root.SilenceErrors = true
	root.SilenceUsage = true
	root.SetArgs(args)
	err := root.Execute()
	return out.String(), err
}

func nonEmptyLines(raw string) []string {
	lines := strings.Split(raw, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			out = append(out, line)
		}
	}
	return out
}

func langfetchKeysForTest() []string {
	return langfetch.Keys()
}
