package ticket

import (
	"strings"
	"testing"
)

func TestUpdateStatus(t *testing.T) {
	f := newFolder(t)
	writeTicket(t, f, NewParams{Title: "X"}, StatusOpen)

	if err := UpdateStatus(f, ID(1), StatusActive); err != nil {
		t.Fatal(err)
	}
	doc, err := f.Read(ID(1))
	if err != nil {
		t.Fatal(err)
	}
	if doc.Frontmatter.Status != StatusActive {
		t.Errorf("status = %q", doc.Frontmatter.Status)
	}
	if doc.Frontmatter.Updated == "" {
		t.Errorf("updated empty")
	}
}

func TestUpdateStatus_Invalid(t *testing.T) {
	f := newFolder(t)
	writeTicket(t, f, NewParams{Title: "X"}, StatusOpen)
	if err := UpdateStatus(f, ID(1), "bogus"); err == nil {
		t.Fatal("want error")
	}
}

func TestAppendLogs_ExistingSection(t *testing.T) {
	f := newFolder(t)
	writeTicket(t, f, NewParams{Title: "X"}, StatusOpen)
	if err := AppendLogs(f, ID(1), []string{"first", "second"}); err != nil {
		t.Fatal(err)
	}
	doc, err := f.Read(ID(1))
	if err != nil {
		t.Fatal(err)
	}
	body := string(doc.Body)
	if !strings.Contains(body, ": first") || !strings.Contains(body, ": second") {
		t.Errorf("missing entries:\n%s", body)
	}
	if strings.Count(body, "## Log") != 1 {
		t.Errorf("expected single Log section:\n%s", body)
	}
}

func TestAppendLogs_MissingSection(t *testing.T) {
	f := newFolder(t)
	doc, err := New(f, NewParams{Title: "X"})
	if err != nil {
		t.Fatal(err)
	}
	doc.Body = []byte("# X\n\n## Notes\n\nstuff\n")
	if err := f.Write(doc); err != nil {
		t.Fatal(err)
	}
	if err := AppendLogs(f, ID(1), []string{"hello"}); err != nil {
		t.Fatal(err)
	}
	got, err := f.Read(ID(1))
	if err != nil {
		t.Fatal(err)
	}
	body := string(got.Body)
	if !strings.Contains(body, "## Log") || !strings.Contains(body, ": hello") {
		t.Errorf("missing log section:\n%s", body)
	}
}

func TestAppendLogs_Empty(t *testing.T) {
	f := newFolder(t)
	writeTicket(t, f, NewParams{Title: "X"}, StatusOpen)
	before, _ := f.Read(ID(1))
	if err := AppendLogs(f, ID(1), nil); err != nil {
		t.Fatal(err)
	}
	after, _ := f.Read(ID(1))
	if string(before.Body) != string(after.Body) {
		t.Error("body changed for empty msgs")
	}
}

func TestAppendLogs_PreservesFollowingSection(t *testing.T) {
	f := newFolder(t)
	doc, err := New(f, NewParams{Title: "X"})
	if err != nil {
		t.Fatal(err)
	}
	doc.Body = []byte("# X\n\n## Log\n\n- 2026-01-01: old\n\n## After\n\ntail\n")
	if err := f.Write(doc); err != nil {
		t.Fatal(err)
	}
	if err := AppendLogs(f, ID(1), []string{"new"}); err != nil {
		t.Fatal(err)
	}
	got, _ := f.Read(ID(1))
	body := string(got.Body)
	if !strings.Contains(body, "old") || !strings.Contains(body, ": new") {
		t.Errorf("entries missing:\n%s", body)
	}
	if !strings.Contains(body, "## After") || !strings.Contains(body, "tail") {
		t.Errorf("after section lost:\n%s", body)
	}
	if strings.Index(body, "## Log") > strings.Index(body, "## After") {
		t.Errorf("section order wrong:\n%s", body)
	}
}
