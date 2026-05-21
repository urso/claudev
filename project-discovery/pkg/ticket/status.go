package ticket

import "fmt"

// StatusSummary contains ticket counts by status and lists of active/open tickets.
type StatusSummary struct {
	Counts map[string]int
	Active []Summary
	Open   []Summary
}

// Status returns a summary of tickets by status.
func Status(folder TicketFolder) (StatusSummary, error) {
	all, err := List(folder, ListFilter{})
	if err != nil {
		return StatusSummary{}, err
	}

	counts := map[string]int{
		StatusOpen:    0,
		StatusActive:  0,
		StatusDone:    0,
		StatusDropped: 0,
		StatusBlocked: 0,
	}

	var active, open []Summary
	for _, t := range all {
		counts[t.Status]++
		switch t.Status {
		case StatusActive:
			active = append(active, t)
		case StatusOpen:
			open = append(open, t)
		}
	}

	return StatusSummary{Counts: counts, Active: active, Open: open}, nil
}

// FormatStatus formats a StatusSummary for display.
func FormatStatus(s StatusSummary) string {
	out := fmt.Sprintf("open: %d | active: %d | done: %d | dropped: %d | blocked: %d\n",
		s.Counts[StatusOpen], s.Counts[StatusActive], s.Counts[StatusDone],
		s.Counts[StatusDropped], s.Counts[StatusBlocked])

	if len(s.Active) > 0 {
		out += "\nActive:\n"
		for _, t := range s.Active {
			out += fmt.Sprintf("  %s  %s\n", t.ID, t.Title)
		}
	}

	if len(s.Open) > 0 {
		out += "\nOpen:\n"
		for _, t := range s.Open {
			out += fmt.Sprintf("  %s  %s\n", t.ID, t.Title)
		}
	}

	return out
}
