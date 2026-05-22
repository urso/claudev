---
name: discover-expand
description: Expand a ticket by linking or creating children for subsystems found in findings
user-invocable: true
allowed-tools: [Bash, Read, Write]
argument-hint: "<ticket-id>"
---

# discover-expand

Analyze a ticket's findings, extract candidates, and either link to existing tickets or create new children.

## Setup

```bash
DISCOVER="${CLAUDE_PLUGIN_ROOT}/scripts/discover.sh"
```

## Process

### 1. Get the ticket

```bash
bash "$DISCOVER" ticket get <ticket-id>
```

If not found, report error and stop.

### 2. Extract candidates

Read the ticket's `## Findings` section. List candidates for child tickets:
- Subsystems mentioned
- Components referenced
- Follow-up questions implied
- Areas marked as "out of scope" or "needs deeper investigation"

Be specific. "Authentication" is too broad. "JWT validation in auth middleware" is a candidate.

### 3. Check existing children

```bash
bash "$DISCOVER" ticket children <ticket-id>
```

Remove candidates already covered by children.

### 4. List all tickets

```bash
bash "$DISCOVER" ticket list
```

Keep this list for matching.

### 5. For each remaining candidate

**Score with recall:**
```bash
bash "$DISCOVER" ticket recall "<candidate topic>"
```

**Decide — prefer linking over creating:**

- **Match found** (same scope, or high relevance, or similar title):
  
  **Link to existing ticket** — this is the preferred outcome:
  1. Read the existing ticket file
  2. Add `<ticket-id>` to its `parents:` frontmatter array
  3. Write the updated file
  4. Add wikilink to current ticket's Relations section using full filename: `[[t-0002-nvme-of-target-zfs-of]]` (without `.md`)
  5. Note: "Linked to existing t-XXXX"
  
  Linking builds graph edges without duplication. Always prefer linking when a relevant ticket exists.

- **No match** — only then create:
  ```bash
  bash "$DISCOVER" ticket new --title "<candidate>" --intention "<what to investigate>" --parent <ticket-id>
  ```
  Note: "Created t-XXXX"

### 6. Report

List all actions taken:
- Linked: t-XXXX ← candidate topic
- Created: t-XXXX — candidate topic

## Notes

- Parent ticket status unchanged (stays done if done)
- New children start as `open`
- Use `/discover-loop` to process new children
