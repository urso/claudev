---
description: Generate PR description from branch changes vs main
user-invocable: true
disable-model-invocation: false
argument-hint: "[--detailed] [--story story-file]"
allowed-tools: ["Read", "Bash(git:*)"]
model: haiku
effort: low
---

# Generate PR Description

Generate a pull request title and description based on changes between the current branch and main. Outputs the description for the user to copy or use with `gh pr create`.

## User Input
```
$ARGUMENTS
```

Parse for:
- `--detailed` for a more thorough description (optional)
- `--story <story-file>` for additional context (optional)

## Pre-computed Context

### Current Branch
!`git rev-parse --abbrev-ref HEAD`

### Merge Base
!`git merge-base main HEAD`

### Commit Log (branch vs main)
!`git log --oneline main..HEAD`

### Diff Summary (branch vs main)
!`git diff --stat main...HEAD`

### Full Diff (branch vs main)
!`git diff main...HEAD`

## Process

### 1. Check for Changes

If the diff summary is empty, report and exit:
```
No changes found between current branch and main.
```

If on main, report and exit:
```
Already on main. Create or switch to a feature branch first.
```

### 2. Load Context (if provided)

If `--story` provided, read the file to understand the motivation behind the work.

### 3. Analyze Changes

From the commits and diff, identify:
- The overall purpose of the branch (feature, bug fix, refactor, etc.)
- Key files and areas changed
- Breaking changes or notable decisions

### 4. Generate Description

**Title:**
- Max 70 characters
- Concise summary of the branch's purpose
- Imperative mood ("Add feature" not "Added feature")

**Body (default / concise mode):**
- A `## Summary` section with 1-3 bullet points covering what changed and why
- A `## Test plan` section with a brief checklist of how to verify

**Body (--detailed mode):**
- A `## Summary` section with a paragraph explaining the motivation and approach
- A `## Changes` section with grouped bullet points by area/file
- A `## Test plan` section with a thorough checklist
- If relevant: `## Breaking changes` or `## Migration notes`

### 5. Output

```
## Suggested PR

**Title:** <title>

**Body:**
<body in markdown>
```

## Guidelines

- **Concise by default**: The summary should help reviewers understand the PR at a glance
- **Accurate**: Describe actual changes, not intended changes
- **Reviewer-focused**: Highlight what a reviewer should pay attention to
- **No auto-create**: User decides when to create the PR
