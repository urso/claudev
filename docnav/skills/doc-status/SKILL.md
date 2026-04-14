---
description: >-
  Check documentation health: orphaned docs, broken links, and tag overview. Use this skill
  when the user asks about doc health ("are there orphan docs", "any broken links",
  "doc status", "documentation health check", "show me all tags", "tag summary",
  "lint docs", "check docs", "doc audit", "documentation issues"),
  wants to audit documentation quality, or needs to find docs that are disconnected
  or have linking issues. This skill wraps docnav health commands — always prefer it
  over manual searching when the user wants a documentation quality overview.
user-invocable: true
argument-hint: ""
allowed-tools: ["Bash", "Read", "Glob"]
---

# doc-status — Documentation health overview

Check documentation health by finding orphaned docs, broken links, and summarizing tags.

## Variables

- **DOCNAV**: `${CLAUDE_PLUGIN_ROOT}/scripts/docnav.sh`

## Available Commands

| Command | Purpose |
|---------|---------|
| `orphans` | Docs with no incoming links from other docs |
| `broken-links` | Links pointing to non-existent files |
| `tags` | List all tags with doc counts |

All commands accept `--path <dir>` to set the root directory (defaults to git root) and `--json` for structured output.

## Workflow

Run all three commands to build a combined health report:

```bash
bash DOCNAV orphans --path <root>
bash DOCNAV broken-links --path <root>
bash DOCNAV tags --path <root>
```

## Presenting Results

Present a single combined report with three sections:

1. **Broken Links** — list each broken link with the source file and target path. This is the most actionable issue — flag prominently.
2. **Orphaned Docs** — list docs with no incoming links. These may need to be linked from a hub or index page, or may be obsolete.
3. **Tag Summary** — show tags and their counts as a quick overview of documentation coverage.

If any section has zero results, say so briefly (e.g., "No broken links found"). Do not run follow-up commands (search, links, etc.) unless the user asks to investigate a specific issue.
