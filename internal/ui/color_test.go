package ui

import "testing"

func TestIconAsciiFallback(t *testing.T) {
	t.Setenv("STACKFETCH_ASCII", "1")
	if got, want := Icon("check"), "[ok]"; got != want {
		t.Fatalf("Icon(check) = %q, want %q", got, want)
	}
	if got, want := Icon("arrow"), "->"; got != want {
		t.Fatalf("Icon(arrow) = %q, want %q", got, want)
	}
}

func TestIconUtf8(t *testing.T) {
	t.Setenv("STACKFETCH_ASCII", "")
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("LC_ALL", "en_US.UTF-8")
	if got := Icon("check"); got != "✓" {
		t.Fatalf("Icon(check) = %q, want %q", got, "✓")
	}
	if got := Icon("branch"); got != "├─" {
		t.Fatalf("Icon(branch) = %q, want %q", got, "├─")
	}
}

func TestColorSuppression(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	if got, want := RedString("%s", "ok"), "ok"; got != want {
		t.Fatalf("RedString = %q, want %q", got, want)
	}

	t.Setenv("NO_COLOR", "")
	if got := RedString("%s", "ok"); got == "ok" {
		t.Fatalf("expected color escape sequence when NO_COLOR disabled")
	}
}
