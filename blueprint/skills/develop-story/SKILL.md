---
description: Develop tasks from a story document with planning and implementation
user-invocable: true
disable-model-invocation: true
argument-hint: "[story or design name/id]"
allowed-tools: ["Read", "Write", "Edit", "Glob", "Grep", "Bash", "AskUserQuestion"]
---

# Develop Story Tasks

Develop tasks from a story document with comprehensive planning and implementation support.

## Variables

- **RESOLVE_DIR**: `${CLAUDE_PLUGIN_ROOT}/scripts/resolve-dir.sh`
- **QUERY_WORK**: `${CLAUDE_PLUGIN_ROOT}/scripts/query-work.sh`
- **FIND_DOC**: `${CLAUDE_PLUGIN_ROOT}/scripts/find-doc.sh`
- **LIST_RULES**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-rules.sh`
- **LIST_WORKFLOWS**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-workflows.sh`

## User Input
```
$ARGUMENTS
```

## Process

### 1. Resolve Directories

```bash
STORIES_DIR=$(bash RESOLVE_DIR story)
```

### 2. Find a Story to Work On

Find actionable work using `query-work.sh`. Always use `--actionable` to only see stories that are ready to implement.

```bash
bash QUERY_WORK --actionable --search "<user input>"
```

If no input provided, omit `--search` to list all actionable work.

Output format: `TYPE | ID | FILENAME | TITLE | STATUS | BLOCKED-BY | DESCRIPTION`

When a design matches the search, its actionable stories are included automatically. Present the story list and let the user pick one. If no actionable stories are found, tell the user — do not try to resolve blockers or work on blocked stories.

Read the selected story document.

### 3. Load Design Context

If the story has a `design:` frontmatter field:
1. Use FIND_DOC to locate and read the parent design
2. Read any designs listed in the parent design's `depends-on` field
3. Read any files listed in the parent design's `references` field (paths relative to git root)

If the story has a `blocked-by:` field, read the Developer Logs sections (Decision Log, Blockers, Deviations, Lessons Learned) from those dependency stories. This surfaces decisions and lessons that are directly relevant to the current work.

This provides architectural context, decisions, and relevant project files for implementation.

### 4. Load Applicable Rules

List all rules:
```bash
bash LIST_RULES
```

Output format (grouped by directory, pipe-separated columns):
```
<directory>:
full/path/to/file.md|name|applies-to|tags|paths|description
```

Read and apply:
- All rules where `applies-to` includes `*` or matches the project's languages/technologies
- Style rules (tagged `style`) for coding conventions
- Build rules (tagged `build`) for build/lint/test commands

Keep these guidelines in context throughout implementation.

### 5. Task Selection Phase

- Display all tasks and sub-tasks from the story document with completion status
- Ask user to specify which tasks they want to work on
- **IMPORTANT**: Only work on explicitly selected tasks - do not touch other tasks
- Validate that selected tasks are not already completed

### 6. Detailed Planning Phase

For each selected task:
- Create a comprehensive development plan
- Break down into implementation steps
- Identify files to modify
- Consider dependencies and integration points
- Present the plan to the user for discussion and approval
- **Do not proceed to implementation until the plan is explicitly approved**

### 7. Implementation Phase (after plan approval)

Implement the approved plan directly:
- Create/modify the identified files
- Follow all applicable style guide rules loaded in step 2
- Write tests as specified
- Handle integration points carefully

### 8. Story Document Updates

**CRITICAL**: Update the story document before completing:

- Mark ALL completed tasks as `[x]` in the story document
- Update Developer Logs sections as needed:
  - **Decision Log**: Key technical decisions and rationale
  - **Blockers Encountered**: Issues faced and resolutions
  - **Deviations from Design**: Changes from original design and why
  - **Lessons Learned**: Insights for future stories
- Update status field to `in-progress`

**This step is mandatory** - the story document must reflect completed work before this command finishes.

### 9. Check for Design Impact

After updating the story document, check if significant deviations were logged:

- Read the "Deviations from Design" section
- If the story has a `design:` field and deviations are substantial (architectural changes, scope changes, discovered blockers):

  "Deviations were logged that may affect the design. Consider running `/design-update` to revise."

## Key Guidelines

- **Task Selection**: Only work on user-selected tasks, never modify others
- **Planning First**: Always plan before implementing, wait for explicit approval
- **Apply Rules**: Follow all applicable style and build rules from step 2
- **Update Story**: Mark tasks complete in the story document as you finish them
- **Stay Collaborative**: Keep user informed and ask for confirmation on major changes
- **Feedback Loop**: Prompt for design update when deviations impact the design
