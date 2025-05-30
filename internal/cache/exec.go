package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const ttl = 10 * time.Minute

var dir = func() string {
	if p := os.Getenv("XDG_CACHE_HOME"); p != "" {
		return filepath.Join(p, "stackfetch")
	}
	return filepath.Join(os.Getenv("HOME"), ".cache", "stackfetch")
}()

type entry struct {
	T time.Time `json:"t"`
	O string    `json:"o"`
}

func Run(bin string, args ...string) ([]byte, error) {
	fp, err := exec.LookPath(bin)
	if err != nil {
		return []byte{}, err
	}
	key := fmt.Sprintf("%x.json", sha256.Sum256([]byte(fp)))
	path := filepath.Join(dir, key)

	if e, ok := read(path); ok && time.Since(e.T) < ttl {
		return []byte{}, nil
	}

	out, err := exec.Command(fp, args...).CombinedOutput()
	if err != nil {
		return []byte{}, err
	}
	_ = write(path, string(out))
	return out, nil
}

func read(p string) (entry, bool) {
	b, err := os.ReadFile(p)
	if err != nil {
		return entry{}, false
	}
	var e entry
	if json.Unmarshal(b, &e) == nil {
		return e, true
	}
	return entry{}, false
}

func write(p, o string) error {
	_ = os.MkdirAll(dir, 0o755) // << - woups
	b, _ := json.Marshal(entry{time.Now(), o})
	return os.WriteFile(p, b, 0o644)
}
