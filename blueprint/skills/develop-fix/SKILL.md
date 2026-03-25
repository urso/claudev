---
description: Fix build and lint errors
user-invocable: true
disable-model-invocation: true
argument-hint: "[--story story-file]"
model: sonnet
allowed-tools: ["Read", "Write", "Edit", "Glob", "Grep", "Bash", "AskUserQuestion"]
---

# Fix Build/Lint Errors

Fix build and lint errors using project build configuration.

## Variables

- **LIST_RULES**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-rules.sh`
- **QUERY_STORIES**: `${CLAUDE_PLUGIN_ROOT}/scripts/query-stories.sh`

## User Input
```
$ARGUMENTS
```

Parse for:
- `--story <story-file>` for context (optional)

## Process

### 1. Load Build Rules

```bash
bash LIST_RULES "" "" build
```

Output format (pipe-separated): `filename|name|applies-to|tags|paths|description`

Read build rules to get project commands (build, lint, test, format).

If no build rules exist, ask user for the commands or try common defaults.

### 2. Load Story Context (if provided)

If `--story` specified, read the story document to understand what's being implemented.

### 3. Run Build/Lint

Execute the build and/or lint commands from the build rules. Capture all errors and warnings.

If build/lint passes with no errors, report success and exit.

### 4. Load Style Rules (only if errors found)

```bash
bash LIST_RULES "" "" style
```

Read applicable style rules for context on expected patterns when fixing.

### 5. Analyze and Fix Errors

For each error/warning:
- Identify the file and line
- Understand what's wrong
- Apply the minimal fix needed
- Follow style guide conventions

### 6. Verify Fix

Re-run the build/lint command to verify:
- All previous errors are resolved
- No new errors introduced

If new errors appear, fix those too (loop back to step 5).

### 7. Summary

Report what was fixed:
```
## Fix Summary

Fixed N issues:

- path/to/file.go:45 - Fixed nil pointer check
- path/to/file.go:72 - Changed Printf to slog

Build: PASS
Lint: PASS
```

### 8. Style Guide Recommendations (if applicable)

If fixes reveal a recurring pattern worth documenting:
- Describe the pattern observed
- Suggest the rule to add
- Ask user if they want to update the style guide via `/style-update`

**Do NOT automatically update style guides** - always discuss first.

## Key Guidelines

- **Minimal fixes**: Only fix what's broken, don't refactor
- **Verify changes**: Always re-run build/lint after fixes
- **User approval for style updates**: Recommend but don't auto-update
