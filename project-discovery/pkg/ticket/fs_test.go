package ticket

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/urso/claudev/project-discovery/internal/document"
)

func TestResolve_FrontmatterAuthoritative(t *testing.T) {
	f := newFolder(t)
	writeTicket(t, f, NewParams{Title: "A"}, StatusOpen) // t-0001

	// Lookup by various input forms — all should find the same file via frontmatter.
	doc, err := f.Read(ID(1))
	if err != nil {
		t.Fatal(err)
	}
	if doc.Frontmatter.Title != "A" {
		t.Errorf("title = %q", doc.Frontmatter.Title)
	}
}

func TestResolve_NotFound(t *testing.T) {
	f := newFolder(t)
	_, err := f.Read(ID(99))
	if !errors.Is(err, ErrTicketNotFound) {
		t.Errorf("err = %v, want ErrTicketNotFound", err)
	}
}

func TestResolve_IgnoresMismatchedFrontmatter(t *testing.T) {
	f := newFolder(t)
	// File named t-0005-foo.md but frontmatter id says t-0007.
	// Lookup of ID(5) should NOT find it; ID(7) should.
	body := []byte("---\nid: t-0007\ntitle: T\n---\nbody\n")
	if err := os.WriteFile(filepath.Join(f.Path(), "t-0005-foo.md"), body, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := f.Read(ID(5)); !errors.Is(err, ErrTicketNotFound) {
		t.Errorf("ID(5): err = %v, want ErrTicketNotFound", err)
	}
	if _, err := f.Read(ID(7)); err != nil {
		t.Errorf("ID(7): unexpected err = %v", err)
	}
}

func TestWrite_NoPathErrors(t *testing.T) {
	f := newFolder(t)
	doc := Document{} // no Path
	if err := f.Write(doc); !errors.Is(err, document.ErrNoPath) {
		t.Errorf("err = %v, want ErrNoPath", err)
	}
}

func TestResolve_DuplicateFrontmatterIDsError(t *testing.T) {
	f := newFolder(t)
	body := []byte("---\nid: t-0001\ntitle: T\n---\n")
	if err := os.WriteFile(filepath.Join(f.Path(), "t-0001-a.md"), body, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(f.Path(), "t-0001-b.md"), body, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := f.Read(ID(1)); err == nil {
		t.Fatal("want error for duplicate ids")
	}
}
