---
description: Review code for style compliance, bugs, and efficiency. Use for staged/unstaged changes or specific files.
user-invocable: true
disable-model-invocation: false
argument-hint: "[files] [--story story-file] [--style-only] [--bugs-only] [--efficiency-only]"
allowed-tools: ["Read", "Grep", "Glob", "Bash(git:*)", "Task"]
---

# Code Review

Comprehensive code review checking style guide compliance, bugs/logic errors, and performance inefficiencies.

## Variables

- **DIFF_FILES**: `${CLAUDE_PLUGIN_ROOT}/scripts/diff-files.sh`
- **REVIEW_PIPELINE**: `${CLAUDE_PLUGIN_ROOT}/resources/review-pipeline.md`

## User Input
```
$ARGUMENTS
```

Parse for:
- Specific files to review (optional, defaults to staged+unstaged changes)
- `--story <story-file>` for context (optional)
- `--style-only` to run only style review
- `--bugs-only` to run only bug review
- `--efficiency-only` to run only efficiency review
- Branch comparison: "against <branch>", "vs <branch>" → pass branch to script

## Process

### 1. Get Changed Files

**If specific files provided:** Use those files directly (skip script, no chunking).

**If branch comparison requested:** Run `${DIFF_FILES} <branch>` to get chunked file list.

**Otherwise (default — staged+unstaged):** Run `${DIFF_FILES}` with no arguments.

### 2. Show Context

```
## Reviewing local changes

Files: X (Y lines total)
Chunks: Z
```

### 3. Run Review Pipeline

Read $REVIEW_PIPELINE and follow its process with the chunk output.

Pass through any flags (`--story`, `--style-only`, etc.).

