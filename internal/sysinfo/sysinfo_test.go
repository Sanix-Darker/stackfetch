package sysinfo

import (
	"strings"
	"testing"
)

func TestCollectBasic(t *testing.T) {
	got := Collect()
	if got.OS == "" || got.Arch == "" || got.CPUs <= 0 {
		t.Fatalf("Collect returned empty fields: %+v", got)
	}
}

func TestStringFormat(t *testing.T) {
	in := Info{
		OS:        "linux",
		Kernel:    "5.4.0",
		Arch:      "amd64",
		CPUs:      8,
		GoVersion: "go1.22",
	}
	s := in.String()
	if !strings.Contains(s, "linux") || !strings.Contains(s, "amd64") {
		t.Fatalf("string formatting incorrect: %s", s)
	}
}
