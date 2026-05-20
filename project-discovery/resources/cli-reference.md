# CLI Reference

```bash
DISCOVER="${CLAUDE_PLUGIN_ROOT}/scripts/discover.sh"
```

Run `bash "$DISCOVER" help <cmd>` for full details.

## Ticket Operations

```bash
bash "$DISCOVER" ticket recall "<intent>"     # recall-first search
bash "$DISCOVER" ticket status                # summary: counts + active/open lists
bash "$DISCOVER" ticket next [--n 3]          # suggest next ticket to work on
bash "$DISCOVER" ticket new --title "..." --intention "..." [--scope "..."] [--tag X] [--parent t-NNNN]
bash "$DISCOVER" ticket list [--status X] [--tag X]
bash "$DISCOVER" ticket get <id>
bash "$DISCOVER" ticket update <id> --status X --log "reason"
bash "$DISCOVER" ticket search "<query>" [--tag X] [--scope "..."]
bash "$DISCOVER" ticket find-overlap --intent "..."
bash "$DISCOVER" ticket tags
```

## Scan Operations

```bash
bash "$DISCOVER" scan repo    # structure signals → .discovery/scan.json
bash "$DISCOVER" scan churn   # code ownership → .discovery/scan.json
bash "$DISCOVER" scan schema  # show scan.json structure
```

## Statuses

`open` | `active` | `done` | `dropped` | `blocked`
