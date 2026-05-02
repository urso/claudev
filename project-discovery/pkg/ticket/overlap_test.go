package ticket

import "testing"

func TestFindOverlap_TagAndScope(t *testing.T) {
	f := newFolder(t)
	writeTicketFull(t, f, NewParams{
		Title:     "Auth login flow",
		Tags:      []string{"security", "auth"},
		Scope:     "internal/auth login session",
		Intention: "understand login flow and session lifecycle",
	}, "login session token cookie")
	writeTicketFull(t, f, NewParams{
		Title:     "Database migrations",
		Tags:      []string{"infra"},
		Scope:     "schema migrations",
		Intention: "track schema migration history",
	}, "alembic schema versions")
	writeTicketFull(t, f, NewParams{
		Title:     "Logging",
		Tags:      []string{"infra"},
		Scope:     "logs",
		Intention: "central logging",
	}, "")

	hits, err := FindOverlap(f, OverlapParams{
		Intent: "investigate login session expiry",
		Scope:  "internal/auth session",
		Tags:   []string{"security"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) == 0 || hits[0].ID != "t-0001" {
		t.Errorf("expected t-0001 top, got %+v", hits)
	}
}

func TestFindOverlap_NoSignal(t *testing.T) {
	f := newFolder(t)
	writeTicketFull(t, f, NewParams{Title: "Unrelated"}, "")
	hits, err := FindOverlap(f, OverlapParams{Intent: "totally different topic", Scope: "elsewhere"})
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 0 {
		t.Errorf("expected no overlap, got %+v", hits)
	}
}
