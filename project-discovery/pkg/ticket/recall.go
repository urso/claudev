package ticket

// RecallResult contains the result of a recall search.
type RecallResult struct {
	Type    string      // "resolved", "active", "open", "overlap", "none"
	Matches []SearchHit // matched tickets (scored)
}

// Recall searches for existing knowledge before creating a new ticket.
// Order: done tickets → active/open tickets → overlap detection.
// Returns on first match. If no match, returns type "none".
func Recall(folder TicketFolder, intent string, scope string, tags []string) (RecallResult, error) {
	filter := SearchFilter{Tags: tags, Scope: scope}

	// Step 1: Search resolved (done) tickets
	hits, err := Search(folder, intent, filter)
	if err != nil {
		return RecallResult{}, err
	}
	resolved := filterByStatus(hits, "done")
	if len(resolved) > 0 {
		return RecallResult{Type: "resolved", Matches: resolved}, nil
	}

	// Step 2: Search active tickets
	active := filterByStatus(hits, "active")
	if len(active) > 0 {
		return RecallResult{Type: "active", Matches: active}, nil
	}

	// Step 3: Search open tickets
	open := filterByStatus(hits, "open")
	if len(open) > 0 {
		return RecallResult{Type: "open", Matches: open}, nil
	}

	// Step 4: Check overlap
	overlapHits, err := FindOverlap(folder, OverlapParams{Intent: intent, Scope: scope, Tags: tags})
	if err != nil {
		return RecallResult{}, err
	}
	if len(overlapHits) > 0 {
		matches := make([]SearchHit, len(overlapHits))
		for i, h := range overlapHits {
			matches[i] = SearchHit{Summary: h.Summary, Score: h.Score}
		}
		return RecallResult{Type: "overlap", Matches: matches}, nil
	}

	return RecallResult{Type: "none", Matches: nil}, nil
}

func filterByStatus(hits []SearchHit, status string) []SearchHit {
	var out []SearchHit
	for _, h := range hits {
		if h.Status == status {
			out = append(out, h)
		}
	}
	return out
}
