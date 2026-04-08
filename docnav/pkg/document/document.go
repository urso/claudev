package document

import (
	"path/filepath"
	"strings"
)

// Frontmatter holds optional YAML frontmatter fields from a markdown document.
type Frontmatter struct {
	Title             string   `yaml:"title"`
	Type              string   `yaml:"type"`
	Tags              []string `yaml:"tags"`
	Scope             string   `yaml:"scope"`
	UpdatesWhen       []string `yaml:"updates-when"`
	Watches           []string `yaml:"watches"`
	AgentInstructions string         `yaml:"agent-instructions"`
	Extra             map[string]any `yaml:",inline"`
}

// Link represents a markdown link extracted from a document.
type Link struct {
	Text string
	Path string // resolved absolute path
}

// Document represents a parsed markdown file.
type Document struct {
	Path        string
	Frontmatter Frontmatter
	Links       []Link
	IsHub       bool
}

// hubFilenames are recognized as hub/entry-point files.
var hubFilenames = map[string]bool{
	"readme.md": true,
	"index.md":  true,
}

// IsHubFile reports whether the given file path is a hub/entry-point file.
func IsHubFile(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	return hubFilenames[base]
}
