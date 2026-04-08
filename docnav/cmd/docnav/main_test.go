package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/urso/claudev/docnav/pkg/document"
	"github.com/urso/claudev/docnav/pkg/parser"
	"github.com/urso/claudev/docnav/pkg/walker"
)

func TestMatchWatch(t *testing.T) {
	tests := []struct {
		watch, changed string
		want           bool
	}{
		// Exact match.
		{"api/v1/types.go", "api/v1/types.go", true},
		{"api/v1/types.go", "api/v1/other.go", false},

		// Directory prefix.
		{"internal/reconciler/", "internal/reconciler/pool.go", true},
		{"internal/reconciler/", "internal/reconciler/sub/deep.go", true},
		{"internal/reconciler/", "internal/other/pool.go", false},

		// Glob patterns.
		{"*.go", "main.go", true},
		{"*.go", "dir/main.go", false}, // filepath.Match doesn't cross /
		{"api/v1/*.go", "api/v1/types.go", true},
		{"api/v1/*.go", "api/v2/types.go", false},

		// Doublestar patterns.
		{"internal/**/*.go", "internal/reconciler/pool.go", true},
		{"internal/**/*.go", "internal/reconciler/sub/deep.go", true},
		{"internal/**/*.go", "internal/other.txt", false},
		{"**/*.go", "any/deep/path/file.go", true},
		{"**/*.go", "file.go", true},
	}

	for _, tt := range tests {
		got := matchWatch(tt.watch, tt.changed)
		if got != tt.want {
			t.Errorf("matchWatch(%q, %q) = %v, want %v", tt.watch, tt.changed, got, tt.want)
		}
	}
}

func writeDoc(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	os.MkdirAll(filepath.Dir(path), 0o755)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestTagAggregation(t *testing.T) {
	dir := t.TempDir()

	writeDoc(t, dir, "a.md", `---
title: Doc A
tags: [go, testing]
---
# A
`)
	writeDoc(t, dir, "b.md", `---
title: Doc B
tags: [go, design]
---
# B
`)
	writeDoc(t, dir, "c.md", `---
title: Doc C
---
# C (no tags)
`)

	w := walker.NewWalker(walker.MarkdownFiles())
	tagDocs := make(map[string][]string)
	for path, err := range w.Walk(dir) {
		if err != nil {
			t.Fatal(err)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		doc := parser.ParseDocument(path, content)
		for _, tag := range doc.Frontmatter.Tags {
			tagDocs[tag] = append(tagDocs[tag], doc.Path)
		}
	}

	if len(tagDocs["go"]) != 2 {
		t.Errorf("tag 'go' count = %d, want 2", len(tagDocs["go"]))
	}
	if len(tagDocs["testing"]) != 1 {
		t.Errorf("tag 'testing' count = %d, want 1", len(tagDocs["testing"]))
	}
	if len(tagDocs["design"]) != 1 {
		t.Errorf("tag 'design' count = %d, want 1", len(tagDocs["design"]))
	}
	if _, ok := tagDocs["nonexistent"]; ok {
		t.Error("unexpected tag 'nonexistent'")
	}
}

func TestStaleDetection(t *testing.T) {
	// Simulate docs with watches and updates-when, matched against changed files.
	docs := []document.Document{
		{
			Path: "/docs/reconciler.md",
			Frontmatter: document.Frontmatter{
				Title:       "Reconciler",
				Watches:     []string{"api/v1/types.go", "internal/reconciler/"},
				UpdatesWhen: []string{"CRD schema changes"},
			},
		},
		{
			Path: "/docs/deploy.md",
			Frontmatter: document.Frontmatter{
				Title:       "Deploy Guide",
				UpdatesWhen: []string{"new environment added"},
			},
		},
		{
			Path: "/docs/readme.md",
			Frontmatter: document.Frontmatter{
				Title: "README",
			},
		},
		{
			Path: "/docs/api.md",
			Frontmatter: document.Frontmatter{
				Title:   "API Reference",
				Watches: []string{"api/v1/*.go"},
			},
		},
	}

	changed := []string{
		"api/v1/types.go",
		"internal/reconciler/pool.go",
		"cmd/main.go",
	}

	var results []staleResult
	for _, doc := range docs {
		var matchedWatches []string
		for _, w := range doc.Frontmatter.Watches {
			for _, c := range changed {
				if matchWatch(w, c) {
					matchedWatches = append(matchedWatches, w)
					break
				}
			}
		}

		hasWatchHit := len(matchedWatches) > 0
		hasUpdatesWhen := len(doc.Frontmatter.UpdatesWhen) > 0

		if !hasWatchHit && !hasUpdatesWhen {
			continue
		}

		r := staleResult{
			Path:  doc.Path,
			Title: doc.Frontmatter.Title,
		}
		if hasWatchHit {
			r.Watches = matchedWatches
		}
		if hasUpdatesWhen {
			r.UpdatesWhen = doc.Frontmatter.UpdatesWhen
		}
		results = append(results, r)
	}

	// Expect: reconciler (watches hit + updates-when), deploy (updates-when only), api (watches hit).
	// readme has neither → excluded.
	if len(results) != 3 {
		t.Fatalf("got %d results, want 3", len(results))
	}

	// reconciler.md: both watches should match.
	r := results[0]
	if r.Path != "/docs/reconciler.md" {
		t.Errorf("result[0] path = %s, want /docs/reconciler.md", r.Path)
	}
	if len(r.Watches) != 2 {
		t.Errorf("reconciler matched watches = %d, want 2", len(r.Watches))
	}
	if len(r.UpdatesWhen) != 1 {
		t.Errorf("reconciler updates-when = %d, want 1", len(r.UpdatesWhen))
	}

	// deploy.md: updates-when only, no watches.
	r = results[1]
	if r.Path != "/docs/deploy.md" {
		t.Errorf("result[1] path = %s, want /docs/deploy.md", r.Path)
	}
	if r.Watches != nil {
		t.Errorf("deploy should have no matched watches")
	}
	if len(r.UpdatesWhen) != 1 {
		t.Errorf("deploy updates-when = %d, want 1", len(r.UpdatesWhen))
	}

	// api.md: glob watch matches.
	r = results[2]
	if r.Path != "/docs/api.md" {
		t.Errorf("result[2] path = %s, want /docs/api.md", r.Path)
	}
	if len(r.Watches) != 1 {
		t.Errorf("api matched watches = %d, want 1", len(r.Watches))
	}
}
