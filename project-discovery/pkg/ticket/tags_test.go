package ticket

import "testing"

func TestTags_CountAndSort(t *testing.T) {
	f := newFolder(t)
	writeTicketFull(t, f, NewParams{Title: "A", Tags: []string{"x", "y"}}, "")
	writeTicketFull(t, f, NewParams{Title: "B", Tags: []string{"x"}}, "")
	writeTicketFull(t, f, NewParams{Title: "C", Tags: []string{"z"}}, "")

	got, err := Tags(f)
	if err != nil {
		t.Fatal(err)
	}
	// "discovery" + "ticket" come from DefaultTags on every ticket.
	want := map[string]int{"discovery": 3, "ticket": 3, "x": 2, "y": 1, "z": 1}
	have := make(map[string]int, len(got))
	for _, tc := range got {
		have[tc.Tag] = tc.Count
	}
	for k, v := range want {
		if have[k] != v {
			t.Errorf("tag %q count = %d, want %d", k, have[k], v)
		}
	}
	// Sort: count desc, then tag asc. discovery vs ticket both 3 → discovery first.
	if got[0].Tag != "discovery" || got[1].Tag != "ticket" {
		t.Errorf("sort: got top %s,%s", got[0].Tag, got[1].Tag)
	}
}
