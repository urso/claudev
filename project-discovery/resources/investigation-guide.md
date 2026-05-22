# Investigation Guide

How to investigate a discovery ticket once selected.

## Search Strategy

Start broad, narrow down:

1. **Grep for keywords** from ticket intention/scope
2. **Glob for file patterns** matching the domain (e.g., `**/auth/*.go`)
3. **Read entry points** — main files, handlers, exported APIs
4. **Trace inward** — follow imports, call sites, data flow

## What to Capture

Record in `## Findings` section:

- **File paths with line numbers** — `src/auth/jwt.go:42`
- **Key functions/types** — entry points, core abstractions
- **Data flow** — where data comes from, where it goes
- **Patterns** — conventions, repeated structures
- **Dependencies** — external packages, internal modules used

Keep concise. Link to code, don't copy it.

## When to Spawn Children

Create child ticket when you find:

- Related but out-of-scope question
- Deeper dive needed on a sub-component
- Follow-up work (refactor opportunity, tech debt)

**Broad scope tickets require children.** If you find 3+ subsystems in an overview ticket, each needs its own child ticket. Listing subsystems in findings without spawning children leaves the graph incomplete.

### Check Before Spawn

Before creating a child ticket, check what already exists:

```bash
bash "$DISCOVER" ticket list
bash "$DISCOVER" ticket recall "<subsystem/topic>"
```

The list shows all tickets. Recall scores them by relevance — but it uses stemming, so "zfs-of" and "NVMe-oF" may not match even though they're the same thing.

**Review both outputs.** If any ticket covers the subsystem:
- `done` — Reference in findings, don't duplicate
- `active`/`open` — Note it, add cross-reference
- Similar scope — Skip creating child

Only create if truly missing.

### Creating Children

Don't investigate children — just queue them with clear intention.

```bash
bash "$DISCOVER" ticket new --title "..." --intention "..." --parent <id>
```

## When a Ticket is Done

Mark `done` when:

- Ticket's intention is answered
- Findings section documents the answer
- Open questions are either resolved or captured as children

Mark `blocked` when:

- Need external system access
- Ambiguous code requires human judgment
- Missing context only humans can provide

## Recording Findings

Update ticket file directly. Structure:

```markdown
## Findings

<discoveries here>

## Open Questions

<unresolved items — move to children or resolve before done>

## Log

- YYYY-MM-DD: Status change reason
```

## Investigation Depth

Match depth to ticket scope:

- **Narrow scope** (single function/file) → trace completely
- **Medium scope** (module/package) → map structure, sample key paths
- **Broad scope** (system/architecture) → overview + spawn children for details
