---
name: discover-init
description: Initialize project discovery — scan repo, gather intent, create initial tickets
user-invocable: true
disable-model-invocation: true
allowed-tools: [Bash, Read, Write, AskUserQuestion, Agent]
argument-hint: "[optional initial intent]"
---

# discover-init

Initialize the `.discovery/` directory and create initial tickets based on repo scan and user intent.

## Variables

- **DISCOVER**: `${CLAUDE_PLUGIN_ROOT}/scripts/discover.sh`
- **DOCNAV_PROBE**: `${CLAUDE_PLUGIN_ROOT}/resources/docnav-probe.md`
- **SCAN_OUTPUT**: `${CLAUDE_PLUGIN_ROOT}/resources/scan-output.md`
- **TICKET_TEMPLATE**: `${CLAUDE_PLUGIN_ROOT}/resources/ticket-template.md`

## Process

### 1. Probe for docnav

```bash
bash "${CLAUDE_PLUGIN_ROOT}/scripts/has-docnav.sh"
```

Report to user:
- `found` → "Docnav detected — Tier 2 search available"
- `not found` → "Docnav not detected — using Tier 1 search"

### 2. Run scan (Sub-agent)

Spawn the scan agent:
```
Agent({
  description: "Run repository scan",
  subagent_type: "discover-scan",
  prompt: "Scan this repository and return a summary."
})
```

The agent returns a 3-5 line summary. Display it to the user.

### 3. Gather intent

Use AskUserQuestion to understand the user's goals:

**Question 1 — Role**: "What's your role with this project?"
- Options: New team member, Contributor from another team, External reviewer, Other

**Question 2 — Depth**: "How deep do you need to go?"
- Options: High-level overview, Specific area deep-dive, Full architectural understanding, Other

**Question 3 — Scope**: "Any particular area of focus?"
- Options: Entire codebase, Specific directory/module (ask which), Following a bug/feature, Other

**Question 4 — Motivation**: "What's driving this investigation?"
- Options: Onboarding, Code review, Bug investigation, Feature planning, Other

If the user provided an initial intent argument, incorporate it — they may have already answered some questions.

### 4. Propose initial tickets

Based on scan results + intent, propose an initial ticket set. Typical pattern:
- 1 "Repo overview" ticket (scope: `.`)
- 2-4 scoped tickets based on intent (e.g., "Auth system", "Build pipeline", "Testing strategy")

Read TICKET_TEMPLATE for frontmatter and body structure.

Present proposals with:
- Title
- Scope
- Intention (the question to answer)

Ask user to confirm, modify, or reject each.

### 5. Create approved tickets

For each approved ticket:

```bash
bash "$DISCOVER" ticket new --title "<title>" --scope "<scope>" --intention "<intention>" --tag <tags>
```

Report created tickets with their IDs.

### 6. Ask about .discovery/ gitignore status

Use AskUserQuestion:

"How should `.discovery/` be handled in git?"
- **Commit as shared scratchpad**: Team shares investigation state
- **Add to .gitignore**: Private investigation, not committed
- **Leave alone**: Decide later

Act on the choice:
- Shared: ensure `.discovery/` is NOT in `.gitignore`
- Private: add `.discovery/` to `.gitignore` if not present
- Leave alone: do nothing

## Notes

- This is the only command that creates `.discovery/` if absent
- Other commands assume `.discovery/` exists and suggest `/discover-init` if not
- Safe to re-run — will skip scan if recent, ask about additional tickets
