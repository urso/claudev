package linkgraph

import (
	"iter"
	"os"

	"github.com/urso/claudev/docnav/pkg/document"
	"github.com/urso/claudev/docnav/pkg/parser"
)

// Graph holds forward and reverse link edges between documents.
type Graph struct {
	// Forward maps a document path to its outgoing links.
	Forward map[string][]document.Link
	// Reverse maps a document path to paths of documents that link to it.
	Reverse map[string][]string
}

// NewGraph creates an empty link graph.
func NewGraph() *Graph {
	return &Graph{
		Forward: make(map[string][]document.Link),
		Reverse: make(map[string][]string),
	}
}

// AddDocument adds a document's links to the graph.
func (g *Graph) AddDocument(doc document.Document) {
	g.Forward[doc.Path] = doc.Links
	for _, link := range doc.Links {
		g.Reverse[link.Path] = append(g.Reverse[link.Path], doc.Path)
	}
}

// BrokenLink represents a link whose target does not exist.
type BrokenLink struct {
	Source string
	Link   document.Link
}

// Backlinks yields paths of documents that link to the given path.
func (g *Graph) Backlinks(path string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, src := range g.Reverse[path] {
			if !yield(src) {
				return
			}
		}
	}
}

// Orphans yields document paths that have no incoming links and are not hub files.
func (g *Graph) Orphans() iter.Seq[string] {
	return func(yield func(string) bool) {
		for path := range g.Forward {
			if document.IsHubFile(path) {
				continue
			}
			if len(g.Reverse[path]) == 0 {
				if !yield(path) {
					return
				}
			}
		}
	}
}

// BrokenLinks yields links whose target path does not exist on disk.
func (g *Graph) BrokenLinks() iter.Seq[BrokenLink] {
	return func(yield func(BrokenLink) bool) {
		for source, links := range g.Forward {
			for _, link := range links {
				if _, err := os.Stat(link.Path); err != nil {
					if !yield(BrokenLink{Source: source, Link: link}) {
						return
					}
				}
			}
		}
	}
}

// Build constructs a link graph by reading and parsing files from the given iterator.
func Build(files iter.Seq2[string, error]) (*Graph, error) {
	g := NewGraph()

	for path, err := range files {
		if err != nil {
			return nil, err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		doc := parser.ParseDocument(path, content)
		g.AddDocument(doc)
	}

	return g, nil
}
