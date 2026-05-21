package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urso/claudev/project-discovery/pkg/scan"
	"github.com/urso/claudev/project-discovery/pkg/ticket"
)

type scanCmd struct {
	Repo   scanRepoCmd   `cmd:"" help:"scan repo facts (raw signals)"`
	Churn  scanChurnCmd  `cmd:"" help:"scan file churn and code ownership"`
	Deps   scanDepsCmd   `cmd:"" help:"deferred — inspect manifests directly"`
	Schema scanSchemaCmd `cmd:"" help:"print annotated JSON shape of scan outputs"`
}

type scanRepoCmd struct {
	Root   string `help:"repo root to scan" default:"."`
	TopExt int    `name:"top-ext" help:"cap on extensions reported" default:"25"`
}

func (c scanRepoCmd) Run() error {
	sig, err := scan.Repo(c.Root, scan.RepoOptions{TopExtensions: c.TopExt})
	if err != nil {
		return err
	}
	payload := map[string]any{"repo": sig}
	buf, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	if err := writeDiscoveryFile(c.Root, "scan.json", buf); err != nil {
		return err
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

type scanChurnCmd struct {
	Root            string `help:"repo root to scan" default:"."`
	Since           string `help:"window (e.g. 90d, 6m, 1y)" default:"90d"`
	TopFiles        int    `name:"top-files" help:"top N files" default:"25"`
	TopContributors int    `name:"top-contributors" help:"top N contributors" default:"25"`
	AuthorsPerDir   int    `name:"authors-per-dir" help:"top N authors per top-level dir" default:"3"`
	LargeMinLines   int    `name:"large-min-lines" help:"min lines for large-commit list" default:"500"`
}

func (c scanChurnCmd) Run() error {
	sig, err := scan.Churn(c.Root, scan.ChurnOptions{
		Since:               c.Since,
		TopFilesN:           c.TopFiles,
		TopContributorsN:    c.TopContributors,
		TopAuthorsPerDirN:   c.AuthorsPerDir,
		LargeCommitMinLines: c.LargeMinLines,
	})
	if err != nil {
		if errors.Is(err, scan.ErrNotGitRepo) {
			fmt.Fprintln(os.Stderr, "scan churn: not a git repository — skipping")
			return nil
		}
		return err
	}
	buf, err := json.MarshalIndent(map[string]any{"churn": sig}, "", "  ")
	if err != nil {
		return err
	}
	if _, err := os.Stdout.Write(buf); err != nil {
		return err
	}
	fmt.Println()
	if sig.Note != "" {
		fmt.Fprintln(os.Stderr, "scan churn: "+sig.Note)
	}
	return nil
}

type scanDepsCmd struct{}

func (scanDepsCmd) Run() error {
	return errors.New("scan deps: not implemented — inspect manifests directly (package.json, go.mod, Cargo.toml, pyproject.toml, ...)")
}

type scanSchemaCmd struct{}

func (scanSchemaCmd) Run() error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	return enc.Encode(scan.Schema())
}

// writeDiscoveryFile writes data to <discovery-root>/<name>, locating or
// creating the .discovery/ directory under root.
func writeDiscoveryFile(root, name string, data []byte) error {
	d, err := ticket.FindDiscovery(root)
	var dir string
	if err != nil {
		if !errors.Is(err, ticket.ErrDiscoveryDirNotFound) {
			return err
		}
		abs, aerr := filepath.Abs(root)
		if aerr != nil {
			return aerr
		}
		dir = filepath.Join(abs, ticket.DiscoveryDirName)
	} else {
		dir = d.Path()
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, name), data, 0o644)
}
