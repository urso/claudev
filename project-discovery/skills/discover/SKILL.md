---
name: discover
description: Intent-driven discovery — recall existing knowledge or create new investigation tickets
user-invocable: true
disable-model-invocation: true
allowed-tools: [Bash, Read, Write, AskUserQuestion, Skill]
argument-hint: "<intent>"
---

# discover

Search for existing knowledge before creating new tickets. Uses recall-first: done → active → open → overlap → new.

## Mental Model

`.discovery/` is a knowledge graph, not a docs generator.
- Tickets with `done` status + populated findings = documentation
- No separate docs files — want reference material? Create a ticket for it
- When updating tickets, use the standard sections: Findings, Open questions, Log
- External docs join the graph via `doc-reference` tickets — summarize coverage, spawn children for gaps
- When findings describe structure (entities, layers, data flow), consider adding a diagram:
  - Mermaid source in ticket body (version-controlled, editable)
  - SVG in `.discovery/diagrams/<ticket-id>-<name>.svg` for local review

## Setup

```bash
DISCOVER="${CLAUDE_PLUGIN_ROOT}/scripts/discover.sh"
```

- **CLI_REFERENCE**: `${CLAUDE_PLUGIN_ROOT}/resources/cli-reference.md` — ticket and scan commands
- **TICKET_TEMPLATE**: `${CLAUDE_PLUGIN_ROOT}/resources/ticket-template.md` — structure for new tickets
- **INVESTIGATION_GUIDE**: `${CLAUDE_PLUGIN_ROOT}/resources/investigation-guide.md` — how to investigate and record findings

Read CLI_REFERENCE for available `ticket` commands (new, get, update, list, search, children, etc.).

## Process

### 1. Check for .discovery/

```bash
test -d .discovery && echo "exists" || echo "missing"
```

If missing, tell user to run `/discover-init` first and stop.

### 2. List and Recall

First list all tickets to see the full graph:

```bash
bash "$DISCOVER" ticket list
```

Then recall for relevance scoring:

```bash
bash "$DISCOVER" ticket recall "<user intent>"
```

Recall uses stemming — names like "zfs-of" and "NVMe-oF" tokenize differently and won't match. Review both outputs.

Output: first line is match type (`resolved`, `active`, `open`, `overlap`, `none`), followed by matching tickets if any.

### 3. Handle Result

**resolved** — Found a done ticket covering this area:
1. Summarize the ticket's findings
2. AskUserQuestion: "Found a resolved ticket. How to proceed?"
   - **Reopen**: `bash "$DISCOVER" ticket update <id> --status active --log "Reopened"`
   - **Spawn child**: Create new ticket with `--parent <id>`
   - **Start fresh**: Create unrelated new ticket

**active** — Found an in-progress investigation:
1. Summarize current state
2. AskUserQuestion: "Found an active ticket. How to proceed?"
   - **Resume**: Continue working on it
   - **Start fresh**: Create unrelated new ticket

**open** — Found a queued ticket:
1. Summarize the ticket
2. AskUserQuestion: "Found an open ticket. How to proceed?"
   - **Resume**: Activate and work on it
   - **Start fresh**: Create unrelated new ticket

**overlap** — No direct match but scope/tag overlap:
1. Show overlapping tickets
2. AskUserQuestion: "Found related tickets. How to proceed?"
   - **Expand existing**: Add to one of them
   - **Create new**: Create despite overlap

**none** — No existing knowledge:
1. Read TICKET_TEMPLATE for structure
2. Propose a new ticket (title, scope, intention)
3. On confirmation: `bash "$DISCOVER" ticket new --title "..." --scope "..." --intention "..." --tag discovery`

### 4. Investigate

Once a ticket is selected:

1. If `open`, activate: `bash "$DISCOVER" ticket update <id> --status active --log "Starting investigation"`
2. Read INVESTIGATION_GUIDE for search strategy and recording findings
3. Investigate the codebase based on the ticket's intention
4. Record findings in the ticket body (`## Findings` section)

### 5. Expand

Run `/discover-expand <id>` to extract candidates from findings and link or create children.

### 6. Completion

When investigation seems complete, AskUserQuestion: "Investigation complete. How to proceed?"
- **Done**: `bash "$DISCOVER" ticket update <id> --status done --log "Investigation complete"`
- **Keep open**: Continue investigation

**When to mark done vs keep open:**
- **Done**: The ticket's intention is answered. Findings section is populated. Open questions are either resolved or captured as child tickets.
- **Keep open**: Still actively investigating. Waiting on external input. Child tickets need completion before parent can be resolved.
- **Rule**: A parent can be marked done even if children are open — the parent's question is answered, children are follow-up work.

## Notes

- Log all status changes via `--log "..."`
- Never auto-close — always confirm with user
