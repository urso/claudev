---
description: Extract learnings from Claude Code session history.
user-invocable: true
user-invocable-only: true
argument-hint: "[--all | --since DATE]"
allowed-tools:
  - Bash
  - Read
  - Agent
---

# Retrospective — Extract Learnings from Session History

You are running the **regular retrospective** workflow. Mine Claude Code session history for patterns,
repeated feedback, and hard-won learnings that should be preserved.

## Variables

- **INIT_DB**: `scripts/init-db.sh`
- **DB_FILE**: `~/.claude/retro/retro.duckdb`

---

## Database Schema

### Tables

**`msgs`** — Main message table
| Column | Type | Description |
|--------|------|-------------|
| `uuid` | uuid | Primary key |
| `parentUuid` | uuid | Parent message (for threading) |
| `sessionId` | uuid | Session identifier |
| `timestamp` | timestamp | Message time |
| `type` | varchar | `user` or `assistant` |
| `cwd` | varchar | Working directory |
| `message` | struct | Message struct with `role`, `content` fields |
| `filename` | varchar | Source file path |

Access message fields: `message.role`, `message.content` (JSON)

**`coding_sessions`** — Aggregated coding activity
| Column | Type | Description |
|--------|------|-------------|
| `session_id` | uuid | Session identifier |
| `org` | varchar | GitHub org |
| `repo` | varchar | Repository name |
| `edits` | int64 | Edit tool calls |
| `writes` | int64 | Write tool calls |
| `extensions` | varchar[] | File extensions edited |
| `started` | timestamp | First message |
| `ended` | timestamp | Last message |

### Macros

| Macro | Returns | Description |
|-------|---------|-------------|
| `extract_text(message.content)` | varchar | Extract text from message content JSON |
| `thread_back(uuid, depth)` | table(uuid, parentUuid, sessionId, timestamp, message, depth) | Walk parent chain |
| `thread_forward(uuid, depth)` | table(uuid, parentUuid, sessionId, timestamp, message, depth) | Walk child chain |
| `search_user(terms, limit)` | table(candidate json) | FTS on user messages |
| `search_assistant(terms, limit)` | table(candidate json) | FTS on assistant messages |
| `search_messages(terms, limit)` | table(candidate json) | Combined search |

### Query Patterns

```sql
-- Filter by role
SELECT uuid, extract_text(message.content) FROM msgs WHERE type = 'user';

-- FTS search
SELECT * FROM search_user('wrong folder', 20);

-- Thread expansion
SELECT uuid, extract_text(message.content), depth FROM thread_back('uuid'::uuid, 5);
```

---

## Step 1 — Database Initialization

Run the init script with any user-provided flags from `$ARGUMENTS`:

```bash
bash INIT_DB $ARGUMENTS
```

The script handles:
- `--all`: Load full history
- `--since DATE`: Load since specific date (YYYY-MM-DD format)
- (default): Load since last run via watermark

**Important**: If the user provides a relative date like "last 4 weeks", "past 30 days", or "2 months ago",
convert it to YYYY-MM-DD format before calling the script. Use today's date to calculate.

If 0 rows loaded, warn the user and stop.

---

## Step 2 — Discovery Phase

Discovery follows a strict pipeline. Each stage feeds the next:

```
┌─────────┐     ┌─────────┐     ┌─────────┐     ┌──────────┐
│  Seeds  │ ──▶ │ Needles │ ──▶ │ Threads │ ──▶ │ Clusters │
└─────────┘     └─────────┘     └─────────┘     └──────────┘
   words          UUIDs        conversations     grouped
  to search      (hits)        (expanded)        threads
```

**Critical rule**: Never cluster needles directly. A needle is just a pointer — the thread is the actual conversation. Always expand before clustering.

### 2a. Load Seeds

Seeds are search terms that find entry points into valuable conversations.

**Load custom seeds first** (REQUIRED if file exists):

```bash
cat ~/.claude/retro/seeds.md 2>/dev/null || echo "No custom seeds"
```

Custom seeds take priority over defaults. Parse the file and use those terms.

**seeds.md format**:
```markdown
# Custom Seeds

## Corrections
- "that's not right"
- "/nota"

## Preferences
- "I always want"
```

**Default seeds** (fallback when no custom seeds):

| Category | Seeds |
|----------|-------|
| Decision points | `"Option A" OR "Option B" OR "Pros:" OR "Cons:"` |
| Corrections ack | `"You're right" OR "Good point" OR "I should have"` |
| User corrections | `"not what I" OR "don't do" OR "no," OR "wrong"` |
| User confirmations | `"perfect" OR "exactly" OR "yes that"` |

### 2b. Find Needles

**MANDATORY**: Query EVERY seed category before proceeding. Quality depends on coverage.

For each category in seeds.md (or defaults), run:

```bash
duckdb DB_FILE -c "LOAD fts; SELECT * FROM search_user('<terms>', 15);"
duckdb DB_FILE -c "LOAD fts; SELECT * FROM search_assistant('<terms>', 15);"
```

Use `search_user` for user seeds, `search_assistant` for assistant seeds.

**Filter noise**: Skip messages containing `task notification`, `/retro:spective`, `/init`.

### 2b-gate. Seed Coverage Check

**STOP**: Before proceeding to 2c, display this summary table:

```
| Seed Category | Source | Queries | Needles |
|---------------|--------|---------|---------|
| Corrections (user) | seeds.md | 1 | 12 |
| Confirmations (user) | seeds.md | 1 | 8 |
| Architecture (assistant) | seeds.md | 1 | 15 |
| ... | ... | ... | ... |
| **TOTAL** | | X | Y |
```

If any category shows 0 queries, you missed it. Go back and query it.

Aim for 30-50 total needles across all categories.

### 2c. Expand Needles to Threads

For EVERY needle, expand to get the full conversation:

```bash
duckdb DB_FILE -c "SELECT * FROM thread_back('UUID'::uuid, 5);"
duckdb DB_FILE -c "SELECT * FROM thread_forward('UUID'::uuid, 10);"
```

**Thread traversal helpers**:
- `thread_back(uuid, depth)` — walk parent chain (find conversation start)
- `thread_forward(uuid, depth)` — walk child chain (find resolution)

Extract the actual message content from expanded threads. This is what you'll cluster.

**Utilities**:
- `extract_text(content)` — extract text from message content JSON, filters tool_use/tool_result

### 2d. Cluster Threads

Now cluster the **expanded threads** (not raw needles).

Use the `classifier` agent (`.claude/agents/classifier.md`):

```
Agent({
  subagent_type: "retro:classifier",
  prompt: `Cluster these conversation threads by topic.

INPUT:
<threads>
${threadData}
</threads>

OUTPUT SCHEMA:
[{"topic": "short description", "items": [1, 2, ...], "sources": ["seed_category", "coding:repo"], "type": "correction|preference|architectural|confirmation"}]

Target 5-10 clusters. Return JSON only.`
})
```

### 2e. Present Clusters

Show clusters with thread summaries (not just needle text). Ask user which to explore.

### 2f. Coding Session Discovery

For sessions with code changes, query the `coding_sessions` table:

| Column | Type | Description |
|--------|------|-------------|
| `session_id` | uuid | Session identifier |
| `org` | varchar | GitHub org |
| `repo` | varchar | Repository name |
| `edits` | int64 | Edit tool calls |
| `extensions` | varchar[] | File extensions edited |

```bash
duckdb DB_FILE -c "
  SELECT session_id, repo, edits, extensions
  FROM coding_sessions
  WHERE edits > 5
  ORDER BY edits DESC
  LIMIT 20;
"
```

For interesting sessions, search for needles within that session, then expand as above.

---

## Step 3 — Classification Phase

For each cluster the user selected:

### 3a. Refine Thread Boundaries

If threads seem to span multiple topics, use `classifier` agent:

```
Agent({
  subagent_type: "retro:classifier",
  prompt: `Detect topic boundaries in this conversation.

INPUT:
${messages}

OUTPUT SCHEMA:
{"boundary": "relevant|topic_changed", "trim_at": "uuid or null"}

Return JSON only.`
})
```

### 3b. Classify Threads

Use `classifier` agent for each thread:

```
Agent({
  subagent_type: "retro:classifier",
  prompt: `Classify this conversation thread.

INPUT:
${threadText}

CATEGORIES:
- correction: user corrected Claude's approach
- preference: user stated a preference
- architectural: user made a design decision
- frustration: user expressed repeated frustration
- confirmation: user confirmed a good approach

OUTPUT SCHEMA:
{"type": "...", "summary": "one-line description", "context": "why this matters"}

Return JSON only.`
})
```

**Track seed attribution**: Note which seed term found each thread. This enables self-improvement in Step 7.

---

## Step 4 — Generalization Phase

Before review, generalize findings to extract reusable principles.

### 4a. Generalize Findings

Use `reasoner` agent (Sonnet) — generalization requires nuanced abstraction:

```
Agent({
  subagent_type: "retro:reasoner",
  prompt: `Extract the general, reusable principle from this specific finding.

INPUT:
Specific finding: ${summary}
Context: ${context}
Project: ${project}

OUTPUT SCHEMA:
{
  "specific": "the concrete finding",
  "general": "the reusable principle",
  "scope": "cross-project|project-specific|context-dependent",
  "confidence": 0.0-1.0
}

Scope meanings:
- cross-project: applies everywhere
- project-specific: tied to this codebase's conventions
- context-dependent: needs more instances to generalize

Return JSON only.`
})
```

### 4b. Handle Low-Confidence Findings

If `confidence < 0.6` or `scope == "context-dependent"`:
- Don't promote to memory yet
- Store in staging area: `~/.claude/retro/staging/<project>/`
- These await reinforcement from future sessions

---

## Step 5 — Deduplication Phase

Before storing, check for semantic overlap with existing memories.

### 5a. Load Existing Memories

```bash
find ~/.claude/retro/memory -name "*.md" -exec cat {} \; 2>/dev/null
```

### 5b. Check for Duplicates

Use `reasoner` agent (Sonnet) — semantic comparison requires judgment:

```
Agent({
  subagent_type: "retro:reasoner",
  prompt: `Compare this new finding against existing memories. Determine overlap.

NEW FINDING:
${generalizedPrinciple}

EXISTING MEMORIES:
${existingMemorySummaries}

OUTPUT SCHEMA:
{
  "action": "new|merge|reinforce|skip",
  "existing_match": "memory-name or null",
  "reason": "why this decision"
}

Actions:
- new: no overlap, create new memory
- merge: overlaps but adds nuance, merge into existing
- reinforce: semantically identical, bump reinforcement count
- skip: already fully captured

Return JSON only.`
})
```

### 5c. Execute Action

Based on dedup result:
- **new**: proceed to review
- **merge**: show user both, propose merged version
- **reinforce**: increment count in existing memory, update `last_seen`
- **skip**: inform user, don't re-store

---

## Step 6 — Review Workflow

Present generalized findings for user approval.

### 6a. Present Findings

For each finding from Step 4, display:

```
## Finding: [type]

**Specific**: [original finding]
**Generalized**: [extracted principle]
**Scope**: [cross-project / project-specific]
**Confidence**: [0.0-1.0]

### Thread Preview
[Show 3-5 key messages from expanded thread]

Action: [approve / reject / skip / edit]?
```

User can **edit** to refine the generalization before storing.

### 6b. Collect Decisions

Track user decisions:
- **approve** — store as-is
- **reject** — discard
- **skip** — move to staging for later
- **edit** — user provides refined generalization

### 6c. Generate Summary

After all findings reviewed, compile approved ones:

```markdown
## Findings from /retro session

### New Memories
- [generalized principle] (scope, confidence)

### Reinforced
- [existing memory] — seen again (+1, total: N)

### Staged (awaiting reinforcement)
- [low-confidence finding] — needs more instances
```

### 6d. Save to Memory

**Memory location**: `~/.claude/retro/memory/<project>/`

**Memory file format**:
```markdown
---
name: <type>-<slug>
description: <generalized principle — one line>
metadata:
  type: <feedback|user|project>
  source: retro
  scope: <cross-project|project-specific>
  reinforcement_count: 1
  first_seen: <date>
  last_seen: <date>
  ttl_days: 90
---

<generalized principle>

**Origin**: <specific finding that led to this>
**Why:** <why this matters>
**How to apply:** <when/where this applies>
```

Map finding types to memory types:
- `correction` → `feedback`
- `preference` → `feedback`
- `architectural` → `project`
- `confirmation` → `feedback`
- `frustration` → `feedback`

Create directory if needed. Update `~/.claude/retro/memory/<project>/MEMORY.md` index.

---

## Step 7 — Memory Decay

Memories fade if not reinforced.

### 7a. Check for Stale Memories

At start of each retrospective session, scan for expired memories:

```bash
find ~/.claude/retro/memory -name "*.md" -mtime +90
```

For each memory file, check:
- `last_seen` + `ttl_days` < today → candidate for decay

### 7b. Decay Rules

| reinforcement_count | ttl_days | action on expiry |
|---------------------|----------|------------------|
| 1                   | 90       | delete           |
| 2-3                 | 180      | demote to staging |
| 4+                  | 365      | flag for review  |

Higher reinforcement = longer TTL. Frequently validated memories persist.

### 7c. Present Decay Candidates

```
## Memories approaching expiry

These haven't been reinforced recently:

- [memory name] — last seen 85 days ago, expires in 5 days
  Action: [keep / let expire / refresh]?
```

- **keep**: reset TTL, keep in memory
- **let expire**: allow natural decay
- **refresh**: user confirms still valid, bump `last_seen`

---

## Step 8 — Self-Improvement

After review, analyze patterns:

1. **Seed effectiveness**: Which seeds led to approved vs rejected findings?
2. **New phrases**: Extract common phrases from approved findings that aren't in current seeds

Use `classifier` agent to identify new seed candidates:

```
Agent({
  subagent_type: "retro:classifier",
  prompt: `Extract distinctive phrases from these approved findings that could be new search seeds.

INPUT:
${approvedFindings}

OUTPUT SCHEMA:
[{"phrase": "...", "category": "corrections|preferences|architectural|confirmation", "reason": "..."}]

Return 3-5 phrases. JSON only.`
})
```

If new seeds found, present to user:

```
## Suggested New Seeds

Based on your approved findings:
- "actually I meant" → corrections (appeared 3x in approved)
- "going forward" → preferences (appeared 2x in approved)

Add these to ~/.claude/retro/seeds.md? [y/n]
```

If approved, append to `seeds.md` (create if doesn't exist).
