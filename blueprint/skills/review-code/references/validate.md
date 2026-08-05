# Validation Instructions

Validate review findings by spawning sub-agents to score each issue's confidence,
then filter out low-confidence findings.

## Confidence Scale

Each validator assigns every issue a **confidence** score from 0 to 100 — how
certain it is that the issue is real, given the code it read.

| Range | Meaning |
|-------|---------|
| 90-100 | Confirmed. Traced the code; the issue provably manifests. |
| 75-89  | Likely real. Strong evidence, minor unverified assumptions. |
| 50-74  | Plausible but unproven. Depends on context the validator could not confirm. |
| 25-49  | Probably wrong. Evidence leans against it. |
| 0-24   | False positive. The code is correct, or the rule does not apply. |

**Threshold: keep issues with confidence >= 75.** Everything below is filtered.

Score the *issue*, not its severity. A cosmetic style nit that definitely applies
scores high; a catastrophic race condition that might not be reachable scores low.

Calibration rules for validators:
- Do not default to 75. If you did not read enough to be sure, the score is below 75.
- A duplicate of another issue in the list scores 0.
- An issue in unchanged / pre-existing code scores 0.
- Uncertainty about whether a code path is reachable caps the score at 60.

## Process

### Parse Review Output

Extract issues grouped by category:
- **Bug errors**: `[error]` items (have `Bug:` / `Impact:` fields)
- **Bug warnings**: `[warning]` items with `Bug:` / `Impact:` fields
- **Style warnings**: `[warning]` items with `Style:` / `Rule:` fields
- **Efficiency warnings**: `[warning]` items with `Efficiency:` / `Context:` / `Suggestion:` fields
- **Config issues**: `[error]` / `[warning]` items on config files with `Issue:` / `Context:` / `Impact:` fields

If no issues found, report "No issues to validate" and exit.

### Spawn Validation Sub-agents

Spawn validation sub-agents based on issue type. Run all validators in parallel.

Every validator returns **all** issues it was given, each with a `Confidence:` line
and a one-line `Reason:`. The orchestrator does the filtering — validators do not
drop issues themselves.

**For style warnings** (if any), spawn:
```
Task: Validate style issues
Subagent type: general-purpose
Model: haiku
Prompt:
You are scoring style review findings for confidence.

## Style Issues to Validate
[list style [warning] issues]

## Files Under Review
[file list]

## Confidence Scale
[paste the Confidence Scale section above]

## Instructions

1. For each unique rule file referenced (e.g., "Rule: common.md - ..."), read that rule file from docs/ai/rules/ to understand the full rule context.

2. For each issue, read the relevant code and assign a confidence score.

Score low when:
- The code actually follows the stated rule when considering full rule context
- The rule doesn't apply to this code pattern
- The issue duplicates another in the list

Return EVERY issue in this format, none omitted:
[warning] path/file.ext:LINE
Style: ...
Rule: ...
Confidence: NN
Reason: one line on what drove the score
```

**For bug warnings** (if any), spawn:
```
Task: Validate bug warnings
Subagent type: general-purpose
Model: sonnet
Prompt:
You are scoring bug review warnings for confidence.

## Bug Warnings to Validate
[list bug [warning] issues]

## Files Under Review
[file list]

## Confidence Scale
[paste the Confidence Scale section above]

## Instructions

For each warning, read the code with surrounding context and score how certain you are the concern is valid.

Score low when:
- The described scenario cannot actually occur given the code flow
- The code handles the case correctly when considering full context
- Language/framework guarantees prevent the issue
- The issue is in unchanged/pre-existing code

Return EVERY issue in this format, none omitted:
[warning] path/file.ext:LINE
Bug: ...
Impact: ...
Confidence: NN
Reason: one line on what drove the score
```

**For bug errors** (if any), spawn:
```
Task: Validate critical bugs
Subagent type: general-purpose
Model: opus
Prompt:
You are scoring critical bug findings for confidence. These require deep analysis - especially for race conditions, ownership issues, and memory safety.

## Critical Issues to Validate
[list bug [error] issues]

## Files Under Review
[file list]

## Confidence Scale
[paste the Confidence Scale section above]

## Instructions

For EACH error, perform thorough analysis before scoring:

1. Read the full file(s) involved, not just the flagged line
2. Trace data flow and control flow around the issue
3. For race conditions: analyze all access points, locks, synchronization
4. For ownership/lifetime issues: trace object creation, usage, and cleanup
5. For null/nil issues: check all paths that reach the flagged code
6. Consider framework/language guarantees that might prevent the issue

Score high (90+) when:
- You traced a concrete path where the bug manifests under realistic conditions
- No existing guards prevent the issue

Score low when:
- Deep analysis shows the issue cannot occur
- Synchronization/guards exist that the original review missed
- Language semantics guarantee safety
- The code path is unreachable

If you could not complete the trace, say so and score below 75 rather than guessing.

Return EVERY issue in this format, none omitted:
[error] path/file.ext:LINE
Bug: ...
Impact: ...
Confidence: NN
Analysis: Brief explanation of the trace and what drove the score.
```

**For config issues** (if any), spawn:
```
Task: Validate config issues
Subagent type: general-purpose
Model: sonnet
Prompt:
You are scoring configuration review findings for confidence.

## Config Issues to Validate
[list config issues]

## Files Under Review
[file list, including context files marked as unchanged]

## Confidence Scale
[paste the Confidence Scale section above]

## Instructions

For each issue, read the config with surrounding context and score it.

Score low when:
- The setting is defined elsewhere — a parent chart, a base overlay, a merged
  values file, a mutating admission controller
- The value is intentional and documented in the chart or repo
- The issue is in a file passed as unchanged context, not under review
- The template renders correctly despite looking suspicious — check the actual
  rendered output, not just the template source
- The permission or capability is genuinely required by the workload

Score high when:
- The issue manifests at render time or deploy time
- Embedded code has a real defect with a concrete trigger
- A security setting is missing with no upstream source

Return EVERY issue in this format, none omitted:
[error|warning] path/file.yaml:LINE
Issue: ...
Impact: ...
Confidence: NN
Reason: one line on what drove the score
```

### Filter and Merge Validated Results

Collect outputs from all validation agents. Then:

1. **Keep** issues with `Confidence >= 75`.
2. **Drop** issues below 75, but keep a count per category for the summary.
3. If a validator returns an issue with no `Confidence:` line, treat it as 50 — dropped.

Carry the confidence value into the final report so the user can judge borderline
findings.

### Present Validated Report

```
## Validation Summary

### Critical Bugs (N)
[kept error-level bug issues, each with its confidence]

### Bug Warnings (N)
[kept warning-level bug issues]

### Style Issues (N)
[kept style issues]

### Efficiency Issues (N)
[kept efficiency issues]

### Config Issues (N)
[kept config issues]

---
Initial findings: X issues
After validation: Y confirmed at confidence >= 75 (Z filtered below threshold)
```

If any dropped issue scored 65-74, list it under a `Borderline (not reported)`
heading with its title and score — near-misses are worth a glance, and hiding them
entirely loses information the validator paid for.
