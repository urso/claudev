---
description: Load story and design context for review or discussion (no implementation)
user-invocable: true
disable-model-invocation: true
argument-hint: "[story or design name/id]"
allowed-tools: ["Read", "Bash", "Glob", "Grep"]
---

# Prime Story Context

Load story and design context for code review, discussion, or planning. Does **not** implement anything.

## Variables

- **RESOLVE_DIR**: `${CLAUDE_PLUGIN_ROOT}/scripts/resolve-dir.sh`
- **QUERY_WORK**: `${CLAUDE_PLUGIN_ROOT}/scripts/query-work.sh`
- **FIND_DOC**: `${CLAUDE_PLUGIN_ROOT}/scripts/find-doc.sh`
- **LIST_RULES**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-rules.sh`

## User Input
```
$ARGUMENTS
```

## Pre-computed Context

### All Rules
!`bash ${CLAUDE_PLUGIN_ROOT}/scripts/list-rules.sh`

Output format (grouped by directory, pipe-separated columns):
```
<directory>:
full/path/to/file.md|name|applies-to|tags|paths|description
```

## Process

### 1. Resolve Directories

```bash
STORIES_DIR=$(bash RESOLVE_DIR story)
```

### 2. Find Story

Find work using `query-work.sh`:

```bash
bash QUERY_WORK --search "<user input>"
```

If no input provided, omit `--search` to list all work.

Output format: `TYPE | ID | FILENAME | TITLE | STATUS | BLOCKED-BY | DESCRIPTION`

Present matching stories and let user pick one. Read the selected story document.

### 3. Load Design Context

If the story has a `design:` frontmatter field:
1. Use FIND_DOC to locate and read the parent design
2. Read any designs listed in the parent design's `depends-on` field
3. Read any files listed in the parent design's `references` field (paths relative to git root)

If the story has a `blocked-by:` field, read the Developer Logs sections (Decision Log, Blockers, Deviations, Lessons Learned) from those dependency stories.

### 4. Load Applicable Rules

Use the rules listing from the pre-computed context above.

Read and apply:
- All rules where `applies-to` includes `*` or matches the project's languages/technologies
- Style rules (tagged `style`) for coding conventions
- Build rules (tagged `build`) for build/lint/test commands

### 5. Summarize Loaded Context

Provide a brief summary:
- Story title, status, and task overview
- Parent design (if any) and key decisions
- Number of rules loaded by type
- Dependency stories referenced (if any)

State that context is loaded and ready for review or discussion.

## Key Guidelines

- **Read-only**: Do not implement, modify code, or update the story document
- **Context sharing**: Goal is to share knowledge for review discussions
- **Concise summary**: Don't dump entire files, summarize what's loaded
- **Ready for review**: Context is now loaded for code review or planning conversations
