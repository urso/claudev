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
	Artifacts  []Artifact `json:"artifacts" jsonschema:"description=discovered documentation, config, and operational artifacts found recursively"`
}

type Artifact struct {
	Type string `json:"type" jsonschema:"description=artifact category: docs, readme, helm, k8s, playbook, scripts, ci"`
	Path string `json:"path" jsonschema:"description=relative path from repo root"`
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

// artifactDirs maps directory basenames to artifact types (searched recursively).
var artifactDirs = map[string]string{
	"docs":          "docs",
	"doc":           "docs",
	"documentation": "docs",
	"helm":          "helm",
	"charts":        "helm",
	"k8s":           "k8s",
	"kubernetes":    "k8s",
	"manifests":     "k8s",
	"playbooks":     "playbook",
	"runbooks":      "playbook",
	"scripts":       "scripts",
	"tools":         "scripts",
	"bin":           "scripts",
}

// artifactFiles maps file basenames (case-insensitive) to artifact types.
var artifactFiles = map[string]string{
	"readme.md":       "readme",
	"readme":          "readme",
	"readme.txt":      "readme",
	"architecture.md": "docs",
	"contributing.md": "docs",
	"changelog.md":    "docs",
	"jenkinsfile":     "ci",
	".gitlab-ci.yml":  "ci",
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
	var artifacts []Artifact
	seenArtifactDirs := map[string]struct{}{}

	for entry, werr := range walkFilesAndDirs(abs) {
		if werr != nil {
			return RepoSignals{}, werr
		}
		if entry.D.IsDir() {
			baseLower := strings.ToLower(entry.D.Name())
			if atype, ok := artifactDirs[baseLower]; ok {
				if _, seen := seenArtifactDirs[entry.Rel]; !seen {
					seenArtifactDirs[entry.Rel] = struct{}{}
					artifacts = append(artifacts, Artifact{Type: atype, Path: entry.Rel})
				}
			}
			continue
		}
		fileCount++
		extCounts[extOf(entry.Rel)]++
		dirCounts[topDirOfPath(entry.Rel)]++
		baseLower := strings.ToLower(filepath.Base(entry.Rel))
		if atype, ok := artifactFiles[baseLower]; ok {
			artifacts = append(artifacts, Artifact{Type: atype, Path: entry.Rel})
		}
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
		Artifacts:  artifacts,
	}, nil
}

type walkEntry struct {
	Rel string
	D   fs.DirEntry
}

// walkFilesAndDirs yields each non-ignored file/dir under root as (entry, error).
// On error, yields (walkEntry{}, err) once and stops. Pruned subtrees use fs.SkipDir.
func walkFilesAndDirs(root string) iter.Seq2[walkEntry, error] {
	return func(yield func(walkEntry, error) bool) {
		_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				if !yield(walkEntry{}, err) {
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
			}
			rel, rerr := filepath.Rel(root, path)
			if rerr != nil {
				return nil
			}
			if path == root {
				return nil
			}
			if !yield(walkEntry{Rel: rel, D: d}, nil) {
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
