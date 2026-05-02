package ticket

import "testing"

func writeTicketFull(t *testing.T, f TicketFolder, p NewParams, body string) {
	t.Helper()
	doc, err := New(f, p)
	if err != nil {
		t.Fatal(err)
	}
	doc.Body = []byte(body)
	if err := f.Write(doc); err != nil {
		t.Fatal(err)
	}
}

func TestSearch_RanksTitleOverBody(t *testing.T) {
	f := newFolder(t)
	writeTicketFull(t, f, NewParams{Title: "Authentication flow", Intention: "investigate auth"}, "body mentions authentication once")
	writeTicketFull(t, f, NewParams{Title: "Database migration", Intention: "schema work"}, "and authentication is also mentioned briefly")
	writeTicketFull(t, f, NewParams{Title: "Logging pipeline"}, "no relevant words here")

	hits, err := Search(f, "authentication", SearchFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) < 2 {
		t.Fatalf("expected >=2 hits, got %d", len(hits))
	}
	if hits[0].ID != "t-0001" {
		t.Errorf("expected t-0001 first (title match), got %s", hits[0].ID)
	}
}

func TestSearch_StemmingMatches(t *testing.T) {
	f := newFolder(t)
	writeTicketFull(t, f, NewParams{Title: "Deployment pipeline"}, "")
	hits, err := Search(f, "deploying", SearchFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 {
		t.Fatalf("expected 1 hit (stem match), got %d", len(hits))
	}
}

func TestSearch_FilterTagsAndScope(t *testing.T) {
	f := newFolder(t)
	writeTicketFull(t, f, NewParams{Title: "Auth login", Tags: []string{"security"}, Scope: "internal/auth login flow"}, "")
	writeTicketFull(t, f, NewParams{Title: "Auth tokens", Tags: []string{"infra"}, Scope: "tokens storage"}, "")

	hits, err := Search(f, "auth", SearchFilter{Tags: []string{"security"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 || hits[0].ID != "t-0001" {
		t.Errorf("tag filter failed: %+v", hits)
	}

	hits, err = Search(f, "auth", SearchFilter{Scope: "logins"})
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 || hits[0].ID != "t-0001" {
		t.Errorf("scope stem filter failed: %+v", hits)
	}
}

func TestSearch_EmptyQuery(t *testing.T) {
	f := newFolder(t)
	writeTicketFull(t, f, NewParams{Title: "Anything"}, "")
	hits, err := Search(f, "", SearchFilter{})
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 0 {
		t.Errorf("empty query should return no hits, got %d", len(hits))
	}
}
