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

## Setup

Read `${CLAUDE_PLUGIN_ROOT}/resources/cli-reference.md` for available operations.

```bash
DISCOVER="${CLAUDE_PLUGIN_ROOT}/scripts/discover.sh"
```

## Process

### 1. Check for .discovery/

```bash
test -d .discovery && echo "exists" || echo "missing"
```

If missing, tell user to run `/discover-init` first and stop.

### 2. Recall Search

```bash
bash "$DISCOVER" ticket recall "<user intent>"
```

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
1. Propose a new ticket (title, scope, intention)
2. On confirmation: `bash "$DISCOVER" ticket new --title "..." --scope "..." --intention "..." --tag discovery`

### 4. Investigate

Once a ticket is selected:

1. If `open`, activate: `bash "$DISCOVER" ticket update <id> --status active --log "Starting investigation"`
2. Investigate the codebase based on the ticket's intention
3. Record findings in the ticket body (`## Findings` section)

### 5. Completion

When investigation seems complete, AskUserQuestion: "Investigation complete. How to proceed?"
- **Done**: `bash "$DISCOVER" ticket update <id> --status done --log "Investigation complete"`
- **Spawn child**: Create follow-up tickets with `--parent <id>`, then mark done
- **Keep open**: Continue investigation

## Notes

- Log all status changes via `--log "..."`
- Never auto-close — always confirm with user
