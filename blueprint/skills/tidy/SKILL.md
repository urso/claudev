---
description: Tidy code for clarity and maintainability while preserving functionality
user-invocable: true
disable-model-invocation: true
argument-hint: "[files] [--story story-file] [--review-only]"
allowed-tools: ["Read", "Write", "Edit", "Glob", "Grep", "Bash(git:*)"]
model: opus
---

# Tidy Code

Tidy and refine code for clarity, consistency, and maintainability while preserving exact functionality. Focuses on uncommitted changes unless specific files are provided.

## Variables

- **LIST_RULES**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-rules.sh`

## User Input
```
$ARGUMENTS
```

Parse for:
- Specific files to tidy (optional)
- `--story <story-file>` for context (optional)
- `--review-only` flag to report issues without modifying files (optional)

## Pre-computed Context

### Changed Files (unstaged)
!`git diff --name-only HEAD`

### Changed Files (staged)
!`git diff --name-only --cached`

### Style Rules
!`bash ${CLAUDE_PLUGIN_ROOT}/scripts/list-rules.sh "" "" style`

## Process

### 1. Determine Files to Tidy

**If files specified:** Use those files.

**If no files specified:** Use the changed files from the pre-computed context above.

**Filter out non-code files:**
- Skip: `*.md`, `*.json`, `*.yaml`, `*.yml`, `*.txt`, `*.lock`
- Skip: Files in `.git/`, `node_modules/`, `vendor/`, `dist/`, `build/`
- Include: Source code files (`.go`, `.py`, `.js`, `.ts`, `.tsx`, `.rs`, etc.)

If no files to tidy, report and exit.

### 2. Load Style Rules

Read applicable style rules from the pre-computed context above to understand project conventions.

Output format (pipe-separated): `filename|name|applies-to|tags|paths|description`

### 3. Load Story Context (if provided)

If `--story` specified, read the story document to understand:
- What task is being implemented
- The intent behind the code

### 4. Analyze and Tidy Each File

For each file, read and look for opportunities to:

**Reduce complexity:**
- Flatten unnecessary nesting
- Reduce conditional logic complexity
- Remove dead code or unused variables
- Eliminate redundant abstractions that add no value
- Consolidate related logic that belongs together

**Improve clarity:**
- Replace nested ternaries with if/else or switch statements
- Use clearer variable/function names where obviously better
- Remove unnecessary comments that describe obvious code
- Prefer explicit over clever/compact code
- Choose clarity over brevity - explicit code is often better than overly compact code

**Apply project standards:**
- Follow style guide conventions
- Use consistent patterns with the rest of the codebase

### 5. Apply Changes

**Skip this step if `--review-only` is set.**

For each change:
1. Verify the change preserves exact functionality
2. Apply the edit
3. Move to next opportunity

### 6. Summary

**If `--review-only`:** Report findings without applying changes:
```
## Tidy Review

Found N issue(s) in M file(s):

### path/to/file.go
- Lines 45-52: Nested if statements could be flattened
- Line 78: Ternary chain could be replaced with switch

### path/to/other.ts
- Lines 12-20: Duplicate validation logic could be consolidated
- Line 35: Unused variable

No issues found:
- path/to/clean.go
```

**Otherwise:** Report what was tidied:
```
## Tidy Summary

Tidied N file(s):

### path/to/file.go
- Lines 45-52: Flattened nested if statements
- Line 78: Replaced ternary chain with switch

### path/to/other.ts
- Lines 12-20: Consolidated duplicate validation logic
- Line 35: Removed unused variable

No changes needed:
- path/to/clean.go
```

## Core Principles

### Preserve Functionality
Never change what the code does - only how it's written. All original features, outputs, and behaviors must remain intact.

### Clarity Over Brevity
Choose readable, explicit code over compact solutions. Three clear lines are better than one clever line.

### Avoid Over-Tidying
Do NOT:
- Reduce code clarity or maintainability
- Create overly clever solutions that are hard to understand
- Combine too many concerns into single functions or components
- Remove helpful abstractions that improve code organization
- Prioritize "fewer lines" over readability (e.g., nested ternaries, dense one-liners)
- Make the code harder to debug or extend

### Minimal Scope
- Only tidy code in the target files
- Do not refactor surrounding code
- Do not add new features or functionality
- Do not fix bugs (use `/develop-fix` for that)

## What to Tidy

**Good candidates:**
- Deeply nested conditionals (3+ levels)
- Repeated code blocks within a file
- Complex ternary expressions
- Obvious dead code
- Overly verbose patterns that have simpler equivalents

**Leave alone:**
- Code that's already clear and simple
- Abstractions that serve a purpose
- Patterns that match the project's established style
- Code outside the target files
