package linkgraph

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuild(t *testing.T) {
	dir := t.TempDir()

	// Create two docs: a.md links to b.md, b.md links to a.md.
	writeFile(t, dir, "a.md", "# A\n\nSee [B](b.md).\n")
	writeFile(t, dir, "b.md", "# B\n\nSee [A](a.md) and [C](c.md).\n")

	aPath := filepath.Join(dir, "a.md")
	bPath := filepath.Join(dir, "b.md")
	cPath := filepath.Join(dir, "c.md") // doesn't exist

	files := func(yield func(string, error) bool) {
		if !yield(aPath, nil) {
			return
		}
		yield(bPath, nil)
	}

	g, err := Build(files)
	if err != nil {
		t.Fatal(err)
	}

	// Forward edges.
	if got := len(g.Forward[aPath]); got != 1 {
		t.Errorf("forward[a] = %d links, want 1", got)
	}
	if got := len(g.Forward[bPath]); got != 2 {
		t.Errorf("forward[b] = %d links, want 2", got)
	}

	// Reverse edges.
	if got := len(g.Reverse[bPath]); got != 1 {
		t.Errorf("reverse[b] = %d, want 1", got)
	}
	if got := len(g.Reverse[aPath]); got != 1 {
		t.Errorf("reverse[a] = %d, want 1", got)
	}
	if got := len(g.Reverse[cPath]); got != 1 {
		t.Errorf("reverse[c] = %d, want 1 (broken link still tracked)", got)
	}
}

func TestOrphans(t *testing.T) {
	dir := t.TempDir()

	// a.md links to b.md; c.md has no incoming links; README.md is a hub (excluded).
	writeFile(t, dir, "a.md", "# A\n\nSee [B](b.md).\n")
	writeFile(t, dir, "b.md", "# B\n")
	writeFile(t, dir, "c.md", "# C\n")
	writeFile(t, dir, "README.md", "# Hub\n")

	g, err := Build(seqFiles(dir, "a.md", "b.md", "c.md", "README.md"))
	if err != nil {
		t.Fatal(err)
	}

	orphans := collect(g.Orphans())

	// a.md has no incoming links → orphan. c.md has no incoming links → orphan.
	// b.md has incoming from a → not orphan. README.md is hub → excluded.
	if len(orphans) != 2 {
		t.Fatalf("got %d orphans, want 2: %v", len(orphans), orphans)
	}

	set := toSet(orphans)
	aPath := filepath.Join(dir, "a.md")
	cPath := filepath.Join(dir, "c.md")
	if !set[aPath] {
		t.Errorf("expected %s in orphans", aPath)
	}
	if !set[cPath] {
		t.Errorf("expected %s in orphans", cPath)
	}
}

func TestBrokenLinks(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, dir, "a.md", "# A\n\n[B](b.md) [Missing](missing.md)\n")
	writeFile(t, dir, "b.md", "# B\n\n[Also missing](nope.md)\n")

	g, err := Build(seqFiles(dir, "a.md", "b.md"))
	if err != nil {
		t.Fatal(err)
	}

	broken := collectBroken(g.BrokenLinks())
	if len(broken) != 2 {
		t.Fatalf("got %d broken links, want 2: %v", len(broken), broken)
	}
}

func seqFiles(dir string, names ...string) func(func(string, error) bool) {
	return func(yield func(string, error) bool) {
		for _, name := range names {
			if !yield(filepath.Join(dir, name), nil) {
				return
			}
		}
	}
}

func collect(seq func(func(string) bool)) []string {
	var out []string
	seq(func(s string) bool {
		out = append(out, s)
		return true
	})
	return out
}

func collectBroken(seq func(func(BrokenLink) bool)) []BrokenLink {
	var out []BrokenLink
	seq(func(b BrokenLink) bool {
		out = append(out, b)
		return true
	})
	return out
}

func toSet(ss []string) map[string]bool {
	m := make(map[string]bool, len(ss))
	for _, s := range ss {
		m[s] = true
	}
	return m
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
