---
description: Load design context for review or discussion (no implementation)
user-invocable: true
disable-model-invocation: true
argument-hint: "[design name/id]"
allowed-tools: ["Read", "Bash", "Glob", "Grep"]
---

# Prime Design Context

Load design context for code review, discussion, or planning. Does **not** implement anything.

## Variables

- **RESOLVE_DIR**: `${CLAUDE_PLUGIN_ROOT}/scripts/resolve-dir.sh`
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

### 1. Find Design

Use FIND_DOC to locate the design:

```bash
bash FIND_DOC design "<user input>"
```

If no input provided, list available designs and let user pick one.

Read the selected design document.

### 2. Load Design Dependencies

If the design has a `depends-on:` field:
- Use FIND_DOC to locate and read each dependency design

If the design has a `references:` field:
- Read any files listed (paths relative to git root)

### 3. Load Applicable Rules

Use the rules listing from the pre-computed context above.

Read and apply:
- All rules where `applies-to` includes `*` or matches the project's languages/technologies
- Style rules (tagged `style`) for coding conventions
- Build rules (tagged `build`) for build/lint/test commands

### 4. Summarize Loaded Context

Provide a brief summary:
- Design title and status
- Key decisions and constraints
- Dependent designs loaded (if any)
- Number of rules loaded by type

State that context is loaded and ready for review or discussion.

## Key Guidelines

- **Read-only**: Do not implement, modify code, or update any documents
- **Context sharing**: Goal is to share knowledge for review discussions
- **Concise summary**: Don't dump entire files, summarize what's loaded
- **Ready for review**: Context is now loaded for code review or planning conversations
