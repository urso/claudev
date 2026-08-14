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
          command: "${CLAUDE_PLUGIN_ROOT}/bin/hook-validate-doc.sh"
---

# Create Design Document

Create a design document that captures problem context, goals, approach, and technical decisions.

## Variables

- **DISCOVERY_GUIDE**: `${CLAUDE_PLUGIN_ROOT}/resources/discovery.md`
- **DESIGN_OPS**: `${CLAUDE_PLUGIN_ROOT}/resources/design-operations.md`
- **LIST_WORKFLOWS**: `list-workflows.sh`
- **LIST_TEMPLATES**: `list-templates.sh`

## User Input
```
$ARGUMENTS
```

## Pre-computed Context

### Design Workflows
!`list-workflows.sh "" design`

## Process

### 1. Load Guides

Read DISCOVERY_GUIDE and DESIGN_OPS for available tools and procedures.

Read the workflow files listed in the pre-computed context above for design guidelines.

### 2. Select Template

Run `bash LIST_TEMPLATES design` to list available templates. Output format: `name|path|description`.

If multiple templates exist, present them to user and ask which to use. If only `default`, use it without asking.

Read the selected template.

### 3. Follow Template Instructions

Follow the `agent-instructions` field in the template's frontmatter. Use `<!-- agent: ... -->` comments as per-section guidance.

The template drives the conversation — gathering problem context, goals, requirements, etc.

### 4. Determine Design Name

Ask user for a short kebab-case name (e.g., `user-auth`, `api-redesign`).

### 5. Check Existing Designs

If related designs exist, ask user if this design `depends-on` any of them.

### 6. Gather References

Ask user if there are project files that should be referenced by this design (e.g., source code, configs, API specs). These are project-relative paths stored in the `references` frontmatter field and will be auto-loaded as context during story expansion and development.

### 7. Create Design Document

Follow the procedures in DESIGN_OPS to create the design. Fill in all sections with the gathered context.

The PostToolUse hook will automatically validate the document after writing.

### 8. Clean Up Template Artifacts

Before finalizing the document:
1. Remove the `agent-instructions` field from frontmatter
2. Remove all `<!-- agent: ... -->` comments (single-line and multi-line)
3. Collapse multiple blank lines to one

### 9. Review

Spawn the `design-reviewer` agent to review the newly created design with fresh context. Pass the design file path as the argument.

### 10. Confirm

- Report creation and review results
- Suggest next step: `/design-expand` to break into stories

## Guidelines

- **Be lean**: Problem/goals/approach/decisions only. Task details live in stories.
- **No task lists**: Designs describe what and why, not step-by-step how
- **Capture uncertainty**: Use "Open Questions" for unresolved issues
