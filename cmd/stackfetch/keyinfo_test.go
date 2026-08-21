package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

type keyInfoPayload struct {
	Key          string   `json:"key"`
	Supported    bool     `json:"supported"`
	Dependencies []string `json:"dependencies"`
	Suggestions  []string `json:"suggestions"`
	Fetch        bool     `json:"fetch"`
	Error        string   `json:"error"`
}

func TestKeyInfoSupportedKey(t *testing.T) {
	resp := buildKeyInfoResponse("go", false)
	var out bytes.Buffer
	err := writeKeyInfoOutput(&out, resp)
	if err != nil {
		t.Fatalf("write output: %v", err)
	}
	if !strings.Contains(out.String(), "key: go") {
		t.Fatalf("missing key: %q", out.String())
	}
	if !strings.Contains(out.String(), "supported: true") {
		t.Fatalf("missing supported marker: %q", out.String())
	}
}

func TestKeyInfoUnsupportedSuggestsKeys(t *testing.T) {
	resp := buildKeyInfoResponse("gom", false)
	var out bytes.Buffer
	err := writeKeyInfoOutput(&out, resp)
	if err != nil {
		t.Fatalf("write output: %v", err)
	}
	if !strings.Contains(out.String(), "supported: false") {
		t.Fatalf("missing unsupported marker: %q", out.String())
	}
	if !strings.Contains(strings.ToLower(out.String()), "suggestions:") {
		t.Fatalf("missing suggestion line: %q", out.String())
	}
	if !strings.Contains(out.String(), "go") {
		t.Fatalf("expected go suggestion: %q", out.String())
	}
}

func TestKeyInfoJSON(t *testing.T) {
	resp := buildKeyInfoResponse("gom", false)
	rawPayload, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	var payload keyInfoPayload
	if err := json.Unmarshal(rawPayload, &payload); err != nil {
		t.Fatalf("decode payload: %v\nraw=%q", err, rawPayload)
	}
	if payload.Key != "gom" {
		t.Fatalf("key=%q", payload.Key)
	}
	if payload.Supported {
		t.Fatalf("expected unsupported")
	}
	if len(payload.Suggestions) == 0 {
		t.Fatalf("expected at least one suggestion")
	}
	if payload.Fetch {
		t.Fatalf("expected fetch=false by default")
	}
	if payload.Error != "" {
		t.Fatalf("unexpected error=%q", payload.Error)
	}
}
