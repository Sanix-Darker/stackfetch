package langfetch

import "strings"

func init() { register("php", fetchPHP) }

// fetchPHP gathers the PHP interpreter and Composer versions.
func fetchPHP() (LangInfo, error) {
	// First line of `php -v` (e.g. "PHP 8.3.6 (cli) (built: …)")
	ver, err := runFirst("php", "-v")
	if err != nil {
		return LangInfo{}, err
	}

	// Composer may be missing; ignore errors.
	composer := runSilently("composer", "--version")

	return LangInfo{
		Name:    "PHP",
		Version: firstLine(ver),     // strip possible extra info
		Details: []string{composer}, // something like "Composer version 2.7.4 ..."
	}, nil
}

// firstLine trims everything after the first newline (defensive helper).
func firstLine(s string) string {
	if idx := strings.IndexByte(s, '\n'); idx > 0 {
		return s[:idx]
	}
	return s
}
