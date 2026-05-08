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

**Body:**

Choose format based on complexity:

**Simple PRs** — single prose block covering why and what changed.

**Complex PRs** — structured format:
```
<Why paragraph — no header, just start writing the motivation/context>

## Changes
<Summary of what changed>

### Details
<Breakdown of changes — only if summary isn't enough>

### Attention
<Tricky, risky, or critical parts for reviewers — only if relevant>
```

**When to use structured format:**
- Multiple logical changes
- Non-obvious decisions or tradeoffs
- Touches risky/sensitive code
- Anything that warrants reviewer attention

**If unsure:** Ask the user which format they prefer.

### 5. Output

```
## Suggested PR

**Title:** <title>

**Body:**
<body in markdown>
```

## Guidelines

- **Reviewer-focused**: Help reviewers understand intent before reading code
- **Accurate**: Describe actual changes, not intended changes
- **No test plan**: User handles testing as part of their workflow
- **No auto-create**: User decides when to create the PR
