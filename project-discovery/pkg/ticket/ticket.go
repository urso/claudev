package ticket

import "github.com/urso/claudev/project-discovery/internal/document"

// Ticket holds the YAML frontmatter of a discovery ticket. The leading fields
// mirror docnav's frontmatter schema; the trailing fields are discovery-specific.
// Unknown keys are preserved via Extra so round-trip writes do not drop data.
type Ticket struct {
	Title       string   `yaml:"title"`
	Type        string   `yaml:"type"`
	Tags        []string `yaml:"tags,omitempty"`
	Scope       string   `yaml:"scope,omitempty"`
	Watches     []string `yaml:"watches,omitempty"`
	UpdatesWhen []string `yaml:"updates-when,omitempty"`

	ID        string   `yaml:"id"`
	Status    string   `yaml:"status"`
	Parents   []string `yaml:"parents,omitempty"`
	Related   []string `yaml:"related,omitempty"`
	Refs      []string `yaml:"refs,omitempty"`
	Intention string   `yaml:"intention,omitempty"`
	Created   string   `yaml:"created"`
	Updated   string   `yaml:"updated,omitempty"`

	Extra map[string]any `yaml:",inline"`
}

// Document is a parsed ticket file: frontmatter plus body bytes.
type Document = document.Document[Ticket]

// Ticket lifecycle status values.
const (
	StatusOpen    = "open"
	StatusActive  = "active"
	StatusDone    = "done"
	StatusDropped = "dropped"
	StatusBlocked = "blocked"
)
