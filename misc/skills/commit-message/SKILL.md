---
description: Generate commit message from staged changes
user-invocable: true
disable-model-invocation: false
argument-hint: "[--story story-file]"
allowed-tools: ["Read", "Bash(git:*)"]
model: haiku
effort: low
---

# Generate Commit Message

Generate a concise commit message based on staged changes. Does NOT auto-commit - outputs message for user to copy.

## User Input
```
$ARGUMENTS
```

Parse for:
- `--story <story-file>` for context (optional)

## Pre-computed Context

### Staged Changes Summary
!`git diff --cached --stat`

### Full Staged Diff
!`git diff --cached`

### Recent Commits (for style detection)
!`git log --oneline -10`

## Process

### 1. Check Staged Changes

If the staged changes summary above is empty, report and exit:
```
No staged changes. Stage files with `git add` first.
```

### 2. Detect Project Commit Style

From the recent commits above, detect if project uses:
- Conventional commits (`feat:`, `fix:`, `refactor:`, etc.)
- Free-form messages
- Other patterns (e.g., ticket prefixes like `[ABC-123]`)

Match the detected style.

### 3. Load Context (if provided)

If `--story` provided, read the file to understand what was being worked on.

### 4. Analyze Changes

From the diff, identify:
- What files changed
- What kind of change (new feature, bug fix, refactor, tests, docs)
- Key modifications

### 5. Generate Message

**Title (first line):**
- Max 50 characters
- Concise summary of change
- Match project style (conventional or free-form)
- Imperative mood ("Add feature" not "Added feature")

**Body (optional, if changes are complex):**
- Brief explanation of what and why
- If story provided, can mention: `Part of: <story name>`
- Blank line between title and body

### 6. Output

```
## Suggested Commit Message

```
<title>

<body if needed>
```

## Guidelines

- **Concise**: Title should be scannable in git log
- **Accurate**: Message should reflect actual changes, not intended changes
- **Match style**: Follow project's existing commit conventions
- **No auto-commit**: User decides when to commit
