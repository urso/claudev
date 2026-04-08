package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/urso/claudev/docnav/pkg/document"
	"github.com/urso/claudev/docnav/pkg/linkgraph"
	"github.com/urso/claudev/docnav/pkg/parser"
	"github.com/urso/claudev/docnav/pkg/related"
	"github.com/urso/claudev/docnav/pkg/search"
	"github.com/urso/claudev/docnav/pkg/walker"
)

type GlobalFlags struct {
	Path string
	JSON bool
}

func (g *GlobalFlags) Register(fs *flag.FlagSet) {
	fs.StringVar(&g.Path, "path", "", "root directory (default: git root, or cwd)")
	fs.BoolVar(&g.JSON, "json", false, "output as JSON lines")
}

// Root returns the resolved root path. If Path is empty, it tries the git
// repo root, falling back to the current working directory.
func (g *GlobalFlags) Root() (string, error) {
	if g.Path != "" {
		return filepath.Abs(g.Path)
	}
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err == nil {
		return strings.TrimSpace(string(out)), nil
	}
	return os.Getwd()
}

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var commands = map[string]func([]string) error{
	"search":       runSearch,
	"links":        runLinks,
	"backlinks":    runBacklinks,
	"orphans":      runOrphans,
	"broken-links": runBrokenLinks,
	"related":      runRelated,
	"tags":         runTags,
	"stale":        runStale,
}

func run(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: docnav <command> [flags]\ncommands: %s", commandList())
	}

	cmd, ok := commands[args[1]]
	if !ok {
		return fmt.Errorf("unknown command: %s\nknown commands: %s", args[1], commandList())
	}
	return cmd(args[2:])
}

func commandList() string {
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

func runSearch(args []string) error {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	var gf GlobalFlags
	gf.Register(fs)

	var typeFilter string
	var tagsFilter string
	fs.StringVar(&typeFilter, "type", "", "filter by document type")
	fs.StringVar(&tagsFilter, "tags", "", "filter by tags (comma-separated)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("usage: docnav search [flags] <query>")
	}
	query := fs.Arg(0)

	root, err := gf.Root()
	if err != nil {
		return err
	}

	var searchOpts []search.Option
	if typeFilter != "" {
		searchOpts = append(searchOpts, search.WithTypeFilter(typeFilter))
	}
	if tagsFilter != "" {
		searchOpts = append(searchOpts, search.WithTagsFilter(strings.Split(tagsFilter, ",")...))
	}

	w := walker.NewWalker(walker.MarkdownFiles())
	s := search.NewSearcher(searchOpts...)
	results, err := s.Search(w.Walk(root), query)
	if err != nil {
		return err
	}

	if gf.JSON {
		enc := json.NewEncoder(os.Stdout)
		for _, r := range results {
			enc.Encode(r)
		}
	} else {
		for _, r := range results {
			fmt.Fprintf(os.Stdout, "%s\t%s\t%.0f\t%s\n", r.Path, r.Title, r.Score, r.MatchedField)
		}
	}
	return nil
}

type linkResult struct {
	Source string `json:"source,omitempty"`
	Target string `json:"target"`
	Text   string `json:"text,omitempty"`
	Broken bool   `json:"broken,omitempty"`
}

func runLinks(args []string) error {
	fs := flag.NewFlagSet("links", flag.ContinueOnError)
	var gf GlobalFlags
	gf.Register(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: docnav links [flags] <file>")
	}

	file, err := filepath.Abs(fs.Arg(0))
	if err != nil {
		return err
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	doc := parser.ParseDocument(file, content)

	if gf.JSON {
		enc := json.NewEncoder(os.Stdout)
		for _, l := range doc.Links {
			_, statErr := os.Stat(l.Path)
			enc.Encode(linkResult{
				Target: l.Path,
				Text:   l.Text,
				Broken: statErr != nil,
			})
		}
	} else {
		for _, l := range doc.Links {
			_, statErr := os.Stat(l.Path)
			status := ""
			if statErr != nil {
				status = "\tBROKEN"
			}
			fmt.Fprintf(os.Stdout, "%s\t%s%s\n", l.Path, l.Text, status)
		}
	}
	return nil
}

func runBacklinks(args []string) error {
	fs := flag.NewFlagSet("backlinks", flag.ContinueOnError)
	var gf GlobalFlags
	gf.Register(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: docnav backlinks [flags] <file>")
	}

	file, err := filepath.Abs(fs.Arg(0))
	if err != nil {
		return err
	}

	root, err := gf.Root()
	if err != nil {
		return err
	}

	w := walker.NewWalker(walker.MarkdownFiles())
	g, err := linkgraph.Build(w.Walk(root))
	if err != nil {
		return err
	}

	if gf.JSON {
		enc := json.NewEncoder(os.Stdout)
		for src := range g.Backlinks(file) {
			enc.Encode(linkResult{Source: src, Target: file})
		}
	} else {
		for src := range g.Backlinks(file) {
			fmt.Fprintln(os.Stdout, src)
		}
	}
	return nil
}

func runOrphans(args []string) error {
	fs := flag.NewFlagSet("orphans", flag.ContinueOnError)
	var gf GlobalFlags
	gf.Register(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}

	root, err := gf.Root()
	if err != nil {
		return err
	}

	w := walker.NewWalker(walker.MarkdownFiles())
	g, err := linkgraph.Build(w.Walk(root))
	if err != nil {
		return err
	}

	if gf.JSON {
		enc := json.NewEncoder(os.Stdout)
		for path := range g.Orphans() {
			enc.Encode(struct {
				Path string `json:"path"`
			}{Path: path})
		}
	} else {
		for path := range g.Orphans() {
			fmt.Fprintln(os.Stdout, path)
		}
	}
	return nil
}

func runBrokenLinks(args []string) error {
	fs := flag.NewFlagSet("broken-links", flag.ContinueOnError)
	var gf GlobalFlags
	gf.Register(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}

	root, err := gf.Root()
	if err != nil {
		return err
	}

	w := walker.NewWalker(walker.MarkdownFiles())
	g, err := linkgraph.Build(w.Walk(root))
	if err != nil {
		return err
	}

	if gf.JSON {
		enc := json.NewEncoder(os.Stdout)
		for bl := range g.BrokenLinks() {
			enc.Encode(linkResult{
				Source: bl.Source,
				Target: bl.Link.Path,
				Text:   bl.Link.Text,
				Broken: true,
			})
		}
	} else {
		for bl := range g.BrokenLinks() {
			fmt.Fprintf(os.Stdout, "%s\t%s\t%s\n", bl.Source, bl.Link.Text, bl.Link.Path)
		}
	}
	return nil
}

func runRelated(args []string) error {
	fs := flag.NewFlagSet("related", flag.ContinueOnError)
	var gf GlobalFlags
	gf.Register(fs)
	var topN int
	fs.IntVar(&topN, "top", 10, "number of results to return")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("usage: docnav related [flags] <file>")
	}

	file, err := filepath.Abs(fs.Arg(0))
	if err != nil {
		return err
	}

	root, err := gf.Root()
	if err != nil {
		return err
	}

	w := walker.NewWalker(walker.MarkdownFiles())
	g, docs, contents, err := linkgraph.BuildWithDocs(w.Walk(root))
	if err != nil {
		return err
	}

	// Find target document.
	var target *document.Document
	for i := range docs {
		if docs[i].Path == file {
			target = &docs[i]
			break
		}
	}
	if target == nil {
		return fmt.Errorf("file not found in document set: %s", file)
	}

	results := related.Related(*target, g, docs, contents, topN)

	if gf.JSON {
		enc := json.NewEncoder(os.Stdout)
		for _, r := range results {
			enc.Encode(r)
		}
	} else {
		for _, r := range results {
			fmt.Fprintf(os.Stdout, "%s\t%.4f\tco-citation=%d tags=%d watches=%d content=%.3f\n",
				r.Path, r.Score, r.CoCitation, r.SharedTags, r.SharedWatches, r.ContentSim)
		}
	}
	return nil
}

// scanDocs walks the root and parses all markdown documents.
func scanDocs(root string) ([]document.Document, error) {
	w := walker.NewWalker(walker.MarkdownFiles())
	var docs []document.Document
	for path, err := range w.Walk(root) {
		if err != nil {
			return nil, err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		docs = append(docs, parser.ParseDocument(path, content))
	}
	return docs, nil
}

func runTags(args []string) error {
	fs := flag.NewFlagSet("tags", flag.ContinueOnError)
	var gf GlobalFlags
	gf.Register(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}

	root, err := gf.Root()
	if err != nil {
		return err
	}

	docs, err := scanDocs(root)
	if err != nil {
		return err
	}

	// Build tag → paths map.
	tagDocs := make(map[string][]string)
	for _, doc := range docs {
		for _, tag := range doc.Frontmatter.Tags {
			tagDocs[tag] = append(tagDocs[tag], doc.Path)
		}
	}

	// If a tag argument is given, filter to that tag.
	if fs.NArg() > 0 {
		tag := fs.Arg(0)
		paths := tagDocs[tag]
		if len(paths) == 0 {
			return nil
		}
		sort.Strings(paths)
		if gf.JSON {
			enc := json.NewEncoder(os.Stdout)
			for _, p := range paths {
				enc.Encode(struct {
					Tag  string `json:"tag"`
					Path string `json:"path"`
				}{Tag: tag, Path: p})
			}
		} else {
			for _, p := range paths {
				fmt.Fprintln(os.Stdout, p)
			}
		}
		return nil
	}

	// List all tags sorted by name.
	tags := make([]string, 0, len(tagDocs))
	for tag := range tagDocs {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	if gf.JSON {
		enc := json.NewEncoder(os.Stdout)
		for _, tag := range tags {
			enc.Encode(struct {
				Tag   string `json:"tag"`
				Count int    `json:"count"`
			}{Tag: tag, Count: len(tagDocs[tag])})
		}
	} else {
		for _, tag := range tags {
			fmt.Fprintf(os.Stdout, "%s\t%d\n", tag, len(tagDocs[tag]))
		}
	}
	return nil
}

// matchWatch reports whether a changed file matches a watch pattern.
// Supports exact paths, directory prefixes (ending in /), and glob patterns.
// All paths are relative to the git root.
func matchWatch(watch, changed string) bool {
	// Directory prefix match.
	if strings.HasSuffix(watch, "/") {
		return strings.HasPrefix(changed, watch)
	}
	// Glob match — handles exact paths, single-star, and ** patterns.
	matched, _ := doublestar.Match(watch, changed)
	return matched
}

// changedFiles returns file paths changed since the given commit.
func changedFiles(since string) ([]string, error) {
	out, err := exec.Command("git", "diff", "--name-only", "--end-of-options", since+"..HEAD").Output()
	if err != nil {
		return nil, fmt.Errorf("git diff: %w", err)
	}
	var files []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

type staleResult struct {
	Path        string   `json:"path"`
	Title       string   `json:"title,omitempty"`
	Watches     []string `json:"matched_watches,omitempty"`
	UpdatesWhen []string `json:"updates_when,omitempty"`
}

func runStale(args []string) error {
	fs := flag.NewFlagSet("stale", flag.ContinueOnError)
	var gf GlobalFlags
	gf.Register(fs)
	var since string
	fs.StringVar(&since, "since", "HEAD~1", "git commit to diff against (default: HEAD~1)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	root, err := gf.Root()
	if err != nil {
		return err
	}

	changed, err := changedFiles(since)
	if err != nil {
		return err
	}

	docs, err := scanDocs(root)
	if err != nil {
		return err
	}

	var results []staleResult
	for _, doc := range docs {
		var matchedWatches []string
		for _, w := range doc.Frontmatter.Watches {
			for _, c := range changed {
				if matchWatch(w, c) {
					matchedWatches = append(matchedWatches, w)
					break
				}
			}
		}

		hasWatchHit := len(matchedWatches) > 0
		hasUpdatesWhen := len(doc.Frontmatter.UpdatesWhen) > 0

		if !hasWatchHit && !hasUpdatesWhen {
			continue
		}

		r := staleResult{
			Path:  doc.Path,
			Title: doc.Frontmatter.Title,
		}
		if hasWatchHit {
			r.Watches = matchedWatches
		}
		if hasUpdatesWhen {
			r.UpdatesWhen = doc.Frontmatter.UpdatesWhen
		}
		results = append(results, r)
	}

	if gf.JSON {
		enc := json.NewEncoder(os.Stdout)
		for _, r := range results {
			enc.Encode(r)
		}
	} else {
		for _, r := range results {
			fmt.Fprintf(os.Stdout, "%s\t%s\n", r.Path, r.Title)
			for _, w := range r.Watches {
				fmt.Fprintf(os.Stdout, "  watch: %s\n", w)
			}
			for _, u := range r.UpdatesWhen {
				fmt.Fprintf(os.Stdout, "  updates-when: %s\n", u)
			}
		}
	}
	return nil
}
