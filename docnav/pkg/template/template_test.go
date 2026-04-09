package template

import (
	"os"
	"path/filepath"
	"slices"
	"sort"
	"testing"
)

func writeTemplate(t *testing.T, dir, name, content string) {
	t.Helper()
	os.MkdirAll(dir, 0o755)
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func collect(t *testing.T, dirs ...string) []Template {
	t.Helper()
	var templates []Template
	for tmpl, err := range Discover(dirs...) {
		if err != nil {
			t.Fatal(err)
		}
		templates = append(templates, tmpl)
	}
	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Name < templates[j].Name
	})
	return templates
}

func TestDiscoverFindsTemplates(t *testing.T) {
	dir := t.TempDir()
	writeTemplate(t, dir, "note.md", `---
scope: A general note
---
# {title}
`)
	writeTemplate(t, dir, "runbook.md", `---
scope: Step-by-step procedure
---
# {title}
`)

	templates := collect(t, dir)
	if len(templates) != 2 {
		t.Fatalf("got %d templates, want 2", len(templates))
	}
	if templates[0].Name != "note" {
		t.Errorf("templates[0].Name = %q, want %q", templates[0].Name, "note")
	}
	if templates[0].Description != "A general note" {
		t.Errorf("templates[0].Description = %q, want %q", templates[0].Description, "A general note")
	}
	if templates[1].Name != "runbook" {
		t.Errorf("templates[1].Name = %q, want %q", templates[1].Name, "runbook")
	}
}

func TestDiscoverOverrideLastDirWins(t *testing.T) {
	bundled := t.TempDir()
	local := t.TempDir()

	writeTemplate(t, bundled, "note.md", `---
scope: Bundled note
---
`)
	writeTemplate(t, local, "note.md", `---
scope: Project-local note
---
`)

	templates := collect(t, bundled, local)
	if len(templates) != 1 {
		t.Fatalf("got %d templates, want 1", len(templates))
	}
	if templates[0].Description != "Project-local note" {
		t.Errorf("description = %q, want %q", templates[0].Description, "Project-local note")
	}
	if templates[0].Path != filepath.Join(local, "note.md") {
		t.Errorf("path = %q, want path in local dir", templates[0].Path)
	}
}

func TestDiscoverMergesDirs(t *testing.T) {
	bundled := t.TempDir()
	local := t.TempDir()

	writeTemplate(t, bundled, "note.md", `---
scope: Bundled note
---
`)
	writeTemplate(t, local, "custom.md", `---
scope: Custom template
---
`)

	templates := collect(t, bundled, local)
	if len(templates) != 2 {
		t.Fatalf("got %d templates, want 2", len(templates))
	}
	if templates[0].Name != "custom" {
		t.Errorf("templates[0].Name = %q, want %q", templates[0].Name, "custom")
	}
	if templates[1].Name != "note" {
		t.Errorf("templates[1].Name = %q, want %q", templates[1].Name, "note")
	}
}

func TestDiscoverSkipsMissingDir(t *testing.T) {
	templates := collect(t, "/nonexistent/path")
	if len(templates) != 0 {
		t.Errorf("got %d templates, want 0", len(templates))
	}
}

func TestDiscoverSkipsNonMarkdown(t *testing.T) {
	dir := t.TempDir()
	writeTemplate(t, dir, "note.md", `---
scope: A note
---
`)
	writeTemplate(t, dir, "readme.txt", "not a template")
	os.MkdirAll(filepath.Join(dir, "subdir"), 0o755)

	templates := collect(t, dir)
	if len(templates) != 1 {
		t.Fatalf("got %d templates, want 1", len(templates))
	}
}

func TestScanDir(t *testing.T) {
	dir := t.TempDir()
	writeTemplate(t, dir, "a.md", "# A")
	writeTemplate(t, dir, "b.md", "# B")
	writeTemplate(t, dir, "c.txt", "not markdown")
	os.MkdirAll(filepath.Join(dir, "subdir"), 0o755)

	var paths []string
	for path, err := range scanDir(dir) {
		if err != nil {
			t.Fatal(err)
		}
		paths = append(paths, filepath.Base(path))
	}
	sort.Strings(paths)

	want := []string{"a.md", "b.md"}
	if !slices.Equal(paths, want) {
		t.Errorf("scanDir = %v, want %v", paths, want)
	}
}

func TestScanDirMissing(t *testing.T) {
	count := 0
	for _, err := range scanDir("/nonexistent") {
		if err != nil {
			t.Fatal(err)
		}
		count++
	}
	if count != 0 {
		t.Errorf("got %d results for missing dir, want 0", count)
	}
}
