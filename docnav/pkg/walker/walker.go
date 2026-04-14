package walker

import (
	"io/fs"
	"iter"
	"path/filepath"
	"strings"
)

// Walker discovers files in directory trees.
type Walker struct {
	skipDirs map[string]struct{}
	include  []func(string) bool
	exclude  []func(string) bool
}

type option func(*Walker)

// NewWalker creates a Walker with the given options.
// If no skip-dir options are provided, DefaultSkipDirs is applied.
func NewWalker(opts ...option) *Walker {
	w := &Walker{}
	for _, opt := range opts {
		opt(w)
	}
	if w.skipDirs == nil {
		DefaultSkipDirs()(w)
	}
	return w
}

// SkipDirs returns an option that configures directories to skip during walks.
func SkipDirs(dirs ...string) option {
	return func(w *Walker) {
		m := make(map[string]struct{}, len(dirs))
		for _, d := range dirs {
			m[d] = struct{}{}
		}
		w.skipDirs = m
	}
}

// defaultSkipDirs are directory names that are always skipped during walks.
var defaultSkipDirs = []string{
	"node_modules",
	"vendor",
	"__pycache__",
	".git",
}

// DefaultSkipDirs returns an option that skips common non-doc directories.
func DefaultSkipDirs() option {
	return SkipDirs(defaultSkipDirs...)
}

// Include adds an include filter. A file is yielded if it matches any include filter.
// If no include filters are set, all files pass.
func Include(fn func(string) bool) option {
	return func(w *Walker) {
		w.include = append(w.include, fn)
	}
}

// Exclude adds an exclude filter. A file is rejected if it matches any exclude filter.
// Exclude filters are applied after include filters.
func Exclude(fn func(string) bool) option {
	return func(w *Walker) {
		w.exclude = append(w.exclude, fn)
	}
}

// MatchSuffix returns an include filter for files with the given suffix (case-insensitive).
func MatchSuffix(suffix string) option {
	lower := strings.ToLower(suffix)
	return Include(func(name string) bool {
		return strings.HasSuffix(strings.ToLower(name), lower)
	})
}

// MarkdownFiles returns an include filter for .md files.
func MarkdownFiles() option {
	return MatchSuffix(".md")
}

// SkipHidden returns an exclude filter for hidden files (names starting with ".").
func SkipHidden() option {
	return Exclude(func(name string) bool {
		return strings.HasPrefix(name, ".")
	})
}

// accepts reports whether a filename passes the walker's include/exclude filters.
func (w *Walker) accepts(name string) bool {
	if len(w.include) > 0 {
		matched := false
		for _, fn := range w.include {
			if fn(name) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	for _, fn := range w.exclude {
		if fn(name) {
			return false
		}
	}
	return true
}

// Walk yields file paths found under the given root directories.
// Hidden directories (starting with ".") and configured skip directories are skipped.
// Files are filtered by include/exclude options.
// Errors are yielded as ("", err) and do not stop the walk.
func (w *Walker) Walk(roots ...string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		for _, root := range roots {
			root, err := filepath.Abs(root)
			if err != nil {
				if !yield("", err) {
					return
				}
				continue
			}
			stopped := false
			filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
				if stopped {
					return fs.SkipAll
				}
				if err != nil {
					if !yield("", err) {
						stopped = true
						return fs.SkipAll
					}
					return nil
				}
				if d.IsDir() {
					name := d.Name()
					if _, skip := w.skipDirs[name]; skip {
						return fs.SkipDir
					}
					// Skip hidden directories (but not the root itself).
					if name != "." && strings.HasPrefix(name, ".") && path != root {
						return fs.SkipDir
					}
					return nil
				}
				if !w.accepts(d.Name()) {
					return nil
				}
				if !yield(path, nil) {
					stopped = true
					return fs.SkipAll
				}
				return nil
			})
			if stopped {
				return
			}
		}
	}
}
