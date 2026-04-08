package parser

import (
	"testing"
)

func TestParseFrontmatter_Valid(t *testing.T) {
	content := []byte(`---
title: Test Doc
type: design
tags: [go, cli]
scope: Testing frontmatter parsing
updates-when:
  - schema changes
watches:
  - pkg/parser/
agent-instructions: do the thing
---

# Body content
`)
	fm := ParseFrontmatter(content)

	if fm.Title != "Test Doc" {
		t.Errorf("Title = %q, want %q", fm.Title, "Test Doc")
	}
	if fm.Type != "design" {
		t.Errorf("Type = %q, want %q", fm.Type, "design")
	}
	if len(fm.Tags) != 2 || fm.Tags[0] != "go" || fm.Tags[1] != "cli" {
		t.Errorf("Tags = %v, want [go cli]", fm.Tags)
	}
	if fm.Scope != "Testing frontmatter parsing" {
		t.Errorf("Scope = %q, want %q", fm.Scope, "Testing frontmatter parsing")
	}
	if len(fm.UpdatesWhen) != 1 || fm.UpdatesWhen[0] != "schema changes" {
		t.Errorf("UpdatesWhen = %v, want [schema changes]", fm.UpdatesWhen)
	}
	if len(fm.Watches) != 1 || fm.Watches[0] != "pkg/parser/" {
		t.Errorf("Watches = %v, want [pkg/parser/]", fm.Watches)
	}
	if fm.AgentInstructions != "do the thing" {
		t.Errorf("AgentInstructions = %q, want %q", fm.AgentInstructions, "do the thing")
	}
}

func TestParseFrontmatter_Missing(t *testing.T) {
	content := []byte(`# No frontmatter here

Just a regular markdown file.
`)
	fm := ParseFrontmatter(content)
	if fm.Title != "" {
		t.Errorf("expected empty frontmatter, got Title=%q", fm.Title)
	}
}

func TestParseFrontmatter_Partial(t *testing.T) {
	content := []byte(`---
title: Only Title
---

# Content
`)
	fm := ParseFrontmatter(content)
	if fm.Title != "Only Title" {
		t.Errorf("Title = %q, want %q", fm.Title, "Only Title")
	}
	if fm.Type != "" {
		t.Errorf("Type should be empty, got %q", fm.Type)
	}
	if len(fm.Tags) != 0 {
		t.Errorf("Tags should be empty, got %v", fm.Tags)
	}
}

func TestParseFrontmatter_Malformed(t *testing.T) {
	content := []byte(`---
title: [invalid yaml
  broken: {
---
`)
	fm := ParseFrontmatter(content)
	if fm.Title != "" {
		t.Errorf("expected empty frontmatter on malformed YAML, got Title=%q", fm.Title)
	}
}

func TestParseFrontmatter_ExtraFields(t *testing.T) {
	content := []byte(`---
title: With Extras
custom-field: hello
priority: 5
---
`)
	fm := ParseFrontmatter(content)
	if fm.Title != "With Extras" {
		t.Errorf("Title = %q, want %q", fm.Title, "With Extras")
	}
	if fm.Extra["custom-field"] != "hello" {
		t.Errorf("Extra[custom-field] = %v, want %q", fm.Extra["custom-field"], "hello")
	}
}

func TestExtractLinks_Basic(t *testing.T) {
	content := []byte(`See [design](../design.md) and [notes](./notes.md) for details.`)
	links := ExtractLinks(content, "/repo/docs/guide")

	if len(links) != 2 {
		t.Fatalf("got %d links, want 2", len(links))
	}
	if links[0].Text != "design" || links[0].Path != "/repo/docs/design.md" {
		t.Errorf("link[0] = %+v, want {design /repo/docs/design.md}", links[0])
	}
	if links[1].Text != "notes" || links[1].Path != "/repo/docs/guide/notes.md" {
		t.Errorf("link[1] = %+v, want {notes /repo/docs/guide/notes.md}", links[1])
	}
}

func TestExtractLinks_IgnoresExternal(t *testing.T) {
	content := []byte(`Check [Google](https://google.com) and [HTTP](http://example.com).`)
	links := ExtractLinks(content, "/repo")

	if len(links) != 0 {
		t.Errorf("got %d links, want 0 (external URLs should be ignored)", len(links))
	}
}

func TestExtractLinks_StripAnchor(t *testing.T) {
	content := []byte(`See [section](doc.md#heading) for more.`)
	links := ExtractLinks(content, "/repo/docs")

	if len(links) != 1 {
		t.Fatalf("got %d links, want 1", len(links))
	}
	if links[0].Path != "/repo/docs/doc.md" {
		t.Errorf("Path = %q, want %q", links[0].Path, "/repo/docs/doc.md")
	}
}

func TestExtractLinks_AnchorOnly(t *testing.T) {
	content := []byte(`See [above](#section) for details.`)
	links := ExtractLinks(content, "/repo/docs")

	if len(links) != 0 {
		t.Errorf("got %d links, want 0 (anchor-only links should be skipped)", len(links))
	}
}

func TestExtractLinks_EmptyTarget(t *testing.T) {
	content := []byte(`An [empty]() link.`)
	links := ExtractLinks(content, "/repo")

	if len(links) != 0 {
		t.Errorf("got %d links, want 0", len(links))
	}
}

func TestExtractLinks_AbsolutePath(t *testing.T) {
	content := []byte(`See [abs](/docs/readme.md).`)
	links := ExtractLinks(content, "/repo/src")

	if len(links) != 1 {
		t.Fatalf("got %d links, want 1", len(links))
	}
	if links[0].Path != "/docs/readme.md" {
		t.Errorf("Path = %q, want %q", links[0].Path, "/docs/readme.md")
	}
}

func TestExtractLinks_MultipleOnSameLine(t *testing.T) {
	content := []byte(`[a](a.md) and [b](b.md) and [c](c.md)`)
	links := ExtractLinks(content, "/repo")

	if len(links) != 3 {
		t.Fatalf("got %d links, want 3", len(links))
	}
}

func TestExtractLinks_IgnoresCodeBlocks(t *testing.T) {
	content := []byte("Real [link](real.md) here.\n\n```\n[not a link](fake.md)\n```\n")
	links := ExtractLinks(content, "/repo")

	if len(links) != 1 {
		t.Fatalf("got %d links, want 1", len(links))
	}
	if links[0].Path != "/repo/real.md" {
		t.Errorf("Path = %q, want %q", links[0].Path, "/repo/real.md")
	}
}

func TestExtractLinks_IgnoresInlineCode(t *testing.T) {
	content := []byte("See `[not a link](fake.md)` and [real](real.md).")
	links := ExtractLinks(content, "/repo")

	if len(links) != 1 {
		t.Fatalf("got %d links, want 1", len(links))
	}
	if links[0].Text != "real" {
		t.Errorf("Text = %q, want %q", links[0].Text, "real")
	}
}

func TestParseDocument(t *testing.T) {
	content := []byte(`---
title: My Doc
tags: [test]
---

See [other](other.md) for details.
`)
	doc := ParseDocument("/repo/docs/README.md", content)

	if doc.Path != "/repo/docs/README.md" {
		t.Errorf("Path = %q", doc.Path)
	}
	if !doc.IsHub {
		t.Error("expected IsHub=true for README.md")
	}
	if doc.Frontmatter.Title != "My Doc" {
		t.Errorf("Title = %q", doc.Frontmatter.Title)
	}
	if len(doc.Links) != 1 || doc.Links[0].Path != "/repo/docs/other.md" {
		t.Errorf("Links = %+v", doc.Links)
	}
}
