package ticket

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/urso/claudev/project-discovery/internal/document"
)

// IDPrefix is the prefix used for all ticket ids.
const IDPrefix = "t-"

// minIDWidth is the minimum zero-padding width when serializing an ID.
// Ids never serialize narrower than this; wider numbers grow naturally.
const minIDWidth = 4

// ID is a ticket identifier. The on-disk form is "t-NNNN" (zero-padded to at
// least 4 digits) but in memory we carry it as an int so comparisons and
// next-id arithmetic are honest.
type ID int

// ParseID accepts any of "1", "01", "0001", "t-1", "t-0001" and returns the
// numeric id. Empty string, non-numeric, or non-positive values error.
func ParseID(s string) (ID, error) {
	raw := strings.TrimPrefix(s, IDPrefix)
	if raw == "" {
		return 0, fmt.Errorf("invalid ticket id %q", s)
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid ticket id %q: %w", s, err)
	}
	if n < 1 {
		return 0, fmt.Errorf("invalid ticket id %q: must be >= 1", s)
	}
	return ID(n), nil
}

// String returns the canonical on-disk form, zero-padded to at least minIDWidth.
func (id ID) String() string {
	return fmt.Sprintf("%s%0*d", IDPrefix, minIDWidth, int(id))
}

// NextID scans the tickets folder and returns the next sequential id.
// Files whose names do not parse as ticket ids are ignored.
func NextID(folder TicketFolder) (ID, error) {
	matches, err := filepath.Glob(filepath.Join(folder.Path(), IDPrefix+"*.md"))
	if err != nil {
		return 0, err
	}

	var max ID
	for _, m := range matches {
		base := filepath.Base(m)
		num, _, ok := strings.Cut(strings.TrimPrefix(base, IDPrefix), "-")
		if !ok {
			continue
		}
		id, err := ParseID(num)
		if err != nil {
			continue
		}
		if id > max {
			max = id
		}
	}
	return max + 1, nil
}

// findByID returns the path of the ticket whose frontmatter id equals want.
// It is the source-of-truth lookup: filename is only a hint, frontmatter is
// authoritative.
func findByID(folder TicketFolder, want ID) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(folder.Path(), IDPrefix+"*.md"))
	if err != nil {
		return nil, err
	}
	var hits []string
	for _, m := range matches {
		fm, err := document.ParseFrontmatterFromFile[Ticket](m)
		if err != nil {
			return nil, fmt.Errorf("read frontmatter %s: %w", m, err)
		}
		got, err := ParseID(fm.ID)
		if err != nil {
			continue
		}
		if got == want {
			hits = append(hits, m)
		}
	}
	return hits, nil
}
