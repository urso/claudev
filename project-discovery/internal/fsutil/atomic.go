package fsutil

import (
	"os"
	"path/filepath"
)

// AtomicWrite writes data to path atomically by writing to a temporary file
// in the same directory and renaming. If the target file exists, its
// permissions are preserved.
func AtomicWrite(path string, data []byte) error {
	dir := filepath.Dir(path)

	perm := os.FileMode(0644)
	if info, err := os.Stat(path); err == nil {
		perm = info.Mode().Perm()
	}

	tmp, err := os.CreateTemp(dir, ".discover-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := tmp.Chmod(perm); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return os.Rename(tmpPath, path)
}
