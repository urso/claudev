---
description: >-
  Search, find, and navigate markdown documentation using the docnav CLI. Use this skill
  whenever the user wants to find docs about a topic ("what docs do we have about X",
  "find documentation about Y", "is there a doc for Z"), list docs by type or tag
  ("show me all design docs", "find docs tagged auth"), explore how docs connect
  ("what links to this doc", "what does this doc link to", "what's related to X"),
  or needs to locate relevant documentation before implementing, fixing, or understanding
  code. This skill wraps a dedicated doc search and navigation tool — always prefer it
  over manual grep/find when the user is looking for project documentation.
user-invocable: true
argument-hint: "<query or file path>"
allowed-tools: ["Bash", "Read", "Glob"]
---

# doc-nav — Find and navigate project documentation

Discover, search, and explore markdown documentation in the current repository using the `docnav` CLI.

## Variables

- **DOCNAV**: `${CLAUDE_PLUGIN_ROOT}/scripts/docnav.sh`

## Available Commands

| Command | Input | Purpose |
|---------|-------|---------|
| `search <query>` | free-text query | Fuzzy search across filenames, frontmatter, headings, and content |
| `related <file>` | absolute file path | Find docs related by shared tags, links, watches, and content similarity |
| `links <file>` | absolute file path | Show outgoing markdown links from a doc (flags broken ones) |
| `backlinks <file>` | absolute file path | Show all docs that link to a given file |

All commands accept `--path <dir>` to set the root directory (defaults to git root) and `--json` for structured output.

### Search flags

- `--type <type>` — filter by document type (e.g., design, runbook, note)
- `--tags <t1,t2>` — filter by tags (comma-separated)

When `--type` or `--tags` is set, the query is optional. This lets you list all docs of a given type or tag without a search term:
```bash
bash DOCNAV search --type design --path /path/to/docs
bash DOCNAV search --tags auth --path /path/to/docs
```

### Related flags

- `--top <n>` — number of results (default: 10)

## How to Choose Commands

Given user input `$ARGUMENTS`:

1. **Free-text query** (e.g., "reconciler architecture", "how does auth work"):
   Run `search` and present the results. Do not automatically follow up with `related`, `links`, or `backlinks` unless the user explicitly asks to explore connections or related docs.

2. **File path** (e.g., `docs/design.md`, a path from search results):
   Run `links` and `backlinks` to map the document's connections. Only add `related` if the user asks for non-obvious connections beyond direct links.

3. **List by type or tag** (e.g., "show me all design docs", "what's tagged auth"):
   Run `search --type <type>` or `search --tags <tags>` without a query term.

4. **No input / browsing**:
   Run `search` with a broad query derived from the current working context, or suggest the user provide a topic.

## Running Commands

```bash
bash DOCNAV search "reconciler"
bash DOCNAV search --type design "auth"
bash DOCNAV search --tags operator,reconciler "pool"
bash DOCNAV related /absolute/path/to/doc.md
bash DOCNAV links /absolute/path/to/doc.md
bash DOCNAV backlinks /absolute/path/to/doc.md
```

## Output Formats

**search** — tab-separated: `path  title  score  matched_field`
**related** — tab-separated: `path  score  co-citation=N tags=N watches=N content=F`
**links** — tab-separated: `target_path  link_text  [BROKEN]`
**backlinks** — one source path per line

## Presenting Results

- Show results as a concise list with paths and titles
- For search results, mention which field matched (title, tags, content, etc.)
- For links, flag any broken links prominently
- If results are empty, suggest broadening the query or trying a different command
- When the user wants to read a doc from the results, use the Read tool on the path
- Keep responses minimal — present the docnav output in a readable format without reading or summarizing the docs unless the user asks for content
