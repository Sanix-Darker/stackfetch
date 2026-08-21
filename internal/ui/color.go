package ui

import (
	"fmt"
	"os"
	"strings"
)

type symbolPair struct {
	unicode string
	ascii   string
}

var symbolMap = map[string]symbolPair{
	"lightbulb": {"💡", "[+]"},
	"satellite": {"🛰", "[cloud]"},
	"lock":      {"🔐", "[sec]"},
	"flag":      {"🏁", "[risk]"},
	"check":     {"✓", "[ok]"},
	"warning":   {"⚠", "[!]"},
	"cross":     {"✗", "[x]"},
	"running":   {"●", "[*]"},
	"partial":   {"◐", "[>]"},
	"blank":     {"◌", "[ ]"},
	"arrow":     {"→", "->"},
	"branch":    {"├─", "|-"},
	"leaf":      {"└─", "`-"},
	"goal":      {"▸", ">"},
	"shield":    {"🛡", "[R]"},
}

func colorsEnabled() bool {
	if os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" || os.Getenv("CLICOLOR") == "0" {
		return false
	}
	if os.Getenv("CLICOLOR_FORCE") == "1" {
		return true
	}
	return true
}

func applyColor(code string, format string, a ...interface{}) string {
	if !colorsEnabled() {
		return fmt.Sprintf(format, a...)
	}
	return fmt.Sprintf("\033["+code+"m"+format+"\033[0m", a...)
}

// AsciiSafe returns true when symbols should be ASCII-only.
func AsciiSafe() bool {
	if os.Getenv("STACKFETCH_ASCII") != "" || os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb" {
		return true
	}

	loc := strings.ToUpper(strings.TrimSpace(strings.Join([]string{os.Getenv("LC_ALL"), os.Getenv("LANG"), os.Getenv("LC_CTYPE")}, "")))
	if loc == "" {
		return false
	}
	if strings.Contains(loc, "UTF-8") || strings.Contains(loc, "UTF8") {
		return false
	}
	if strings.HasPrefix(loc, "C") && loc != "C.UTF-8" && loc != "C.UTF8" {
		return true
	}
	return false
}

// Icon returns a terminal-safe symbol for the requested token.
func Icon(name string) string {
	item, ok := symbolMap[name]
	if !ok {
		return name
	}
	if AsciiSafe() {
		return item.ascii
	}
	return item.unicode
}

// HasColor reports whether we can emit ANSI color.
func HasColor() bool {
	return colorsEnabled()
}

// RedString wraps the input string with ANSI codes for red color
func RedString(format string, a ...interface{}) string {
	return applyColor("31", format, a...)
}

// GreenString wraps the input string with ANSI codes for green color
func GreenString(format string, a ...interface{}) string {
	return applyColor("32", format, a...)
}

// YellowString wraps the input string with ANSI codes for yellow color
func YellowString(format string, a ...interface{}) string {
	return applyColor("33", format, a...)
}

// BlueString wraps the input string with ANSI codes for blue color
func BlueString(format string, a ...interface{}) string {
	return applyColor("34", format, a...)
}

// CyanString wraps the input string with ANSI codes for cyan color
func CyanString(format string, a ...interface{}) string {
	return applyColor("36", format, a...)
}

// GrayString wraps the input string with ANSI codes for gray color
func GrayString(format string, a ...interface{}) string {
	return applyColor("90", format, a...)
}
