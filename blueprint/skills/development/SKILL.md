---
description: This skill should be used when the user asks to "implement X", "add X feature", "create X", "refactor X", "fix X bug", "develop this story", "work on story", or needs guidance on development work - both ad-hoc tasks and story-based workflows with style guide compliance.
---

# Development Skill

This skill provides guidance for both ad-hoc development tasks and structured story-based development with style guide support.

## Overview

### Story-Based Workflow

1. Select a story document
2. Load applicable style guides and build config
3. Choose specific tasks to implement
4. Plan the implementation (get approval)
5. Execute the plan
6. Update the story with progress and learnings

### Ad-Hoc Workflow

1. Detect development task (implement/add/fix/refactor)
2. Load applicable style guides and build config
3. Plan the implementation (get approval)
4. Execute the plan
5. Verify with build/lint

## When to Use

**Story-Based:**
- Have a story document ready
- Need to implement specific planned tasks
- Want structured development with progress tracking

**Ad-Hoc:**
- Quick feature additions
- Bug fixes
- Refactoring work
- One-off implementations
- Any coding task without a story document

## Rules

Rules in `docs/ai/rules/` provide conventions and build instructions:

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
```

**Frontmatter fields:**
- `name` - display name
- `applies-to` - languages/technologies (use `["*"]` for universal)
- `tags` - categories: `style`, `build`, `testing`, etc.
- `paths` - optional glob patterns for path-based filtering
- `description` - brief description

Rules with `applies-to: ["*"]` apply to all code. Use `tags` to filter by category (style vs build).

## Commands

### Loading Development Context

To manually load style guides and build config:

```
/prime [tech...] [build|style]
```

Examples:
- `/prime` - Load everything
- `/prime go` - Load go-related style guides + build.md
- `/prime build` - Only load build.md
- `/prime style` - Only load style guides
- `/prime go rust` - Load go and rust style guides + build.md

**Note:** For ad-hoc development tasks, the skill automatically loads context. Use this command when you want explicit control or to reload context.

### Developing a Story

To implement tasks from a story document:

```
/develop-story [story-file or story-name]
```

The command will:
1. Find and load the story document
2. Load applicable style guides and build config
3. Show tasks and completion status
4. Ask which tasks to work on
5. Create detailed implementation plan
6. Wait for approval before implementing
7. Execute and update story

### Ad-Hoc Development

For ad-hoc tasks without a story document, simply describe what you want:

```
"Implement user authentication"
"Add a delete button to the profile page"
"Refactor the API client to use async/await"
"Fix the memory leak in the event handler"
```

The skill will:
1. Auto-load applicable style guides and build config
2. Create detailed implementation plan
3. Wait for approval before implementing
4. Execute the plan
5. Verify with build/lint commands

### Fixing Build/Lint Errors

To fix build or lint errors:

```
/develop-fix [--story story-file]
```

The command will:
1. Load build configuration
2. Run build/lint commands
3. Fix any errors found
4. Verify fixes pass
5. Optionally suggest style guide updates

### Full Development Loop

For automated end-to-end development:

```
/develop-loop [story-file]
```

This orchestrates:
1. `/develop-story` - Implement tasks (interactive)
2. `/develop-fix` - Fix build errors
3. `/review-loop` - Simplify, fix, and review-fix cycle
4. Present final results

## Key Principles

**Story-Based Development:**
- **Task selection**: Only work on selected tasks, never touch others
- **Planning first**: Always create and approve a plan before implementing
- **Style compliance**: Follow all loaded style guide rules
- **Update story**: Mark tasks complete immediately, update developer logs
- **Stay collaborative**: Confirm major changes with user

**Ad-Hoc Development:**
- **Auto-load context**: Load applicable style guides and build config automatically
- **Planning first**: Always create and approve a plan before implementing
- **Style compliance**: Follow all loaded style guide rules
- **Verification**: Run build/lint after implementation
- **Stay collaborative**: Confirm major changes with user

## Developer Logs

After completing tasks, update the story's developer logs:

- **Decision Log**: Key technical decisions and rationale
- **Blockers Encountered**: Issues and resolutions
- **Deviations from Design**: Changes and reasons
- **Lessons Learned**: Insights for future work

These logs feed into future story creation and design updates.
