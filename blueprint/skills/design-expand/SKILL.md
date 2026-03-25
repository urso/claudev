---
description: Break a design into multiple stories with dependency wiring
user-invocable: true
disable-model-invocation: true
argument-hint: "[design-id or file]"
allowed-tools: ["Read", "Write", "Edit", "Glob", "Grep", "Bash", "AskUserQuestion", "Agent"]
hooks:
  PostToolUse:
    - matcher: "Write|Edit"
      hooks:
        - type: command
          command: "${CLAUDE_PLUGIN_ROOT}/scripts/hook-validate-doc.sh"
---

# Expand Design into Stories

Break a design document into multiple stories, setting up `blocked-by` dependencies between them.

## Variables

- **DISCOVERY_GUIDE**: `${CLAUDE_PLUGIN_ROOT}/resources/discovery.md`
- **DESIGN_OPS**: `${CLAUDE_PLUGIN_ROOT}/resources/design-operations.md`
- **STORY_OPS**: `${CLAUDE_PLUGIN_ROOT}/resources/story-operations.md`
- **QUERY_STORIES**: `${CLAUDE_PLUGIN_ROOT}/scripts/query-stories.sh`
- **SET_STATUS**: `${CLAUDE_PLUGIN_ROOT}/scripts/set-status.sh`
- **LIST_WORKFLOWS**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-workflows.sh`
- **FIND_DOC**: `${CLAUDE_PLUGIN_ROOT}/scripts/find-doc.sh`

## User Input

```
$ARGUMENTS
```

Parse for design file path, name, or ID.

## Process

### 1. Load Guides

Read DISCOVERY_GUIDE, DESIGN_OPS, and STORY_OPS for available tools and procedures.

### 2. Find Design

Locate the design based on user input. Read the design document fully.

### 3. Load Design Context

Load additional context from the design's frontmatter:

- **`depends-on`**: For each dependency ID, use `find-doc.sh` to locate and read the dependent design. These provide architectural context and decisions that this design builds on.
- **`references`**: Read each referenced file (paths are relative to git root). These are project files that inform the design.

Keep this context available when proposing stories — it helps identify integration points and constraints.

### 4. Check Existing Stories

```bash
bash QUERY_STORIES --design <design-id>
```

Show existing stories for this design so we don't create duplicates.

### 5. Load Story Guidelines

```bash
bash LIST_WORKFLOWS "" story
```

Read the listed workflow files for story guidelines.

### 6. Propose Story Breakdown

Analyze the design's goals, requirements, and technical approach. Propose a set of stories that:
- Cover the full scope of the design
- Account for existing stories (don't duplicate)
- Are self-contained and independently implementable where possible
- Have clear dependency ordering

Present the proposed stories as a numbered list:
```
1. [Story Title] - [one-line description]
   blocked-by: (none)
2. [Story Title] - [one-line description]
   blocked-by: 1
3. [Story Title] - [one-line description]
   blocked-by: 1, 2
```

Ask user to review, modify, add, or remove stories before creating them.

### 7. Create Stories

Follow the procedures in STORY_OPS to create each approved story. Link each to this design via `design:` frontmatter. Set `blocked-by` relationships between them.

If the user requests changes to stories after creation (reorder, modify, cancel), follow the relevant procedures in STORY_OPS.

### 8. Update Design Status

If the design is in `draft` status, offer to move it to `in-progress`:
```bash
bash SET_STATUS design <design-id> in-progress
```

### 9. Review Stories

Spawn `story-review` agents in parallel for each created story file. Pass the story file path as the argument to each agent.

### 10. Summary

Report all created stories with their IDs, titles, dependency graph, and review results. Suggest:
- `/develop-story` to start implementing the first actionable story
- `/design-review` if the design hasn't been reviewed yet

## Guidelines

- Stories should be small enough to implement in one session where possible
- Each story must have at least one `## Tasks` section with checkboxes
- Dependencies flow forward — later stories depend on earlier ones, not vice versa
- Set status to `ready` for stories with no blockers, `draft` if they need refinement
- Don't over-decompose — 3-7 stories per design is typical
