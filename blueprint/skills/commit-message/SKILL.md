---
description: Generate commit message from staged changes
user-invocable: true
disable-model-invocation: false
argument-hint: "[--story story-file]"
allowed-tools: ["Read", "Bash(git:*)"]
model: sonnet
---

# Generate Commit Message

Generate a concise commit message based on staged changes. Does NOT auto-commit - outputs message for user to copy.

## User Input
```
$ARGUMENTS
```

Parse for:
- `--story <story-file>` for context (optional)

## Process

### 1. Check Staged Changes

```bash
git diff --cached --stat
```

If nothing staged, report and exit:
```
No staged changes. Stage files with `git add` first.
```

### 2. Get Diff Content

```bash
git diff --cached
```

### 3. Detect Project Commit Style

Look at recent commits:
```bash
git log --oneline -10
```

Detect if project uses:
- Conventional commits (`feat:`, `fix:`, `refactor:`, etc.)
- Free-form messages
- Other patterns (e.g., ticket prefixes like `[ABC-123]`)

Match the detected style.

### 4. Load Context (if provided)

If `--story` provided, read the file to understand what was being worked on.

### 5. Analyze Changes

From the diff, identify:
- What files changed
- What kind of change (new feature, bug fix, refactor, tests, docs)
- Key modifications

### 6. Generate Message

**Title (first line):**
- Max 50 characters
- Concise summary of change
- Match project style (conventional or free-form)
- Imperative mood ("Add feature" not "Added feature")

**Body (optional, if changes are complex):**
- Brief explanation of what and why
- If story provided, can mention: `Part of: <story name>`
- Blank line between title and body

### 7. Output

```
## Suggested Commit Message

```
<title>

<body if needed>
```

Copy and use with:
  git commit -m "<title>" -m "<body>"

Or:
  git commit
  # Then paste into editor
```

## Guidelines

- **Concise**: Title should be scannable in git log
- **Accurate**: Message should reflect actual changes, not intended changes
- **Match style**: Follow project's existing commit conventions
- **No auto-commit**: User decides when to commit
