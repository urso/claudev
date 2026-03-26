---
description: Create agent handover prompt for continuing work in a new session
user-invocable: true
disable-model-invocation: true
argument-hint: "[additional context]"
allowed-tools: ["Read", "Glob"]
---

# Create Agent Handover Prompt

Generate a handover prompt that can be copied to continue work in a new Claude session. This command outputs text only - it does NOT write any files.

## User Input
```
$ARGUMENTS
```

## Process

### 1. Gather Current Context

Identify relevant documents and state:
- Active story document (if any)
- Related design document
- Style guides in use
- Key files being worked on

### 2. Generate Handover Prompt

Output a prompt in this format that the user can copy/paste:

```
# Handover Context

## Current Task
[Brief description of what we were working on]

## Documents to Read
- [Path to active story document]
- [Path to related design]
- [Path to relevant style guides]
- [Other key files]

## Status Update
[Current progress on tasks, what's completed, what's in progress]

## Requirements & Decisions
[Key requirements identified during discussion]
[Important decisions made and rationale]

## Developer Logs
[Key entries from the story's developer logs — decisions, blockers encountered, deviations from design, lessons learned. Include anything the new agent needs to avoid re-discovering.]

## Next Steps
1. [First thing the new agent should do]
2. [Second step]
3. [Continue from here...]

## Notes
[Any blockers, concerns, or context the new agent needs]
```

### 3. Output Only

- Print the handover prompt to the conversation
- Do NOT write to any file
- User will copy/paste this into a new session

## Guidelines

- Be concise but include all critical context
- Reference specific file paths the new agent should read
- Summarize decisions so they don't need to be re-discussed
- Start planning mode first - don't jump into implementation
