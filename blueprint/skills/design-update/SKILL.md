---
description: Update a design with learnings from stories
user-invocable: true
disable-model-invocation: true
argument-hint: "[design-id or file]"
allowed-tools: ["Read", "Write", "Edit", "Glob", "Grep", "Bash", "Task", "AskUserQuestion"]
hooks:
  PostToolUse:
    - matcher: "Write|Edit"
      hooks:
        - type: command
          command: "${CLAUDE_PLUGIN_ROOT}/scripts/hook-validate-doc.sh"
---

# Update Design

Update a design and its stories based on implementation learnings. Can update the design document, add new stories, adjust story dependencies, or cancel stories that are no longer needed.

## Variables

- **DISCOVERY_GUIDE**: `${CLAUDE_PLUGIN_ROOT}/resources/discovery.md`
- **DESIGN_OPS**: `${CLAUDE_PLUGIN_ROOT}/resources/design-operations.md`
- **STORY_OPS**: `${CLAUDE_PLUGIN_ROOT}/resources/story-operations.md`
- **QUERY_STORIES**: `${CLAUDE_PLUGIN_ROOT}/scripts/query-stories.sh`
- **SET_STATUS**: `${CLAUDE_PLUGIN_ROOT}/scripts/set-status.sh`
- **LIST_WORKFLOWS**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-workflows.sh`

## User Input
```
$ARGUMENTS
```

Parse for design file path, name, or ID.

## Process

### 1. Find Design

Read DISCOVERY_GUIDE, DESIGN_OPS, and STORY_OPS for available tools and procedures. Locate the design based on user input. Read the design document.

### 2. Extract Insights from Stories

Use a sub-agent to read all stories and extract insights:

```bash
bash QUERY_STORIES --design <design-id>
```

Spawn a sub-agent with the story file paths and the design document:

```
Task: Extract design-relevant insights from stories
subagent_type: general-purpose

Read the design document at [design-file-path].
Read each of these story files: [story-file-paths]

Extract insights that may require updating the design:
- **Deviations**: Decisions in stories that changed the original design approach
- **New decisions**: Architecture decisions made during implementation not captured in the design
- **Blockers resolved**: Issues that required design-level changes
- **Open questions answered**: Questions from the design that stories resolved
- **Gaps**: Work discovered during implementation that isn't covered by any existing story
- **Obsolete stories**: Stories that are no longer needed due to approach changes
- **Dependency issues**: Stories whose blocked-by relationships need adjusting

Also report story statuses (how many done vs in-progress vs other).

Return a structured summary of findings. Only include items that are relevant — skip routine implementation details.
```

### 3. Load Design Guidelines

```bash
bash LIST_WORKFLOWS "" design
```

Read the listed workflow files for design guidelines.

### 4. Identify Updates

Present the sub-agent's findings to the user.

Ask user which updates to apply.

### 5. Apply Design Updates

Update the design document with approved changes:
- Update Technical Approach if the approach changed
- Add to Architecture Decisions with new decisions
- Remove answered items from Open Questions
- Keep the design lean — don't add implementation details or tasks

### 6. Apply Story Changes

Follow the procedures in STORY_OPS for all story modifications. Based on the sub-agent's findings and user approval:

- **New stories**: Create and link to this design
- **Cancel stories**: Cancel and clean up blocked-by references
- **Reorder dependencies**: Adjust blocked-by relationships

Only make changes the user explicitly approves.

### 7. Check Design Completion

Check story statuses from step 2:
- If ALL stories for this design have `status: done`, suggest marking the design as done:

  "All stories for this design are done. Mark design as done?"

- If user agrees:
  ```bash
  bash SET_STATUS design <id> done
  ```

### 8. Confirm

Report what was updated (design changes, new stories, cancelled stories, dependency changes) and suggest next steps if applicable.

## Guidelines

- Keep the design lean — update approach and decisions, don't add task details
- Preserve the original structure, only modify sections that need updating
- Design completion is suggested, not enforced — user decides
- Follow STORY_OPS for all story creation, cancellation, and dependency management
