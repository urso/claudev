package related

import (
	"testing"

	"github.com/urso/claudev/docnav/pkg/document"
	"github.com/urso/claudev/docnav/pkg/linkgraph"
)

func TestRelated_CoCitation(t *testing.T) {
	g := linkgraph.NewGraph()
	docs := []document.Document{
		{Path: "/a.md", Links: []document.Link{{Path: "/c.md"}}},
		{Path: "/b.md", Links: []document.Link{{Path: "/c.md"}}},
		{Path: "/c.md"},
	}
	for _, d := range docs {
		g.AddDocument(d)
	}

	results := Related(docs[0], g, docs, nil, 10)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Path != "/b.md" {
		t.Errorf("expected /b.md, got %s", results[0].Path)
	}
	if results[0].CoCitation != 1 {
		t.Errorf("expected co-citation=1, got %d", results[0].CoCitation)
	}
}

func TestRelated_SharedTags(t *testing.T) {
	g := linkgraph.NewGraph()
	docs := []document.Document{
		{Path: "/a.md", Frontmatter: document.Frontmatter{Tags: []string{"go", "cli"}}},
		{Path: "/b.md", Frontmatter: document.Frontmatter{Tags: []string{"go", "web"}}},
		{Path: "/c.md", Frontmatter: document.Frontmatter{Tags: []string{"python"}}},
	}
	for _, d := range docs {
		g.AddDocument(d)
	}

	results := Related(docs[0], g, docs, nil, 10)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d: %+v", len(results), results)
	}
	if results[0].Path != "/b.md" {
		t.Errorf("expected /b.md, got %s", results[0].Path)
	}
	if results[0].SharedTags != 1 {
		t.Errorf("expected shared_tags=1, got %d", results[0].SharedTags)
	}
}

func TestRelated_SharedWatches(t *testing.T) {
	g := linkgraph.NewGraph()
	docs := []document.Document{
		{Path: "/a.md", Frontmatter: document.Frontmatter{Watches: []string{"src/main.go", "pkg/"}}},
		{Path: "/b.md", Frontmatter: document.Frontmatter{Watches: []string{"src/main.go"}}},
	}
	for _, d := range docs {
		g.AddDocument(d)
	}

	results := Related(docs[0], g, docs, nil, 10)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].SharedWatches != 1 {
		t.Errorf("expected shared_watches=1, got %d", results[0].SharedWatches)
	}
}

func TestRelated_ExcludesDirectLinks(t *testing.T) {
	g := linkgraph.NewGraph()
	docs := []document.Document{
		{Path: "/a.md", Links: []document.Link{{Path: "/b.md"}}, Frontmatter: document.Frontmatter{Tags: []string{"go"}}},
		{Path: "/b.md", Frontmatter: document.Frontmatter{Tags: []string{"go"}}},
		{Path: "/c.md", Frontmatter: document.Frontmatter{Tags: []string{"go"}}},
	}
	for _, d := range docs {
		g.AddDocument(d)
	}

	results := Related(docs[0], g, docs, nil, 10)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d: %+v", len(results), results)
	}
	if results[0].Path != "/c.md" {
		t.Errorf("expected /c.md, got %s", results[0].Path)
	}
}

func TestRelated_ExcludesBacklinks(t *testing.T) {
	g := linkgraph.NewGraph()
	docs := []document.Document{
		{Path: "/a.md", Frontmatter: document.Frontmatter{Tags: []string{"go"}}},
		{Path: "/b.md", Links: []document.Link{{Path: "/a.md"}}, Frontmatter: document.Frontmatter{Tags: []string{"go"}}},
		{Path: "/c.md", Frontmatter: document.Frontmatter{Tags: []string{"go"}}},
	}
	for _, d := range docs {
		g.AddDocument(d)
	}

	results := Related(docs[0], g, docs, nil, 10)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d: %+v", len(results), results)
	}
	if results[0].Path != "/c.md" {
		t.Errorf("expected /c.md, got %s", results[0].Path)
	}
}

func TestRelated_TopN(t *testing.T) {
	g := linkgraph.NewGraph()
	docs := []document.Document{
		{Path: "/a.md", Frontmatter: document.Frontmatter{Tags: []string{"go"}}},
		{Path: "/b.md", Frontmatter: document.Frontmatter{Tags: []string{"go"}}},
		{Path: "/c.md", Frontmatter: document.Frontmatter{Tags: []string{"go"}}},
		{Path: "/d.md", Frontmatter: document.Frontmatter{Tags: []string{"go"}}},
	}
	for _, d := range docs {
		g.AddDocument(d)
	}

	results := Related(docs[0], g, docs, nil, 2)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestRelated_CombinedStructuralScore(t *testing.T) {
	// B has co-citation + shared tag with A. C has only shared tag.
	// B should rank higher.
	g := linkgraph.NewGraph()
	docs := []document.Document{
		{Path: "/a.md", Links: []document.Link{{Path: "/x.md"}}, Frontmatter: document.Frontmatter{Tags: []string{"go"}}},
		{Path: "/b.md", Links: []document.Link{{Path: "/x.md"}}, Frontmatter: document.Frontmatter{Tags: []string{"go"}}},
		{Path: "/c.md", Frontmatter: document.Frontmatter{Tags: []string{"go"}}},
		{Path: "/x.md"},
	}
	for _, d := range docs {
		g.AddDocument(d)
	}

	results := Related(docs[0], g, docs, nil, 10)
	if len(results) < 2 {
		t.Fatalf("expected at least 2 results, got %d", len(results))
	}
	if results[0].Path != "/b.md" {
		t.Errorf("expected /b.md first, got %s", results[0].Path)
	}
}

func TestRelated_NoResults(t *testing.T) {
	g := linkgraph.NewGraph()
	docs := []document.Document{
		{Path: "/a.md"},
		{Path: "/b.md"},
	}
	for _, d := range docs {
		g.AddDocument(d)
	}

	results := Related(docs[0], g, docs, nil, 10)
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestRelated_ContentSimilarity(t *testing.T) {
	g := linkgraph.NewGraph()
	docs := []document.Document{
		{Path: "/a.md"},
		{Path: "/b.md"},
		{Path: "/c.md"},
	}
	for _, d := range docs {
		g.AddDocument(d)
	}
	contents := map[string][]byte{
		"/a.md": []byte("# Kubernetes deployment strategies for production clusters"),
		"/b.md": []byte("# Kubernetes production cluster management and deployment"),
		"/c.md": []byte("# Cooking recipes for Italian pasta dishes"),
	}

	results := Related(docs[0], g, docs, contents, 10)
	if len(results) == 0 {
		t.Fatal("expected results from content similarity")
	}
	if results[0].Path != "/b.md" {
		t.Errorf("expected /b.md first (similar content), got %s", results[0].Path)
	}
	if results[0].ContentSim == 0 {
		t.Error("expected non-zero content similarity for /b.md")
	}
}

func TestRelated_ContentAndStructuralFusion(t *testing.T) {
	// B is related structurally (shared tag). C is related by content.
	// Both should appear in results.
	g := linkgraph.NewGraph()
	docs := []document.Document{
		{Path: "/a.md", Frontmatter: document.Frontmatter{Tags: []string{"infra"}}},
		{Path: "/b.md", Frontmatter: document.Frontmatter{Tags: []string{"infra"}}},
		{Path: "/c.md"},
	}
	for _, d := range docs {
		g.AddDocument(d)
	}
	contents := map[string][]byte{
		"/a.md": []byte("# Kubernetes deployment strategies"),
		"/b.md": []byte("# Completely different topic about cooking"),
		"/c.md": []byte("# Kubernetes deployment patterns and strategies"),
	}

	results := Related(docs[0], g, docs, contents, 10)
	if len(results) < 2 {
		t.Fatalf("expected at least 2 results, got %d", len(results))
	}
	foundB, foundC := false, false
	for _, r := range results {
		if r.Path == "/b.md" {
			foundB = true
		}
		if r.Path == "/c.md" {
			foundC = true
		}
	}
	if !foundB {
		t.Error("expected /b.md in results (structural match)")
	}
	if !foundC {
		t.Error("expected /c.md in results (content match)")
	}
}
