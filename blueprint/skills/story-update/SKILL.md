---
description: Verify and update a story document — check tasks are marked, logs are filled in, status is correct
user-invocable: true
disable-model-invocation: true
argument-hint: "[story-id or file]"
allowed-tools: ["Read", "Edit", "Glob", "Grep", "Bash", "AskUserQuestion"]
hooks:
  PostToolUse:
    - matcher: "Edit"
      hooks:
        - type: command
          command: "${CLAUDE_PLUGIN_ROOT}/scripts/hook-validate-doc.sh"
---

# Update Story

Verify that a story document is up to date. Check that completed tasks are marked, developer logs capture decisions and blockers, and status is correct.

Typically called within the same session after working on tasks, so the agent already has context about what was done.

## Variables

- **DISCOVERY_GUIDE**: `${CLAUDE_PLUGIN_ROOT}/resources/discovery.md`

## User Input
```
$ARGUMENTS
```

Parse for story file path, name, or ID.

## Process

### 1. Find Story

Read DISCOVERY_GUIDE for available tools and strategy, then locate the story based on user input. Read the story document.

### 2. Check Task Completion

Review each task and sub-task in the story. Based on work done in this session (or the current state of uncommitted changes if context is unclear), flag:

- Tasks that were completed but not yet marked `[x]`
- Tasks marked `[x]` that don't appear to be done

Present findings and ask the user to confirm before updating.

### 3. Check Developer Logs

Review the Developer Logs sections. Based on work done in this session, propose entries for:

- **Decision Log**: Technical decisions made and rationale
- **Blockers Encountered**: Issues hit and how they were resolved
- **Deviations from Design**: Divergence from the original approach (only if story has a `design:` field)
- **Lessons Learned**: Insights for future stories

Present proposed entries to the user. Let them accept, edit, or skip each one before writing.

### 4. Update Status

Based on task state:
- All tasks complete → suggest `done`
- Some tasks worked on → ensure `in-progress`

Ask user to confirm any status change.

### 5. Check Design Impact

If the story has a `design:` field and deviations or blockers were logged:

"Deviations or blockers were logged. Consider running `/design-update` to keep the design in sync."

## Guidelines

- Never fabricate log entries — only record what the user confirms
- Keep log entries concise
- This skill updates the story document only, it does not modify code
