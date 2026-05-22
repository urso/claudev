package ticket

import (
	"path/filepath"
	"slices"
	"sort"

	"github.com/urso/claudev/project-discovery/internal/document"
)

// ListFilter narrows the result of List. Empty fields match everything.
type ListFilter struct {
	Status string
	Tag    string
}

// Summary is the per-ticket row returned by List. Body bytes are not loaded.
type Summary struct {
	ID     string
	Status string
	Title  string
	Tags   []string
	Path   string
}

// List enumerates tickets in folder, reading only their frontmatter, and
// returns rows matching filter. Results are sorted by id ascending.
func List(folder TicketFolder, filter ListFilter) ([]Summary, error) {
	matches, err := filepath.Glob(filepath.Join(folder.Path(), IDPrefix+"*.md"))
	if err != nil {
		return nil, err
	}
	out := make([]Summary, 0, len(matches))
	for _, path := range matches {
		fm, err := document.ParseFrontmatterFromFile[Ticket](path)
		if err != nil {
			return nil, err
		}
		if filter.Status != "" && fm.Status != filter.Status {
			continue
		}
		if filter.Tag != "" && !slices.Contains(fm.Tags, filter.Tag) {
			continue
		}
		out = append(out, Summary{
			ID:     fm.ID,
			Status: fm.Status,
			Title:  fm.Title,
			Tags:   fm.Tags,
			Path:   path,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// Children returns tickets that have parentID in their Parents field.
func Children(folder TicketFolder, parentID string) ([]Summary, error) {
	matches, err := filepath.Glob(filepath.Join(folder.Path(), IDPrefix+"*.md"))
	if err != nil {
		return nil, err
	}
	out := make([]Summary, 0)
	for _, path := range matches {
		fm, err := document.ParseFrontmatterFromFile[Ticket](path)
		if err != nil {
			return nil, err
		}
		if !slices.Contains(fm.Parents, parentID) {
			continue
		}
		out = append(out, Summary{
			ID:     fm.ID,
			Status: fm.Status,
			Title:  fm.Title,
			Tags:   fm.Tags,
			Path:   path,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
