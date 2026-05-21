package ticket

import (
	"path/filepath"
	"sort"

	"github.com/urso/claudev/project-discovery/internal/document"
)

// NextSuggestion is a ticket suggestion with reason.
type NextSuggestion struct {
	Summary
	Reason string
	Depth  int
}

// Next returns up to n tickets to work on, ordered by heuristic:
// 1. Active tickets (resume in-progress)
// 2. Open tickets by parent depth (shallowest first)
func Next(folder TicketFolder, n int) ([]NextSuggestion, error) {
	all, err := List(folder, ListFilter{})
	if err != nil {
		return nil, err
	}

	// Build parent lookup
	idToParents := make(map[string][]string)
	for _, t := range all {
		doc, err := document.ParseFrontmatterFromFile[Ticket](t.Path)
		if err != nil {
			return nil, err
		}
		idToParents[t.ID] = doc.Parents
	}

	var suggestions []NextSuggestion

	// Active tickets first
	for _, t := range all {
		if t.Status == StatusActive {
			suggestions = append(suggestions, NextSuggestion{
				Summary: t,
				Reason:  "active, resume in-progress",
				Depth:   0,
			})
		}
	}

	// Open tickets by parent depth
	type openTicket struct {
		Summary
		depth int
	}
	var openTickets []openTicket
	for _, t := range all {
		if t.Status == StatusOpen {
			depth := computeDepth(t.ID, idToParents, make(map[string]bool))
			openTickets = append(openTickets, openTicket{Summary: t, depth: depth})
		}
	}
	sort.Slice(openTickets, func(i, j int) bool {
		if openTickets[i].depth != openTickets[j].depth {
			return openTickets[i].depth < openTickets[j].depth
		}
		return openTickets[i].ID < openTickets[j].ID
	})

	for _, t := range openTickets {
		reason := "open, no parents"
		if t.depth > 0 {
			reason = "open, child ticket"
		}
		suggestions = append(suggestions, NextSuggestion{
			Summary: t.Summary,
			Reason:  reason,
			Depth:   t.depth,
		})
	}

	if len(suggestions) > n {
		suggestions = suggestions[:n]
	}
	return suggestions, nil
}

func computeDepth(id string, parents map[string][]string, seen map[string]bool) int {
	if seen[id] {
		return 0 // cycle protection
	}
	seen[id] = true

	ps := parents[id]
	if len(ps) == 0 {
		return 0
	}

	maxParentDepth := 0
	for _, p := range ps {
		d := computeDepth(p, parents, seen)
		if d > maxParentDepth {
			maxParentDepth = d
		}
	}
	return maxParentDepth + 1
}

// FindTicketByGlob finds a ticket file matching the id pattern.
func FindTicketByGlob(folder TicketFolder, id string) (string, error) {
	pattern := filepath.Join(folder.Path(), id+"*.md")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", nil
	}
	return matches[0], nil
}
