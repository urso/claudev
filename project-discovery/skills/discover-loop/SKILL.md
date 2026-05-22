---
name: discover-loop
description: Loop through open discovery tickets until queue is empty
user-invocable: true
disable-model-invocation: true
allowed-tools: [Bash, Agent]
---

# discover-loop

Process all open tickets by spawning sub-agents. Each ticket gets isolated context.

**Important**: Run only ONE sub-agent at a time. The knowledge graph is walked sequentially — parallel agents would race on ticket selection and duplicate work.

## Setup

```bash
DISCOVER="${CLAUDE_PLUGIN_ROOT}/scripts/discover.sh"
```

## Process

### 1. Check state

```bash
bash "$DISCOVER" ticket status
```

If no open tickets, report "No open tickets" and stop.

### 2. Loop

Repeat until done:

1. Spawn sub-agent:
   ```
   Agent({
     description: "Execute next discovery ticket",
     subagent_type: "discover-exec-next",
     prompt: "Execute the next open ticket. Return only: ok, done, or blocked."
   })
   ```

2. Check result:
   - `ok` → Continue loop
   - `done` → Report "All tickets processed" and stop
   - `blocked` → Log which ticket blocked, continue loop

3. After each iteration: `bash "$DISCOVER" ticket status`

### 3. Summary

When loop ends, report:
- Tickets completed this session
- Tickets still blocked (if any)
- Tickets spawned as children (if any)

## Limits

Stop early if:
- 40 tickets processed (prevent runaway)
- 3 consecutive blocked results (needs human attention)

Report reason for early stop.
