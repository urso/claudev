package fsutil

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ErrNotFound is returned when FindProjRoot cannot locate the named directory.
var ErrNotFound = errors.New("directory not found")

// FindProjRoot walks up from start looking for a directory named name. If none
// is found, it falls back to <git-root>/<name>. Returns ErrNotFound if neither
// lookup succeeds.
func FindProjRoot(start, name string) (string, error) {
	abs, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	dir := abs
	for {
		candidate := filepath.Join(dir, name)
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	if root, err := gitRoot(abs); err == nil {
		candidate := filepath.Join(root, name)
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}
	}

	return "", ErrNotFound
}

func gitRoot(start string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = start
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
