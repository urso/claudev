package ticket

import (
	"testing"
)

func writeTicket(t *testing.T, f TicketFolder, p NewParams, status string) {
	t.Helper()
	doc, err := New(f, p)
	if err != nil {
		t.Fatal(err)
	}
	doc.Frontmatter.Status = status
	if err := f.Write(doc); err != nil {
		t.Fatal(err)
	}
}

func TestList_FilterAndSort(t *testing.T) {
	f := newFolder(t)
	writeTicket(t, f, NewParams{Title: "A", Tags: []string{"x"}}, StatusOpen)
	writeTicket(t, f, NewParams{Title: "B", Tags: []string{"y"}}, StatusActive)
	writeTicket(t, f, NewParams{Title: "C", Tags: []string{"x"}}, StatusDone)

	all, err := List(f, ListFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 3 {
		t.Fatalf("len = %d", len(all))
	}
	if all[0].ID != "t-0001" || all[2].ID != "t-0003" {
		t.Errorf("not sorted: %v", []string{all[0].ID, all[1].ID, all[2].ID})
	}

	open, err := List(f, ListFilter{Status: StatusOpen})
	if err != nil {
		t.Fatal(err)
	}
	if len(open) != 1 || open[0].Title != "A" {
		t.Errorf("status filter: %+v", open)
	}

	xs, err := List(f, ListFilter{Tag: "x"})
	if err != nil {
		t.Fatal(err)
	}
	if len(xs) != 2 {
		t.Errorf("tag filter len = %d", len(xs))
	}

	combo, err := List(f, ListFilter{Status: StatusDone, Tag: "x"})
	if err != nil {
		t.Fatal(err)
	}
	if len(combo) != 1 || combo[0].Title != "C" {
		t.Errorf("combo: %+v", combo)
	}
}
