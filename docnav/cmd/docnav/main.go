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

	"github.com/urso/claudev/docnav/pkg/linkgraph"
	"github.com/urso/claudev/docnav/pkg/parser"
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
	"search":    runSearch,
	"links":     runLinks,
	"backlinks": runBacklinks,
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
