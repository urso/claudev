package walker

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestWalkFindsMarkdownFiles(t *testing.T) {
	dir := t.TempDir()
	touch(t, dir, "a.md")
	touch(t, dir, "sub/b.md")
	touch(t, dir, "sub/deep/c.md")
	touch(t, dir, "not-md.txt")

	w := NewWalker(MarkdownFiles())
	got := collectPaths(t, w.Walk(dir))
	want := []string{
		filepath.Join(dir, "a.md"),
		filepath.Join(dir, "sub/b.md"),
		filepath.Join(dir, "sub/deep/c.md"),
	}
	assertPaths(t, want, got)
}

func TestWalkSkipsHiddenDirs(t *testing.T) {
	dir := t.TempDir()
	touch(t, dir, "visible.md")
	touch(t, dir, ".hidden/secret.md")

	got := collectPaths(t, NewWalker(MarkdownFiles()).Walk(dir))
	want := []string{filepath.Join(dir, "visible.md")}
	assertPaths(t, want, got)
}

func TestWalkSkipsCommonNonDocDirs(t *testing.T) {
	dir := t.TempDir()
	touch(t, dir, "doc.md")
	touch(t, dir, "node_modules/pkg.md")
	touch(t, dir, "vendor/lib.md")
	touch(t, dir, "__pycache__/cache.md")

	got := collectPaths(t, NewWalker(MarkdownFiles()).Walk(dir))
	want := []string{filepath.Join(dir, "doc.md")}
	assertPaths(t, want, got)
}

func TestWalkMultipleRoots(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()
	touch(t, dir1, "a.md")
	touch(t, dir2, "b.md")

	got := collectPaths(t, NewWalker(MarkdownFiles()).Walk(dir1, dir2))
	want := []string{
		filepath.Join(dir1, "a.md"),
		filepath.Join(dir2, "b.md"),
	}
	assertPaths(t, want, got)
}

func TestWalkNonExistentRoot(t *testing.T) {
	var errs []error
	for _, err := range NewWalker(MarkdownFiles()).Walk("/nonexistent/path") {
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		t.Fatal("expected error for non-existent root")
	}
}

func TestWalkCaseInsensitiveMd(t *testing.T) {
	dir := t.TempDir()
	touch(t, dir, "upper.MD")
	touch(t, dir, "mixed.Md")

	got := collectPaths(t, NewWalker(MarkdownFiles()).Walk(dir))
	want := []string{
		filepath.Join(dir, "mixed.Md"),
		filepath.Join(dir, "upper.MD"),
	}
	assertPaths(t, want, got)
}

func TestWalkEarlyBreak(t *testing.T) {
	dir := t.TempDir()
	touch(t, dir, "a.md")
	touch(t, dir, "b.md")
	touch(t, dir, "c.md")

	var count int
	for _, err := range NewWalker(MarkdownFiles()).Walk(dir) {
		if err != nil {
			t.Fatal(err)
		}
		count++
		if count >= 1 {
			break
		}
	}
	if count != 1 {
		t.Fatalf("expected 1 file, got %d", count)
	}
}

func TestWalkAllFiles(t *testing.T) {
	dir := t.TempDir()
	touch(t, dir, "a.md")
	touch(t, dir, "b.txt")
	touch(t, dir, "sub/c.go")

	got := collectPaths(t, NewWalker().Walk(dir))
	want := []string{
		filepath.Join(dir, "a.md"),
		filepath.Join(dir, "b.txt"),
		filepath.Join(dir, "sub/c.go"),
	}
	assertPaths(t, want, got)
}

func TestWalkMultipleIncludes(t *testing.T) {
	dir := t.TempDir()
	touch(t, dir, "a.md")
	touch(t, dir, "b.txt")
	touch(t, dir, "c.go")

	got := collectPaths(t, NewWalker(MarkdownFiles(), MatchSuffix(".txt")).Walk(dir))
	want := []string{
		filepath.Join(dir, "a.md"),
		filepath.Join(dir, "b.txt"),
	}
	assertPaths(t, want, got)
}

func TestWalkSkipHiddenFiles(t *testing.T) {
	dir := t.TempDir()
	touch(t, dir, "visible.md")
	touch(t, dir, ".hidden.md")

	got := collectPaths(t, NewWalker(SkipHidden()).Walk(dir))
	want := []string{filepath.Join(dir, "visible.md")}
	assertPaths(t, want, got)
}

func TestWalkIncludeAndExclude(t *testing.T) {
	dir := t.TempDir()
	touch(t, dir, "a.md")
	touch(t, dir, ".secret.md")
	touch(t, dir, "b.txt")

	got := collectPaths(t, NewWalker(MarkdownFiles(), SkipHidden()).Walk(dir))
	want := []string{filepath.Join(dir, "a.md")}
	assertPaths(t, want, got)
}

// helpers

func touch(t *testing.T, base string, rel string) {
	t.Helper()
	p := filepath.Join(base, rel)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, nil, 0o644); err != nil {
		t.Fatal(err)
	}
}

func collectPaths(t *testing.T, seq func(func(string, error) bool)) []string {
	t.Helper()
	var paths []string
	for path, err := range seq {
		if err != nil {
			t.Fatal(err)
		}
		paths = append(paths, path)
	}
	slices.Sort(paths)
	return paths
}

func assertPaths(t *testing.T, want, got []string) {
	t.Helper()
	slices.Sort(want)
	if !slices.Equal(want, got) {
		t.Errorf("paths mismatch\nwant: %v\ngot:  %v", want, got)
	}
}
