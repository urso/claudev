---
description: Simplify code for clarity and maintainability while preserving functionality
user-invocable: true
disable-model-invocation: true
argument-hint: "[files] [--story story-file]"
allowed-tools: ["Read", "Write", "Edit", "Glob", "Grep", "Bash(git:*)"]
model: opus
---

# Simplify Code

Simplify and refine code for clarity, consistency, and maintainability while preserving exact functionality. Focuses on uncommitted changes unless specific files are provided.

## Variables

- **LIST_RULES**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-rules.sh`

## User Input
```
$ARGUMENTS
```

Parse for:
- Specific files to simplify (optional)
- `--story <story-file>` for context (optional)

## Process

### 1. Determine Files to Simplify

**If files specified:** Use those files.

**If no files specified:** Get uncommitted changes:
```bash
git diff --name-only HEAD
git diff --name-only --cached
```

**Filter out non-code files:**
- Skip: `*.md`, `*.json`, `*.yaml`, `*.yml`, `*.txt`, `*.lock`
- Skip: Files in `.git/`, `node_modules/`, `vendor/`, `dist/`, `build/`
- Include: Source code files (`.go`, `.py`, `.js`, `.ts`, `.tsx`, `.rs`, etc.)

If no files to simplify, report and exit.

### 2. Load Style Rules

```bash
bash LIST_RULES "" "" style
```

Output format (pipe-separated): `filename|name|applies-to|tags|paths|description`

Read applicable style rules to understand project conventions.

### 3. Load Story Context (if provided)

If `--story` specified, read the story document to understand:
- What task is being implemented
- The intent behind the code

### 4. Analyze and Simplify Each File

For each file, read and look for opportunities to:

**Reduce complexity:**
- Flatten unnecessary nesting
- Simplify conditional logic
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

For each simplification:
1. Verify the change preserves exact functionality
2. Apply the edit
3. Move to next opportunity

### 6. Summary

Report what was simplified:
```
## Simplification Summary

Simplified N file(s):

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

### Avoid Over-Simplification
Do NOT:
- Reduce code clarity or maintainability
- Create overly clever solutions that are hard to understand
- Combine too many concerns into single functions or components
- Remove helpful abstractions that improve code organization
- Prioritize "fewer lines" over readability (e.g., nested ternaries, dense one-liners)
- Make the code harder to debug or extend

### Minimal Scope
- Only simplify code in the target files
- Do not refactor surrounding code
- Do not add new features or functionality
- Do not fix bugs (use `/develop-fix` for that)

## What to Simplify

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
