package scan

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func writeTree(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestRepo_ExtensionsAndDirs(t *testing.T) {
	root := t.TempDir()
	writeTree(t, root, map[string]string{
		"go.mod":          "module x",
		"main.go":         "package main",
		"pkg/a/a.go":      "package a",
		"pkg/a/a_test.go": "package a",
		"pkg/b/b.go":      "package b",
		"docs/r.md":       "# r",
		"docs/x.md":       "# x",
	})

	sig, err := Repo(root, RepoOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if sig.FileCount != 7 {
		t.Errorf("file count = %d, want 7", sig.FileCount)
	}
	got := map[string]int{}
	for _, e := range sig.Extensions {
		got[e.Ext] = e.Count
	}
	if got[".go"] != 4 || got[".md"] != 2 {
		t.Errorf("ext counts = %v", got)
	}
	dirs := map[string]int{}
	for _, d := range sig.TopDirs {
		dirs[d.Dir] = d.Files
	}
	if dirs["pkg"] != 3 || dirs["docs"] != 2 || dirs["(root)"] != 2 {
		t.Errorf("dir counts = %v", dirs)
	}
}

func TestRepo_IgnoreDirs(t *testing.T) {
	root := t.TempDir()
	writeTree(t, root, map[string]string{
		"main.go":              "package main",
		"node_modules/x/y.js":  "x",
		".git/HEAD":            "ref",
		"vendor/foo/bar.go":    "package foo",
		".discovery/scan.json": "{}",
	})
	sig, err := Repo(root, RepoOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if sig.FileCount != 1 {
		t.Errorf("file count = %d, want 1 (ignored dirs leaked)", sig.FileCount)
	}
}

func TestRepo_Markers(t *testing.T) {
	root := t.TempDir()
	writeTree(t, root, map[string]string{
		"go.mod":                              "module x",
		"package.json":                        "{}",
		"docker-compose.yml":                  "x",
		".github/workflows/ci.yml":            "x",
		"k8s/deployment.yaml":                 "x",
		"src/main.go":                         "package main",
	})
	sig, err := Repo(root, RepoOptions{})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		".github/workflows/", "docker-compose.yml", "go.mod", "k8s/", "package.json",
	}
	for _, w := range want {
		if !slices.Contains(sig.Markers, w) {
			t.Errorf("missing marker %q in %v", w, sig.Markers)
		}
	}
}

func TestRepo_TopExtCap(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{}
	for i := range 10 {
		files[filepath.Join("d", "f"+string(rune('a'+i))+".x"+string(rune('a'+i)))] = "x"
	}
	writeTree(t, root, files)
	sig, err := Repo(root, RepoOptions{TopExtensions: 3})
	if err != nil {
		t.Fatal(err)
	}
	if len(sig.Extensions) != 3 {
		t.Errorf("ext len = %d, want 3", len(sig.Extensions))
	}
}
