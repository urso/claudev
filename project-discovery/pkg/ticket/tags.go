package ticket

import (
	"path/filepath"
	"sort"

	"github.com/urso/claudev/project-discovery/internal/document"
)

// TagCount is one row in the tag census returned by Tags.
type TagCount struct {
	Tag   string
	Count int
}

// Tags scans every ticket frontmatter in folder and returns a tag census,
// sorted by count desc then tag asc.
func Tags(folder TicketFolder) ([]TagCount, error) {
	matches, err := filepath.Glob(filepath.Join(folder.Path(), IDPrefix+"*.md"))
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, path := range matches {
		fm, err := document.ParseFrontmatterFromFile[Ticket](path)
		if err != nil {
			return nil, err
		}
		for _, t := range fm.Tags {
			counts[t]++
		}
	}
	out := make([]TagCount, 0, len(counts))
	for tag, n := range counts {
		out = append(out, TagCount{Tag: tag, Count: n})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Tag < out[j].Tag
	})
	return out, nil
}
