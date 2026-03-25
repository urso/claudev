---
description: This skill should be used when the user asks to "add to style guide", "update style guide", "create style guide", "add a convention", "learned from review", "add coding rule", or needs guidance on maintaining project style guides for consistent development. Also use after code reviews reveal patterns worth capturing.
user-invocable: true
disable-model-invocation: false
argument-hint: "[rule-name] [rule]"
allowed-tools: ["Read", "Write", "Edit", "Glob", "Bash", "AskUserQuestion"]
---

# Style Guide

Create and update style guides in `docs/ai/rules/` to capture conventions and learnings.

## Variables

- **LIST_RULES**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-rules.sh`

## User Input
```
$ARGUMENTS
```

Parse for:
- Rule file name (e.g., `common`, `go`, `kubernetes`)
- Rule or learning to add

## Process

### 1. Check Existing Rules

```bash
bash LIST_RULES "" "" style
```

This lists all rules with the 'style' tag.

### 2. Create or Update

**New rule file** — create `docs/ai/rules/<name>.md`:

```markdown
---
name: <Name> Style Guide
applies-to:
  - <tech>  # or ["*"] for universal
tags:
  - style
paths: []  # optional: glob patterns for path-based filtering
description: <Brief description of what this guide covers>
---

# <Name> Style Guide

- [First rule]
```

**Existing rule file** — append the new rule to appropriate section, or create a section if needed.

For details on frontmatter fields, rule types, and integration with development, consult `references/rule-structure.md`.

### 3. Keep Rules Concise

Each rule should be:
- One clear sentence or short bullet
- Include a brief example only if the rule is ambiguous
- No verbose explanations — assume reader understands context

Good: `- Use structured logging with slog, not fmt.Printf`
Bad: `- When logging in Go, you should always use the structured logging package slog instead of using fmt.Printf because it provides better...`

Add examples only when ambiguous:
```
- Wrap errors with context: `fmt.Errorf("failed to parse config: %w", err)`
```

### 4. Confirm

Show what was added and confirm.

## Capturing Learnings

After code reviews reveal repeated feedback:

1. Identify the pattern or anti-pattern
2. Run `/style-update` to add the rule
3. Future development will follow the convention

This creates a feedback loop: reviews improve rules, rules improve code.

## Reference Files

- **`references/rule-structure.md`** — Rule file frontmatter fields, rule types, integration with development, symlink setup

## Usage Examples

```
/style-update common always use structured logging
/style-update go prefer table-driven tests
/style-update kubernetes    # interactive creation
```
