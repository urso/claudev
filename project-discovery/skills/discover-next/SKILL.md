---
name: discover-next
description: Suggest next discovery ticket to work on
user-invocable: true
disable-model-invocation: true
allowed-tools: [Bash, Read]
---

# discover-next

Read `${CLAUDE_PLUGIN_ROOT}/resources/cli-reference.md` for available operations.

```bash
DISCOVER="${CLAUDE_PLUGIN_ROOT}/scripts/discover.sh"
test -d .discovery && bash "$DISCOVER" ticket next || echo "Run /discover-init first"
```

Run `/discover <intent>` to start working on a suggestion.
