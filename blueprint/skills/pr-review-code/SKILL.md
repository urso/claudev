---
description: Review all code and config changes on the current branch vs base. Detects PR and uses correct base branch.
user-invocable: true
disable-model-invocation: true
argument-hint: "[--story story-file] [--style-only] [--bugs-only] [--efficiency-only]"
allowed-tools: ["Read", "Grep", "Glob", "Bash(git:*)", "Bash(gh:*)", "Task"]
---

# PR Code Review

Review all code changes on the current branch compared to the base branch.

## Variables

- **PR_DIFF_FILES**: `pr-diff-files.sh`
- **REVIEW_PIPELINE**: `${CLAUDE_PLUGIN_ROOT}/resources/review-pipeline.md`

## User Input
```
$ARGUMENTS
```

Parse for:
- `--story <story-file>` for context (optional)
- `--style-only`, `--bugs-only`, `--efficiency-only` flags

## Process

### 1. Get PR Context and Changed Files

Run `${PR_DIFF_FILES}` to get PR info and chunked file list.

The output includes:
- PR info (number, title) or "No PR" message
- Base branch being compared against
- Files grouped into ~500 line chunks

### 2. Report Context

If PR exists:
```
## Reviewing PR #<number>: <title>

Base: <base branch>
Files: X (Y lines total)
Chunks: Z
```

If no PR:
```
## Reviewing branch changes

Comparing against: main
Files: X (Y lines total)
Chunks: Z
```

### 3. Run Review Pipeline

Read $REVIEW_PIPELINE and follow its process with the chunk output.

Pass through any flags (`--story`, `--style-only`, etc.).

## Notes

- Works whether or not a PR exists — compares against `main` if no PR
- Reviews local changes, including uncommitted/unpushed code
- To review a remote PR, first run `gh pr checkout <PR#>`
