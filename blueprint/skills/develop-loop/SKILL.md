---
description: Full development loop - implement story, fix build, review until clean
user-invocable: true
disable-model-invocation: true
argument-hint: "[story-file or story-name]"
allowed-tools: ["Read", "Glob", "Bash", "Task"]
---

# Development Loop

Full orchestrated development workflow: implement story tasks, fix build, review+fix cycle.

## Variables

- **QUERY_STORIES**: `${CLAUDE_PLUGIN_ROOT}/scripts/query-stories.sh`

## User Input
```
$ARGUMENTS
```

Parse for story file or name.

## Process

### 1. Resolve Story

If no story specified, list available stories:
```bash
bash QUERY_STORIES --actionable
```

Determine the story file path to pass to sub-agents.

### 2. Implement Story (Sub-agent)

Spawn sub-agent to implement story tasks:
```
Task: Implement story tasks
Use the Skill tool to invoke "blueprint:develop-story" with arguments: <story-file>
Return confirmation that story document was updated with completed tasks.
```

This handles:
- Task selection (interactive with user)
- Planning (interactive with user)
- Implementation
- **Updating story document with completed tasks** (mandatory)

Wait for completion.

### 3. Fix Build (Sub-agent)

Spawn sub-agent to fix any build/lint errors:
```
Task: Fix build errors
Use the Skill tool to invoke "blueprint:develop-fix" with arguments: --story <story-file>
Return build status (pass/fail).
```

Wait for completion.

### 4. Review Loop (Sub-agent)

Spawn sub-agent to run the review loop (simplify + fix + review-fix cycle):
```
Task: Run review loop
Use the Skill tool to invoke "blueprint:review-loop" with arguments: --story <story-file>
This only simplifies and fixes review issues - it does NOT implement new functionality.
Return final status and any remaining issues.
```

Wait for completion.

### 5. Present Results

Read the story document to get final task status.

Report:
```
## Development Loop Complete

### Story
[story file path]

### Tasks Implemented
[List tasks marked [x] in story]

### Build Status
[Pass / Fail]

### Review Status
[Pass / N remaining issues]

### Story Status
[Current status field value]
```

If issues remain:
```
### Remaining Issues
[List from review-fix output]

### Next Steps
- Run `/review-loop --story <file>` to continue
- Or fix manually and run `/review-code` to verify
```

## Notes

- Step 2's sub-agent interacts with user for task selection and plan approval
- Steps 3-5 run automatically without user interaction
- Story document is updated by /develop-story before continuing
- /review-loop handles simplify, build fix, and review-fix cycle
