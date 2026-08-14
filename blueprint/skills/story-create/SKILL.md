---
description: Create a story document, standalone or linked to a design
user-invocable: true
disable-model-invocation: true
argument-hint: "[description] or [--design design-id]"
allowed-tools: ["Read", "Write", "Edit", "Glob", "Grep", "Bash", "AskUserQuestion"]
hooks:
  PostToolUse:
    - matcher: "Write|Edit"
      hooks:
        - type: command
          command: "hook-validate-doc.sh"
---

# Create Story

Create a story document for a unit of implementable work. Can be standalone or linked to a parent design.

## Variables

- **DISCOVERY_GUIDE**: `${CLAUDE_PLUGIN_ROOT}/resources/discovery.md`
- **STORY_OPS**: `${CLAUDE_PLUGIN_ROOT}/resources/story-operations.md`
- **QUERY_STORIES**: `query-stories.sh`
- **LIST_WORKFLOWS**: `list-workflows.sh`

## User Input
```
$ARGUMENTS
```

Parse for:
- `--design <id>` — link to a parent design
- Or: description of the work for a standalone story

## Pre-computed Context

### Story Workflows
!`list-workflows.sh "" story`

## Process

### 1. Load Guides

Read DISCOVERY_GUIDE and STORY_OPS for available tools and procedures.

Read the workflow files listed in the pre-computed context above for story guidelines.

### 2. Determine Design Link

**If `--design` provided:**
- Locate and read the design for context
- Show existing stories for this design:
  ```bash
  bash QUERY_STORIES --design <id>
  ```
  This helps the user see what's already covered and what gaps remain.

**If no design specified:**
- Ask user if this story should be linked to a design
- If yes, list designs and let user choose
- If no, proceed as standalone

### 3. Gather Story Details

Ask clarifying questions:
- **What**: What work does this story cover?
- **Acceptance Criteria**: What must be true when done? Name types, signatures, behaviors. (e.g., "`AuthMiddleware` type exists with `Wrap(http.Handler) http.Handler` method")
- **Subtasks**: What specific work steps? Name files, functions, operations. (e.g., "Define `AuthMiddleware` struct in `pkg/auth/middleware.go`")
- **Dependencies**: Does this story depend on other stories?

Get user confirmation before proceeding.

### 4. Capture Technical Context

Scan the conversation for code blocks, API signatures, and technical decisions discussed.

**If technical content found:**
- Show the user what was found
- Ask which items should go in Technical Notes
- Include selected content in the story document

**If nothing found:** Skip to next step.

### 5. Determine Story Name

Ask user for a short kebab-case name (e.g., `auth-middleware`, `test-suite`).

### 6. Create Story

Follow the procedures in STORY_OPS to create the story:
- Generate ID, write the document, set dependencies
- Link to design if applicable

### 7. Confirm

- Report successful creation
- Suggest next steps: `/develop-story` to start implementation

## Guidelines

- Each story must have at least one `## Tasks` section with `- [ ]` checkboxes
- Keep stories small enough to implement in one session where possible
- Use `status: ready` if the story can be worked on now, `draft` if it needs refinement
- When linked to a design, reference the design's approach in Technical Notes
