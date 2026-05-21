---
name: discover-status
description: Show discovery ticket status overview and counts
user-invocable: true
disable-model-invocation: true
allowed-tools: [Bash, Read]
---

# discover-status

Read `${CLAUDE_PLUGIN_ROOT}/resources/cli-reference.md` for available operations.

```bash
DISCOVER="${CLAUDE_PLUGIN_ROOT}/scripts/discover.sh"
test -d .discovery && bash "$DISCOVER" ticket status || echo "Run /discover-init first"
```
