---
name: discover-exec-next
description: Internal agent for orchestrator — pick next open ticket and work it to completion
tools: [Bash, Read, Write, Glob, Grep, Skill]
---

# discover-exec-next

Autonomous ticket execution. Pick next open ticket, investigate, complete it.

## Mental Model

`.discovery/` is a knowledge graph, not a docs generator.
- Tickets with `done` status + populated findings = documentation
- No separate docs files — findings ARE the docs
- Children expand the graph with follow-up questions
- Each completed ticket adds to the project's documented knowledge

## Return Codes

Print exactly one at the end:
- `ok` — Completed a ticket, more may remain
- `done` — No open tickets left
- `blocked` — Ticket requires human input, marked as blocked

## Setup

```bash
DISCOVER="${CLAUDE_PLUGIN_ROOT}/scripts/discover.sh"
```

- **INVESTIGATION_GUIDE**: `${CLAUDE_PLUGIN_ROOT}/resources/investigation-guide.md` — how to investigate and record findings
- **TICKET_TEMPLATE**: `${CLAUDE_PLUGIN_ROOT}/resources/ticket-template.md` — structure for new/child tickets

## Process

### 1. Get next ticket

```bash
bash "$DISCOVER" ticket next --n 1
```

If output is empty or says "no tickets": print `done` and stop.

### 2. Parse and activate

Extract ticket ID from output. Read the ticket:

```bash
bash "$DISCOVER" ticket get <id>
```

Activate it:

```bash
bash "$DISCOVER" ticket update <id> --status active --log "Auto-started by exec-next"
```

### 3. Investigate

Read INVESTIGATION_GUIDE for approach. Based on ticket's intention and scope:

1. Search the codebase (grep, glob, read files)
2. Trace dependencies, call sites, data flow as needed
3. Build understanding incrementally

**Stay focused**: Answer the ticket's intention, nothing more.

### 4. Record findings

Update the ticket file directly. Add to `## Findings` section:
- Key discoveries with file paths and line numbers
- Code patterns found
- Architecture insights

Keep findings concise but complete — these become documentation.

### 5. Check for blockers

If investigation requires:
- External API access you don't have
- Human judgment on ambiguous code
- Access to systems outside the repo

Then:
```bash
bash "$DISCOVER" ticket update <id> --status blocked --log "Reason: <what's needed>"
```
Print `blocked` and stop.

### 6. Expand the knowledge graph

Run `/discover-expand <id>` to:
- Extract candidates from findings
- Link to existing tickets or create children
- Build graph connections

### 7. Complete

```bash
bash "$DISCOVER" ticket update <id> --status done --log "Investigation complete"
```

Print `ok` and stop.

## Constraints

- No user interaction
- One ticket per invocation
- Mark blocked rather than guessing
- Children inherit parent's tags unless specified
