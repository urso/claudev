package ticket

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/urso/claudev/project-discovery/internal/document"
)

func TestSlug(t *testing.T) {
	cases := map[string]string{
		"Hello World":            "hello-world",
		"  Multiple   Spaces  ":  "multiple-spaces",
		"Punct!?@#Mark":          "punct-mark",
		"":                       "ticket",
		"---":                    "ticket",
		strings.Repeat("a", 100): strings.Repeat("a", 50),
		"Auth/Session Tokens v2": "auth-session-tokens-v2",
	}
	for in, want := range cases {
		if got := slug(in); got != want {
			t.Errorf("slug(%q) = %q want %q", in, got, want)
		}
	}
}

func TestNew_FieldsAndDedup(t *testing.T) {
	f := newFolder(t)
	doc, err := New(f, NewParams{
		Title:     "First Ticket",
		Scope:     "pkg/foo",
		Tags:      []string{"discovery", "foo", "foo", " bar "},
		Parents:   []string{"t-0000"},
		Intention: "why",
	})
	if err != nil {
		t.Fatal(err)
	}
	fm := doc.Frontmatter
	if fm.ID != "t-0001" {
		t.Errorf("id = %q", fm.ID)
	}
	if fm.Status != StatusOpen {
		t.Errorf("status = %q", fm.Status)
	}
	if fm.Type != "discovery-ticket" {
		t.Errorf("type = %q", fm.Type)
	}
	want := []string{"discovery", "ticket", "foo", "bar"}
	if !slices.Equal(fm.Tags, want) {
		t.Errorf("tags = %v want %v", fm.Tags, want)
	}
	if fm.Created == "" || fm.Updated != fm.Created {
		t.Errorf("created/updated = %q/%q", fm.Created, fm.Updated)
	}
	wantPath := filepath.Join(f.Path(), "t-0001-first-ticket.md")
	if doc.Path != wantPath {
		t.Errorf("path = %q want %q", doc.Path, wantPath)
	}
}

func TestNew_RoundTrip(t *testing.T) {
	f := newFolder(t)
	doc, err := New(f, NewParams{Title: "Round Trip"})
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Write(doc); err != nil {
		t.Fatal(err)
	}
	parsed, err := document.ParseDocumentFile[Ticket](doc.Path)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Frontmatter.ID != doc.Frontmatter.ID {
		t.Errorf("id mismatch")
	}
	if parsed.Frontmatter.Title != "Round Trip" {
		t.Errorf("title mismatch")
	}
	if parsed.Path != doc.Path {
		t.Errorf("parsed.Path = %q want %q", parsed.Path, doc.Path)
	}
}

func TestNew_EmptyTitle(t *testing.T) {
	f := newFolder(t)
	if _, err := New(f, NewParams{Title: "  "}); err == nil {
		t.Fatal("want error for empty title")
	}
}

func TestNew_ExplicitSlugAndBody(t *testing.T) {
	f := newFolder(t)
	body := []byte("# anything\n\ncustom body\n")
	doc, err := New(f, NewParams{
		Title: "Some Long Title",
		Slug:  "auth-flow",
		Body:  body,
	})
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(doc.Path) != "t-0001-auth-flow.md" {
		t.Errorf("path basename = %q", filepath.Base(doc.Path))
	}
	if string(doc.Body) != string(body) {
		t.Errorf("body not preserved")
	}
}

func TestNew_InvalidSlug(t *testing.T) {
	f := newFolder(t)
	bad := []string{"Bad Slug", "-leading", "trailing-", "double--dash", "Up", "with_underscore"}
	for _, s := range bad {
		if _, err := New(f, NewParams{Title: "X", Slug: s}); err == nil {
			t.Errorf("slug %q: want error", s)
		}
	}
}

func TestNew_EmptyBody(t *testing.T) {
	f := newFolder(t)
	doc, err := New(f, NewParams{Title: "X"})
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Body) != 0 {
		t.Errorf("body = %q want empty", doc.Body)
	}
}
