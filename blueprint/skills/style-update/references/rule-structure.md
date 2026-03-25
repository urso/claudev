# Rule File Structure

Rules in `docs/ai/rules/` capture conventions, learnings from code reviews, build instructions, and project-specific rules. Style rules (tagged with `style`) are loaded during development to ensure consistent implementation.

## Frontmatter Fields

Each rule file is a markdown file with YAML frontmatter:

```yaml
---
name: Go Style Guide
applies-to:
  - go
tags:
  - style
paths: []  # optional: glob patterns for path-based filtering
description: Go coding conventions and patterns.
---

# Go Style Guide

- Use structured logging with slog
- Prefer table-driven tests
```

- **name**: Display name for the rule file
- **applies-to**: Array of languages/technologies (use `["*"]` for universal)
- **tags**: Array of categories (`style`, `build`, `testing`, etc.)
- **paths**: Optional array of glob patterns for path-specific rules
- **description**: Brief description of what this rule file covers

## Rule Types

- **common.md** - Universal style rules (`applies-to: ["*"]`, `tags: [style]`)
- **Language rules** - go.md, python.md (`applies-to: [go]`, `tags: [style]`)
- **Technology rules** - kubernetes.md, helm.md (`applies-to: [kubernetes]`)
- **Build rules** - build.md (`tags: [build]`)
- **Path-specific rules** - frontend.md (`paths: ["frontend/**"]`)

## Integration with Development

During development, applicable rules are loaded based on:

1. **applies-to** - matches project languages/technologies
2. **tags** - filters by category (style, build, etc.)
3. **paths** - optionally matches file patterns being worked on

Rules from these files inform implementation decisions.

## Optional: Symlink to .claude/rules/

For auto-loading by Claude Code based on `paths:` patterns:

```bash
ln -s ../docs/ai/rules .claude/rules
```

This enables dual loading: skills use `applies-to` + `tags`, Claude Code uses `paths:`.
