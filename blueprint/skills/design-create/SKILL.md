---
description: Create a design document for requirements, approach, and technical decisions
user-invocable: true
disable-model-invocation: true
argument-hint: "[feature/problem description]"
allowed-tools: ["Read", "Write", "Edit", "Glob", "Grep", "Bash", "AskUserQuestion", "Agent"]
hooks:
  PostToolUse:
    - matcher: "Write|Edit"
      hooks:
        - type: command
          command: "${CLAUDE_PLUGIN_ROOT}/scripts/hook-validate-doc.sh"
---

# Create Design Document

Create a design document that captures problem context, goals, approach, and technical decisions.

## Variables

- **DISCOVERY_GUIDE**: `${CLAUDE_PLUGIN_ROOT}/resources/discovery.md`
- **DESIGN_OPS**: `${CLAUDE_PLUGIN_ROOT}/resources/design-operations.md`
- **LIST_WORKFLOWS**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-workflows.sh`

## User Input
```
$ARGUMENTS
```

## Process

### 1. Load Guides

Read DISCOVERY_GUIDE and DESIGN_OPS for available tools and procedures.

Load workflow guidelines: `bash LIST_WORKFLOWS "" design`

Read the listed workflow files for design guidelines.

### 2. Understand the Problem

- Parse the user's initial description
- Search for relevant existing documents or code
- Identify scope and complexity

### 3. Gather Context

Ask clarifying questions to understand:
- **Problem**: What problem are we solving? Why now?
- **Goals**: What does success look like?
- **Non-goals**: What are we explicitly not doing?
- **Approach**: What's the high-level strategy?

Get user confirmation before proceeding.

### 4. Determine Design Name

Ask user for a short kebab-case name (e.g., `user-auth`, `api-redesign`).

### 5. Check Existing Designs

If related designs exist, ask user if this design `depends-on` any of them.

### 6. Gather References

Ask user if there are project files that should be referenced by this design (e.g., source code, configs, API specs). These are project-relative paths stored in the `references` frontmatter field and will be auto-loaded as context during story expansion and development.

### 7. Create Design Document

Follow the procedures in DESIGN_OPS to create the design. Fill in all sections with the gathered context.

The PostToolUse hook will automatically validate the document after writing.

### 8. Review

Spawn the `design-review` agent to review the newly created design with fresh context. Pass the design file path as the argument.

### 9. Confirm

- Report creation and review results
- Suggest next step: `/design-expand` to break into stories

## Guidelines

- **Be lean**: Problem/goals/approach/decisions only. Task details live in stories.
- **No task lists**: Designs describe what and why, not step-by-step how
- **Capture uncertainty**: Use "Open Questions" for unresolved issues
