---
description: Initialize docs/ai structure for development workflow
user-invocable: true
disable-model-invocation: true
allowed-tools: ["Read", "Write", "Bash", "Glob", "Skill"]
---

# Initialize Development Workflow

Create the `docs/ai/` directory structure for designs, stories, and style guides.

## Process

### 1. Create Directory Structure

```bash
mkdir -p docs/ai/stories docs/ai/designs docs/ai/rules docs/ai/workflows
```

### 2. Create Starter Files (only if they don't exist)

**Only create if file doesn't exist:**

Check and create `docs/ai/rules/common.md`:
```bash
test -f docs/ai/rules/common.md && echo "EXISTS" || echo "NEW"
```

If NEW, create:
```markdown
---
name: Common Style Guide
applies-to: ["*"]
tags: [style]
description: Universal coding conventions for all languages.
---

# Common Style Guide

```

Check and create `docs/ai/workflows/design.md`:
```bash
bash ${CLAUDE_PLUGIN_ROOT}/scripts/get-guidelines.sh design docs/ai/workflows > /dev/null && echo "OK"
```

This creates the design workflow guidelines from template if missing.

Check and create `docs/ai/workflows/story.md`:
```bash
bash ${CLAUDE_PLUGIN_ROOT}/scripts/get-guidelines.sh story docs/ai/workflows > /dev/null && echo "OK"
```

This creates the story workflow guidelines from template if missing.

Check and create `docs/ai/rules/build.md`:
```bash
test -f docs/ai/rules/build.md && echo "EXISTS" || echo "NEW"
```

If NEW, create:
```markdown
---
name: Build Configuration
applies-to: ["*"]
tags: [build]
description: Build, lint, test, and format commands for this project.
---

# Build Configuration

## Build

## Lint

## Test

## Format

```

### 3. Update .gitignore (only if entries don't exist)

Check if gitignore entries already exist:
```bash
grep -q "docs/ai/designs/" .gitignore 2>/dev/null && echo "EXISTS" || echo "NEW"
```

If NEW, append to `.gitignore`:
```
# AI development artifacts
docs/ai/designs/
docs/ai/stories/
```

### 4. Report

```
## Initialized docs/ai/

Structure:
  docs/ai/
  ├── designs/    (gitignored)
  ├── stories/    (gitignored)
  ├── rules/
  │   ├── common.md  [created/existed]  - universal style conventions
  │   └── build.md   [created/existed]  - build/lint/test commands
  └── workflows/
      ├── design.md  [created/existed]  - design workflow guidelines
      └── story.md   [created/existed]  - story workflow guidelines

Next steps:
- Edit docs/ai/rules/build.md with your build/lint/test commands
- Edit docs/ai/rules/common.md with project conventions
- Add language-specific rules (e.g., docs/ai/rules/go.md)
- Run /design-create to start designing
```
