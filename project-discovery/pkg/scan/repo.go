package scan

import (
	"cmp"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// RepoSignals is the raw output of a single tree walk. No classification —
// agent infers languages/frameworks from markers + extensions.
type RepoSignals struct {
	Root       string     `json:"root" jsonschema:"description=absolute path of the scanned root"`
	FileCount  int        `json:"file_count" jsonschema:"description=total non-ignored files walked"`
	Extensions []ExtCount `json:"extensions" jsonschema:"description=top N file extensions by count (capped by --top-ext, default 25)"`
	TopDirs    []DirCount `json:"top_dirs" jsonschema:"description=top-level directory inventory with file counts; '(root)' bucket = files directly in repo root; sorted by count desc"`
	Markers    []string   `json:"markers" jsonschema:"description=presence list of well-known marker files/dirs at repo root; dirs end with '/'; use to infer build system / ecosystem"`
}

type ExtCount struct {
	Ext   string `json:"ext" jsonschema:"description=lowercased file extension including leading '.'; '(none)' for extensionless"`
	Count int    `json:"count"`
}

type DirCount struct {
	Dir   string `json:"dir" jsonschema:"description=top-level directory name; '(root)' for files directly in repo root"`
	Files int    `json:"files"`
}

// RepoOptions tunes Repo.
type RepoOptions struct {
	TopExtensions int // 0 -> default 25
}

// ignoreDirs are directory basenames pruned during walk.
var ignoreDirs = map[string]struct{}{
	".git":         {},
	"node_modules": {},
	"vendor":       {},
	".bin":         {},
	"dist":         {},
	"build":        {},
	"target":       {},
	"out":          {},
	".next":        {},
	".cache":       {},
	".discovery":   {},
}

// fileMarkers are exact-name marker files at any depth (we only check root).
var fileMarkers = []string{
	"go.mod", "go.work",
	"package.json", "pnpm-workspace.yaml",
	"Cargo.toml",
	"pyproject.toml", "requirements.txt",
	"Makefile",
	"Dockerfile",
	"lerna.json",
	"build.gradle",
	"pom.xml",
}

// fileMarkerGlobs are glob patterns checked at root.
var fileMarkerGlobs = []string{
	"docker-compose.*",
}

// dirMarkers are directory paths (relative to root) whose presence is a marker.
var dirMarkers = []string{
	".github/workflows",
	"k8s",
	"terraform",
}

// Repo walks root once and returns raw signals.
func Repo(root string, opts RepoOptions) (RepoSignals, error) {
	if opts.TopExtensions <= 0 {
		opts.TopExtensions = 25
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return RepoSignals{}, err
	}

	extCounts := map[string]int{}
	dirCounts := map[string]int{}
	fileCount := 0

	for rel, werr := range walkFiles(abs) {
		if werr != nil {
			return RepoSignals{}, werr
		}
		fileCount++
		extCounts[extOf(rel)]++
		dirCounts[topDirOfPath(rel)]++
	}

	exts := make([]ExtCount, 0, len(extCounts))
	for k, v := range extCounts {
		exts = append(exts, ExtCount{Ext: k, Count: v})
	}
	slices.SortFunc(exts, func(a, b ExtCount) int {
		return cmp.Or(cmp.Compare(b.Count, a.Count), cmp.Compare(a.Ext, b.Ext))
	})
	if len(exts) > opts.TopExtensions {
		exts = exts[:opts.TopExtensions]
	}

	dirs := make([]DirCount, 0, len(dirCounts))
	for k, v := range dirCounts {
		dirs = append(dirs, DirCount{Dir: k, Files: v})
	}
	slices.SortFunc(dirs, func(a, b DirCount) int {
		return cmp.Or(cmp.Compare(b.Files, a.Files), cmp.Compare(a.Dir, b.Dir))
	})

	markers := detectMarkers(abs)

	return RepoSignals{
		Root:       abs,
		FileCount:  fileCount,
		Extensions: exts,
		TopDirs:    dirs,
		Markers:    markers,
	}, nil
}

// walkFiles yields each non-ignored file under root as a relative path.
// On error, yields ("", err) once and stops. Pruned subtrees use fs.SkipDir.
func walkFiles(root string) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				if !yield("", err) {
					return filepath.SkipAll
				}
				return err
			}
			if d.IsDir() {
				if path != root {
					if _, skip := ignoreDirs[d.Name()]; skip {
						return fs.SkipDir
					}
				}
				return nil
			}
			rel, rerr := filepath.Rel(root, path)
			if rerr != nil {
				return nil
			}
			if !yield(rel, nil) {
				return filepath.SkipAll
			}
			return nil
		})
	}
}

func extOf(rel string) string {
	ext := strings.ToLower(filepath.Ext(rel))
	if ext == "" {
		return "(none)"
	}
	return ext
}

func topDirOfPath(rel string) string {
	head, _, ok := strings.Cut(rel, string(filepath.Separator))
	if !ok {
		return "(root)"
	}
	return head
}

func detectMarkers(root string) []string {
	var out []string
	for _, name := range fileMarkers {
		if fi, err := os.Stat(filepath.Join(root, name)); err == nil && !fi.IsDir() {
			out = append(out, name)
		}
	}
	entries, err := os.ReadDir(root)
	if err == nil {
		for _, pat := range fileMarkerGlobs {
			for _, e := range entries {
				if e.IsDir() {
					continue
				}
				if ok, _ := filepath.Match(pat, e.Name()); ok {
					out = append(out, e.Name())
				}
			}
		}
	}
	for _, dir := range dirMarkers {
		if fi, err := os.Stat(filepath.Join(root, dir)); err == nil && fi.IsDir() {
			out = append(out, dir+"/")
		}
	}
	slices.Sort(out)
	return out
}
