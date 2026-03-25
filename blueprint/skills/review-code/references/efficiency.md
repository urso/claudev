# Efficiency Review Instructions

Review code changes for performance inefficiencies — unnecessary allocations, redundant work, poor data structure choices, and wasteful I/O patterns.

**Important:** Efficiency is context-dependent. Only flag issues where the cost is likely meaningful given the surrounding code. A small inefficiency on a cold path or with bounded small input is not worth reporting. Focus on patterns that could cause real performance problems at scale or under load.

## Variables

- **LIST_RULES**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-rules.sh`

## Process

### Load Efficiency Rules (if any)

```bash
bash LIST_RULES "" "" efficiency
```

Output format (pipe-separated): `filename|name|applies-to|tags|paths|description`

Read all applicable efficiency rules (those with `applies-to: ["*"]` or matching file types).

### Review Categories

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

## Report Issues

All efficiency issues are reported as `[warning]` — these are suggestions, not bugs.

Output format:
```
[warning] path/to/file.ext:LINE
Efficiency: Description of the inefficiency.
Context: Why this matters here (e.g., called per-request, unbounded input, hot loop).
Suggestion: What to do instead.
```

Example:
```
[warning] src/api/handler.go:85
Efficiency: Slice appended in loop without pre-allocation — results fetched from DB with known count.
Context: Called per API request; result set can be large.
Suggestion: Pre-allocate with make([]T, 0, count).

[warning] src/sync/worker.go:112
Efficiency: Mutex held across network call on line 115.
Context: All workers contend on this lock; network latency makes this a bottleneck.
Suggestion: Narrow lock scope to protect only the shared state update after the call.
```

## When NOT to Report

Do not flag an issue if:
- The data size is small and bounded (e.g., iterating a fixed config list)
- The code is in initialization, setup, or teardown (cold path)
- The simpler code is clearer and the performance difference is negligible
- The issue is in unchanged/pre-existing code
- You're not confident the alternative would actually be faster in practice

When in doubt, don't report it.
