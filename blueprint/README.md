# Blueprint

Design-driven development workflow plugin for Claude Code.

Blueprint structures your work into **designs** (what and why) and **stories** (concrete tasks), with built-in code review, style guides, and automated development loops.

## Getting Started

```
/init
```

This creates the workspace structure:

```
docs/ai/
  designs/    # Design documents (gitignored)
  stories/    # Story documents (gitignored)
  rules/      # Style guides and build config (tracked)
  workflows/  # Workflow guidelines (tracked)
```

## Workflows

### Design-Driven Development

The full workflow takes an idea from design through implementation:


```
/design-create       Create a design document (problem, goals, approach)
       |
/design-review       Review design for structure and clarity
       |
/design-expand       Break the design into stories with dependencies
       |
/develop-story       Implement selected tasks from a story
       |
/story-update        Verify story doc is up to date (tasks, logs, status)
       |
/design-update       Incorporate learnings back into the design
       |
/design-archive      Archive completed designs and stories
```

Use `/develop-loop` instead of `/develop-story` when you want to chain implementation, build fix, and review in one go. Prefer `/develop-story` for ambiguous tasks where you want to review the approach before committing to a full cycle.

Stories can be looked up by ID number or by design name — e.g. `/develop-story 3` or `/develop-story auth`.

New stories can be added to a design at any time via `/story-create` or `/design-update` — for example when blockers or new requirements surface during development.

### Iterative Refinement

When implementation reveals that a story was too ambiguous or the approach is off:

1. Run `/develop-story` (not `/develop-loop`) to implement a subset of tasks
2. Review what the agent produced
3. Discuss with the agent to refine the story document
4. Discard code changes if needed, re-run `/develop-story`
5. Once a story is completed, run `/design-update` to keep the design in sync

### Story-First Development

For smaller work that doesn't need a full design:

```
/story-create        Create a standalone story with tasks
       |
/develop-story       Implement selected tasks
       |
/story-update        Verify story doc is up to date
```



### Ad-Hoc Development

For quick tasks without formal tracking, use `/development` for guidance on implementing directly with style guide compliance.

## Skills Reference

### Design

| Skill | Description |
|-------|-------------|
| `/design-create` | Create a design document through guided questions |
| `/design-review` | Review a design for structure, clarity, and compliance |
| `/design-expand` | Break a design into stories with dependency ordering |
| `/design-update` | Update a design with implementation learnings |
| `/design-archive` | Archive completed designs and their stories |

### Stories

| Skill | Description |
|-------|-------------|
| `/story-create` | Create a story (standalone or linked to a design) |
| `/story-review` | Validate story structure and task quality |
| `/story-update` | Verify and update a story — tasks marked, logs filled in, status correct |

### Implementation

| Skill | Description |
|-------|-------------|
| `/develop-story` | Implement selected tasks from a story |
| `/develop-fix` | Fix build and lint errors |
| `/develop-simplify` | Simplify code for clarity without changing behavior |
| `/develop-loop` | Full loop: implement -> fix -> review (orchestrates everything) |

### Code Review

| Skill | Description |
|-------|-------------|
| `/review-code` | Review code for style, bugs, and efficiency |
| `/review-fix` | Review and fix issues iteratively (max 5 cycles) |
| `/review-loop` | Simplify -> fix -> review-fix cycle |

### Utilities

| Skill | Description |
|-------|-------------|
| `/style-update` | Add conventions to style guides from review learnings |
| `/commit-message` | Generate a commit message from staged changes |
| `/prime` | Load style guides and build config into context |
| `/handover` | Generate a handover prompt for continuing in a new session |

## How It Works

### Documents

**Designs** capture the problem, goals, approach, and architecture decisions. They answer *what* to build and *why*.

**Stories** capture concrete tasks with checkboxes, technical notes, and developer logs. They answer *how* and track implementation progress.

Both use YAML frontmatter for metadata (id, title, status, dependencies) and are managed through shell scripts that handle ID generation, status transitions, and dependency tracking.

### Status Lifecycle

Designs: `draft` -> `ready` -> `in-progress` -> `done` (also: `cancelled`, `archived`)

Stories: `ready` -> `in-progress` -> `done` (also: `cancelled`, `archived`)

Stories can declare `blocked-by` dependencies. Only unblocked stories show up as actionable work.

### Quality Gates

1. **Design review** validates clarity and completeness before story creation
2. **Story review** validates task quality before implementation
3. **Code review** checks style compliance, bugs, and efficiency
4. **Build verification** ensures code compiles and passes lint

Reviews are triggered automatically during creation workflows and can also be run independently.

### Feedback Loops

- Implementation deviations feed back into designs via `/design-update`
- Recurring code review patterns become style rules via `/style-update`
- Developer logs in stories capture decisions and lessons for future work
- `/develop-story` reads developer logs from direct dependency stories, so lessons and decisions carry forward
- `/handover` includes developer log entries so new sessions don't lose context

### Implementation In Practice

A typical session works on a subset of tasks from a story:

1. `/develop-story` — select tasks, review the plan, implement
2. `/story-update` — verify the story doc reflects what was done
3. `/develop-fix` — fix any build/lint errors
4. Optionally `/review-loop` — simplify and review code

`/develop-loop` chains all of these automatically (implement -> fix -> review). It's convenient when the task is well-defined, but `/develop-story` gives more control for ambiguous or exploratory work.

## Style Guides and Rules

Rules live in `docs/ai/rules/` as markdown files with YAML frontmatter. Think of them as CLAUDE.md split into separate, tagged files — each rule file covers one topic (a language, build commands, a specific area) and frontmatter controls when it gets loaded.

```yaml
---
name: Go Style Guide
applies-to:
  - go
tags:
  - style
paths: []
description: Go coding conventions and patterns.
---

# Go Style Guide

- Use structured logging with slog
- Prefer table-driven tests
```

### Frontmatter Fields

- **name**: Display name
- **applies-to**: Languages/technologies this rule covers (use `["*"]` for universal)
- **tags**: Categories — `style`, `build`, `testing`, etc.
- **paths**: Optional glob patterns for path-specific rules (e.g. `["frontend/**"]`)
- **description**: Brief description

### Rule Types

| File | applies-to | tags | Purpose |
|------|-----------|------|---------|
| `common.md` | `["*"]` | `[style]` | Universal style rules |
| `go.md` | `[go]` | `[style]` | Language-specific conventions |
| `build.md` | `["*"]` | `[build]` | Build, lint, test commands |
| `frontend.md` | `[typescript]` | `[style]` | Path-specific rules via `paths: ["frontend/**"]` |

Skills like `/develop-story` and `/review-code` load applicable rules automatically. You can also prime the context manually with `/prime` at the start of a session — this is useful for ad-hoc work or when you want rules loaded before starting. Use `/style-update` to add new conventions from code review learnings.

### Optional: Auto-loading via Claude Code

Symlink the rules directory to enable Claude Code's built-in path-based rule loading:

```bash
ln -s ../docs/ai/rules .claude/rules
```

This enables dual loading: skills use `applies-to` + `tags`, Claude Code uses `paths:`.

## `/prime` - Loading Context

Use `/prime` to load style guides and build configuration before working:

- `/prime` - Load everything
- `/prime go` - Load only Go rules
- `/prime build` - Load only build configuration
- `/prime style` - Load only style guides
