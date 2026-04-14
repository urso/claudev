package template

import (
	"iter"
	"os"
	"path/filepath"
	"strings"

	"github.com/urso/claudev/docnav/pkg/parser"
)

// Template represents a discovered document template.
type Template struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Path        string `json:"path"`
}

// Discover scans the given directories for markdown template files.
// Directories are processed in order; when multiple directories contain
// a template with the same name, the last one wins.
// The iterator yields templates in discovery order (not sorted).
func Discover(dirs ...string) iter.Seq2[Template, error] {
	return func(yield func(Template, error) bool) {
		byName := make(map[string]Template)
		for _, dir := range dirs {
			for path, err := range scanDir(dir) {
				if err != nil {
					yield(Template{}, err)
					return
				}
				name := strings.TrimSuffix(filepath.Base(path), ".md")
				byName[name] = Template{
					Name: name,
					Path: path,
				}
			}
		}

		for name, t := range byName {
			content, err := os.ReadFile(t.Path)
			if err != nil {
				if !yield(Template{}, err) {
					return
				}
				continue
			}
			fm := parser.ParseFrontmatter(content)
			t.Description = fm.Scope
			byName[name] = t
			if !yield(t, nil) {
				return
			}
		}
	}
}

// scanDir yields paths to markdown files in a single directory.
// Returns an empty sequence if the directory does not exist.
func scanDir(dir string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		entries, err := os.ReadDir(dir)
		if os.IsNotExist(err) {
			return
		}
		if err != nil {
			yield("", err)
			return
		}

		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			if !yield(filepath.Join(dir, e.Name()), nil) {
				return
			}
		}
	}
}
