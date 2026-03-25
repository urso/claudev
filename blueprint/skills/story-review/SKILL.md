---
description: Review a story document against guidelines
user-invocable: true
disable-model-invocation: true
argument-hint: "[story-id or file]"
context: fork
allowed-tools: ["Read", "Edit", "Glob", "Bash"]
hooks:
  PostToolUse:
    - matcher: "Edit"
      hooks:
        - type: command
          command: "${CLAUDE_PLUGIN_ROOT}/scripts/hook-validate-doc.sh"
---

# Review Story

Validate a story against guidelines and check for completeness.

## Variables

- **DISCOVERY_GUIDE**: `${CLAUDE_PLUGIN_ROOT}/resources/discovery.md`
- **LIST_WORKFLOWS**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-workflows.sh`

## User Input
```
$ARGUMENTS
```

Parse for story file path, name, or ID.

## Process

### 1. Find Story

Read DISCOVERY_GUIDE for available tools and strategy, then locate the story based on user input.

### 2. Validate Structure

```bash
bash VALIDATE_DOC story <story-file>
```

If validation errors, report them immediately.

### 3. Load Story Guidelines

```bash
bash LIST_WORKFLOWS "" story
```

Read the listed workflow files for story guidelines.

### 4. Review Story

Review the story against the loaded guidelines. Check for:

1. **Guideline compliance**:
   - Has Tasks or Phase sections with checkboxes
   - Has Developer Logs sections (Decision Log, Blockers, Deviations, Lessons Learned)
   - Has proper frontmatter fields (id, title, status, created, description)
   - Tasks are concrete and actionable (not vague)

2. **Task quality**:
   - Tasks are broken into reasonable sub-tasks
   - Sub-tasks are specific enough to implement
   - No duplicate or overlapping tasks
   - Dependencies between tasks are logical

3. **Completeness**:
   - If linked to a design, does the story cover a coherent slice of the design?
   - Are Technical Notes present with implementation guidance?
   - For done stories: are developer logs filled in?

For each issue found, explain:
- Which section has the issue
- What the problem is
- Suggested fix

### 5. Present Results

Display findings to the user.

If issues found:
- List each issue with suggested fixes
- Ask if user wants to fix them now
- If yes, apply the fixes directly

If passed:
- Suggest next steps: `/develop-story` to implement

## Guidelines

- Review focuses on structure and task quality, not implementation correctness
- Fixes are applied directly when user approves
