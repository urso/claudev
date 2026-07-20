# Efficiency Review Instructions

Review code changes for performance inefficiencies — unnecessary allocations, redundant work, poor data structure choices, and wasteful I/O patterns.

**Important:** Efficiency is context-dependent. Only flag issues where the cost is likely meaningful given the surrounding code. A small inefficiency on a cold path or with bounded small input is not worth reporting. Focus on patterns that could cause real performance problems at scale or under load.

## Variables

- **LIST_RULES**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-rules.sh`

## Process

### 1. Identify File Types

From the file list provided to you, note the file extensions (e.g., `.go`, `.ts`, `.py`).
These determine which rules apply.

### 2. Load Efficiency Rules (if any)

For each file type, run LIST_RULES with the applies-to filter:

```bash
bash LIST_RULES "" "go" efficiency    # for .go files
bash LIST_RULES "" "ts" efficiency    # for .ts files
bash LIST_RULES "" "*" efficiency     # universal rules (always load these)
```

Output format (pipe-separated): `filename|name|applies-to|tags|paths|description`

### 3. Read Rule Content

If any rules were found, you MUST read each rule file using the Read tool to get the actual requirements. The LIST_RULES output only shows metadata.

### 4. Review Categories

For each file, read the full file and look for the following categories of issues. For every potential issue, consider whether it actually matters given the context — is this a hot path? Is the data set large or unbounded? Is this called frequently?

**Memory & Allocations:**
- Unnecessary copies where references or pointers would suffice
- Allocations inside loops that could be hoisted or pre-allocated
- Missing capacity hints on collections when the size is known or estimable
- String building via repeated concatenation in loops

**Redundant Work:**
- Repeated computations that could be cached or moved out of loops
- Redundant conversions or serialization round-trips
- Multiple passes over the same data when one pass would do
- Unnecessary sorting or searching when a better data structure would avoid it

**I/O & System Resources:**
- Unbuffered I/O where buffering would help
- Missing connection or resource pooling
- N+1 query patterns (database or API calls in loops)
- Not leveraging available batch APIs

**Concurrency:**
- Lock scope broader than necessary (holding locks during I/O, etc.)
- Unnecessary synchronization on uncontended paths
- Blocking operations where async/buffered alternatives exist

**Data Structure Choice:**
- Linear search on large or unbounded collections where a map/set would work
- Using a complex structure where a simpler one fits (and vice versa)

## 5. Report Format

All efficiency issues are `[warning]` — suggestions, not bugs. Use richer format for issues where the performance impact isn't immediately obvious.

### Standard Format

```
[warning] path/to/file.ext:LINE

**Issue:** What the inefficiency is.
**Context:** Why it matters here — hot path? unbounded input? called per-request?
**Impact:** How bad is it? Rough cost estimate if possible (O(n²), extra allocation per call, lock contention under load).
**Severity:** minor/moderate/significant — and whether it's worth fixing now vs. later.
**Fix:** What to do instead.
```

### Example

```
[warning] src/sync/worker.go:112

**Issue:** Mutex held across network call (line 115).
**Context:** All workers contend on this lock. Network calls take 10-100ms typical.
**Impact:** Under load, workers serialize on network latency instead of just the state update. Throughput limited to ~10-100 ops/sec regardless of worker count.
**Severity:** significant if this path is hot; moderate otherwise.
**Fix:** Narrow lock scope — release before network call, re-acquire for state update.
```

## When NOT to Report

Skip if:
- Data size is small and bounded (iterating fixed config)
- Cold path (init, setup, teardown)
- Simpler code is clearer and perf difference is negligible
- Pre-existing in unchanged code
- Not confident the alternative is actually faster

When in doubt, don't report.
