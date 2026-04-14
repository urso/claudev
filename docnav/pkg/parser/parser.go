package parser

import (
	"path/filepath"
	"strings"

	"github.com/urso/claudev/docnav/pkg/document"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	gmparser "github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
)

var md = goldmark.New(
	goldmark.WithExtensions(&frontmatter.Extender{}),
)

// ParseDocument parses a markdown file from its content and path,
// returning a Document with frontmatter, links, and hub status.
// The content is parsed once to extract both frontmatter and links.
func ParseDocument(path string, content []byte) document.Document {
	ctx := gmparser.NewContext()
	reader := text.NewReader(content)
	tree := md.Parser().Parse(reader, gmparser.WithContext(ctx))

	doc := document.Document{
		Path:  path,
		IsHub: document.IsHubFile(path),
	}

	// Extract frontmatter from parser context.
	if data := frontmatter.Get(ctx); data != nil {
		_ = data.Decode(&doc.Frontmatter)
	}

	// Extract links and headings from AST.
	doc.Links = extractLinksFromTree(tree, content, filepath.Dir(path))
	doc.Headings = extractHeadingsFromTree(tree, content)
	return doc
}

// ParseFrontmatter extracts and parses YAML frontmatter from markdown content.
// Returns an empty Frontmatter if none is found or if parsing fails.
func ParseFrontmatter(content []byte) document.Frontmatter {
	var fm document.Frontmatter

	ctx := gmparser.NewContext()
	reader := text.NewReader(content)
	md.Parser().Parse(reader, gmparser.WithContext(ctx))

	data := frontmatter.Get(ctx)
	if data == nil {
		return fm
	}
	_ = data.Decode(&fm)
	return fm
}

// ExtractLinks extracts markdown links from content by walking the goldmark AST.
// External URLs (http/https) are ignored. Relative paths are resolved against fileDir.
func ExtractLinks(content []byte, fileDir string) []document.Link {
	reader := text.NewReader(content)
	tree := md.Parser().Parse(reader)
	return extractLinksFromTree(tree, content, fileDir)
}

func extractHeadingsFromTree(tree ast.Node, source []byte) []string {
	var headings []string
	ast.Walk(tree, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if _, ok := n.(*ast.Heading); !ok {
			return ast.WalkContinue, nil
		}
		var parts []string
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			if t, ok := c.(*ast.Text); ok {
				parts = append(parts, string(t.Value(source)))
			}
		}
		if text := strings.Join(parts, ""); text != "" {
			headings = append(headings, text)
		}
		return ast.WalkSkipChildren, nil
	})
	return headings
}

func extractLinksFromTree(tree ast.Node, source []byte, fileDir string) []document.Link {
	var links []document.Link
	ast.Walk(tree, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		link, ok := n.(*ast.Link)
		if !ok {
			return ast.WalkContinue, nil
		}

		target := string(link.Destination)

		// Strip anchor fragments.
		if i := strings.Index(target, "#"); i >= 0 {
			target = target[:i]
		}
		if target == "" {
			return ast.WalkContinue, nil
		}

		// Skip external URLs.
		lower := strings.ToLower(target)
		if strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://") {
			return ast.WalkContinue, nil
		}

		resolved := target
		if !filepath.IsAbs(target) {
			resolved = filepath.Join(fileDir, target)
		}
		resolved = filepath.Clean(resolved)

		// Collect link text from child text nodes.
		var textParts []string
		for c := link.FirstChild(); c != nil; c = c.NextSibling() {
			if t, ok := c.(*ast.Text); ok {
				textParts = append(textParts, string(t.Value(source)))
			}
		}

		links = append(links, document.Link{
			Text: strings.Join(textParts, ""),
			Path: resolved,
		})
		return ast.WalkSkipChildren, nil
	})
	return links
}
