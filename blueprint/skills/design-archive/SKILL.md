---
description: Archive a done design and optionally its stories
user-invocable: true
disable-model-invocation: false
argument-hint: "[design-id or file]"
allowed-tools: ["Read", "Bash", "AskUserQuestion"]
---

# Archive Design

Move a completed design and optionally its stories to the archive directory.

## Variables

- **DISCOVERY_GUIDE**: `${CLAUDE_PLUGIN_ROOT}/resources/discovery.md`
- **DESIGN_OPS**: `${CLAUDE_PLUGIN_ROOT}/resources/design-operations.md`
- **STORY_OPS**: `${CLAUDE_PLUGIN_ROOT}/resources/story-operations.md`
- **QUERY_STORIES**: `query-stories.sh`
- **ARCHIVE_DOC**: `archive-doc.sh`

## User Input
```
$ARGUMENTS
```

## Process

### 1. Load Guides

Read DISCOVERY_GUIDE, DESIGN_OPS, and STORY_OPS for available tools and procedures.

### 2. Archive

- Locate the design based on user input (prefer `--status done` designs when no specific input given)
- If status is not `done`, warn the user before proceeding
- Show associated stories and their statuses via `bash QUERY_STORIES --design <design-id>`
- Ask user what to archive: design only, design and all stories, or design and specific stories
- If any selected stories are not `done`, warn the user
- Use ARCHIVE_DOC to set status and move files:
    - `bash ARCHIVE_DOC design <id>` — archive a single design
    - `bash ARCHIVE_DOC design <id> --with-stories` — archive design and all its stories
    - `bash ARCHIVE_DOC design <id> --stories-only` — archive only the stories
- Report what was archived

## Guidelines

- Archiving is non-destructive — files are moved, not deleted
