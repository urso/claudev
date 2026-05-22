package ticket

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// NewParams describes the user-supplied inputs for creating a ticket.
//
// Body and Slug are optional. If Slug is empty it is derived from Title.
// Body is written verbatim; the caller (skill, template, agent) owns its content.
// ParentSlugs maps parent IDs to their full slugs for Obsidian wikilinks.
type NewParams struct {
	Title       string
	Slug        string
	Body        []byte
	Scope       string
	Tags        []string
	Parents     []string
	ParentSlugs map[string]string // id -> full filename slug (e.g. "t-0002" -> "t-0002-nvme-of-target")
	Refs        []string
	Intention   string
}

// DefaultTags are merged into every new ticket so docnav-style filters work
// out of the box.
var DefaultTags = []string{"discovery", "ticket"}

// slugPattern is the validation regex for an explicit Slug.
var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// New builds a ticket Document for the next available id in folder.
// The returned document carries Path = "<folder>/<id>-<slug>.md" so callers
// can persist it via folder.Write(doc) or doc.Write().
func New(folder TicketFolder, p NewParams) (Document, error) {
	if strings.TrimSpace(p.Title) == "" {
		return Document{}, fmt.Errorf("title is required")
	}
	s := p.Slug
	if s == "" {
		s = slug(p.Title)
	} else if !slugPattern.MatchString(s) {
		return Document{}, fmt.Errorf("invalid slug %q: want [a-z0-9-], no leading/trailing/double dashes", s)
	}
	id, err := NextID(folder)
	if err != nil {
		return Document{}, err
	}
	doc := build(id, p, time.Now())
	doc.Path = filepath.Join(folder.Path(), id.String()+"-"+s+".md")
	return doc, nil
}

func build(id ID, p NewParams, now time.Time) Document {
	date := now.Format("2006-01-02")
	ticketType := "discovery-ticket"
	if len(p.Refs) > 0 {
		ticketType = "doc-reference"
	}
	t := Ticket{
		Title:     p.Title,
		Type:      ticketType,
		Tags:      mergeTags(DefaultTags, p.Tags),
		Scope:     p.Scope,
		ID:        id.String(),
		Status:    StatusOpen,
		Parents:   p.Parents,
		Refs:      p.Refs,
		Intention: p.Intention,
		Created:   date,
		Updated:   date,
	}
	body := p.Body
	if len(p.Parents) > 0 {
		body = appendRelations(body, p.Parents, p.ParentSlugs)
	}
	return Document{Frontmatter: t, Body: body}
}

// appendRelations adds a Relations section with wikilinks for Obsidian graph.
func appendRelations(body []byte, parents []string, slugs map[string]string) []byte {
	var b strings.Builder
	if len(body) > 0 {
		b.Write(body)
		if body[len(body)-1] != '\n' {
			b.WriteByte('\n')
		}
		b.WriteByte('\n')
	}
	b.WriteString("## Relations\n\n")
	for _, p := range parents {
		link := p
		if slug, ok := slugs[p]; ok {
			link = slug
		}
		fmt.Fprintf(&b, "- Parent: [[%s]]\n", link)
	}
	return []byte(b.String())
}

func mergeTags(base, extra []string) []string {
	seen := make(map[string]struct{}, len(base)+len(extra))
	out := make([]string, 0, len(base)+len(extra))
	for _, t := range append(append([]string{}, base...), extra...) {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}
	return out
}

// slug lowercases the title and reduces it to [a-z0-9-] separated tokens,
// capped at 50 characters. Used as a fallback when no explicit slug is given.
func slug(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	b.Grow(len(s))
	prevDash := true
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevDash = false
		default:
			if !prevDash {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}
	out := strings.TrimRight(b.String(), "-")
	if len(out) > 50 {
		out = strings.TrimRight(out[:50], "-")
	}
	if out == "" {
		out = "ticket"
	}
	return out
}
