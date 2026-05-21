package ticket

import (
	"os"
	"path/filepath"
	"testing"
)

func newFolder(t *testing.T) TicketFolder {
	t.Helper()
	dir := t.TempDir()
	return TicketFolder{dir: dir}
}

func touch(t *testing.T, dir, name string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte("---\ntitle: x\n---\n"), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestNextID_Empty(t *testing.T) {
	f := newFolder(t)
	id, err := NextID(f)
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatalf("got %v want 1", id)
	}
	if id.String() != "t-0001" {
		t.Fatalf("String() = %q want t-0001", id.String())
	}
}

func TestNextID_GapsAndIgnored(t *testing.T) {
	f := newFolder(t)
	touch(t, f.Path(), "t-0001-a.md")
	touch(t, f.Path(), "t-0003-b.md")
	touch(t, f.Path(), "readme.md")
	touch(t, f.Path(), "t-bogus.md")
	id, err := NextID(f)
	if err != nil {
		t.Fatal(err)
	}
	if id != 4 {
		t.Fatalf("got %v want 4", id)
	}
}

func TestParseID(t *testing.T) {
	cases := map[string]ID{
		"1":      1,
		"01":     1,
		"0001":   1,
		"t-1":    1,
		"t-0001": 1,
		"t-42":   42,
	}
	for in, want := range cases {
		got, err := ParseID(in)
		if err != nil {
			t.Errorf("ParseID(%q) err: %v", in, err)
			continue
		}
		if got != want {
			t.Errorf("ParseID(%q) = %v want %v", in, got, want)
		}
	}
	bad := []string{"", "t-", "abc", "t-abc", "0", "-1", "t--1"}
	for _, in := range bad {
		if _, err := ParseID(in); err == nil {
			t.Errorf("ParseID(%q): want error", in)
		}
	}
}

func TestID_String_WidthGrowsAbove4(t *testing.T) {
	if got := ID(99999).String(); got != "t-99999" {
		t.Errorf("got %q want t-99999", got)
	}
	if got := ID(7).String(); got != "t-0007" {
		t.Errorf("got %q want t-0007", got)
	}
}
