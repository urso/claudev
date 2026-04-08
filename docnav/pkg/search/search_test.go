package search

import (
	"iter"
	"os"
	"path/filepath"
	"testing"

	"github.com/urso/claudev/docnav/pkg/walker"
)

func setupTestdata(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "README.md"), `---
title: Project Overview
type: reference
tags: [overview, getting-started]
scope: Main entry point for project documentation
---

# Project Overview

This is the main hub file.
`)

	writeFile(t, filepath.Join(dir, "design.md"), `---
title: Architecture Design
type: design
tags: [architecture, go]
scope: System architecture and component design
---

# Architecture Design

## Components

The system has a walker and parser.
`)

	writeFile(t, filepath.Join(dir, "plain.md"), `# Just a Plain File

Some content about deployment and infrastructure.
No frontmatter here.
`)

	sub := filepath.Join(dir, "subdir")
	os.MkdirAll(sub, 0o755)
	writeFile(t, filepath.Join(sub, "notes.md"), `---
title: Development Notes
type: note
tags: [dev, go]
---

# Development Notes

## Walker Implementation

Notes about the walker.
`)

	return dir
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func walkDir(dir string) iter.Seq2[string, error] {
	return walker.NewWalker(walker.MarkdownFiles()).Walk(dir)
}

func TestSearch_FilenameMatch(t *testing.T) {
	dir := setupTestdata(t)
	results, err := NewSearcher().Search(walkDir(dir), "plain")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	if results[0].MatchedField != "filename" {
		t.Errorf("MatchedField = %q, want %q", results[0].MatchedField, "filename")
	}
}

func TestSearch_TitleMatch(t *testing.T) {
	dir := setupTestdata(t)
	results, err := NewSearcher().Search(walkDir(dir), "architecture")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	if results[0].Title != "Architecture Design" {
		t.Errorf("Title = %q, want %q", results[0].Title, "Architecture Design")
	}
	if results[0].MatchedField != "title" {
		t.Errorf("MatchedField = %q, want %q", results[0].MatchedField, "title")
	}
}

func TestSearch_TagsMatch(t *testing.T) {
	dir := setupTestdata(t)
	results, err := NewSearcher().Search(walkDir(dir), "getting-started")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	if results[0].MatchedField != "tags" {
		t.Errorf("MatchedField = %q, want %q", results[0].MatchedField, "tags")
	}
}

func TestSearch_HeadingsMatch(t *testing.T) {
	dir := setupTestdata(t)
	results, err := NewSearcher().Search(walkDir(dir), "walker implementation")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	if results[0].MatchedField != "headings" {
		t.Errorf("MatchedField = %q, want %q", results[0].MatchedField, "headings")
	}
}

func TestSearch_ContentMatch(t *testing.T) {
	dir := setupTestdata(t)
	results, err := NewSearcher().Search(walkDir(dir), "deployment")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	if results[0].MatchedField != "content" {
		t.Errorf("MatchedField = %q, want %q", results[0].MatchedField, "content")
	}
}

func TestSearch_HubPriority(t *testing.T) {
	dir := setupTestdata(t)
	results, err := NewSearcher().Search(walkDir(dir), "overview")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	if !results[0].IsHub {
		t.Error("expected hub file to be first result")
	}
}

func TestSearch_TypeFilter(t *testing.T) {
	dir := setupTestdata(t)
	results, err := NewSearcher(WithTypeFilter("design")).Search(walkDir(dir), "go")
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Title == "Development Notes" {
			t.Error("type=note should be filtered out when TypeFilter=design")
		}
	}
}

func TestSearch_TagsFilter(t *testing.T) {
	dir := setupTestdata(t)
	results, err := NewSearcher(WithTagsFilter("architecture")).Search(walkDir(dir), "go")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Title != "Architecture Design" {
		t.Errorf("Title = %q, want %q", results[0].Title, "Architecture Design")
	}
}

func TestSearch_NoFrontmatter(t *testing.T) {
	dir := setupTestdata(t)
	results, err := NewSearcher().Search(walkDir(dir), "plain")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected results for file with no frontmatter")
	}
	if results[0].Title != "plain.md" {
		t.Errorf("Title = %q, want %q", results[0].Title, "plain.md")
	}
}

func TestSearch_NoResults(t *testing.T) {
	dir := setupTestdata(t)
	results, err := NewSearcher().Search(walkDir(dir), "zzzznonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearch_ScoreOrder(t *testing.T) {
	dir := setupTestdata(t)
	results, err := NewSearcher().Search(walkDir(dir), "notes")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	if results[0].MatchedField != "filename" {
		t.Errorf("expected filename match first, got %q", results[0].MatchedField)
	}
}
