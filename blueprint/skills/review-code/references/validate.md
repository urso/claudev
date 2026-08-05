# Validation Instructions

Validate review findings by spawning sub-agents to check each issue category for false positives.

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

**For style warnings** (if any), spawn:
```
Task: Validate style issues
Subagent type: general-purpose
Model: haiku
Prompt:
You are validating style review findings for false positives.

## Style Issues to Validate
[list style [warning] issues]

## Files Under Review
[file list]

## Instructions

1. For each unique rule file referenced (e.g., "Rule: common.md - ..."), read that rule file from docs/ai/rules/ to understand the full rule context.

2. For each issue, read the relevant code and determine if it's a TRUE or FALSE positive.

FALSE POSITIVE if:
- The code actually follows the stated rule when considering full rule context
- The rule doesn't apply to this code pattern
- The issue duplicates another in the list

Return ONLY the true positive issues in the original format:
[warning] path/file.ext:LINE
Style: ...
Rule: ...

If all issues are false positives, return: "No style issues."
```

**For bug warnings** (if any), spawn:
```
Task: Validate bug warnings
Subagent type: general-purpose
Model: sonnet
Prompt:
You are validating bug review warnings for false positives.

## Bug Warnings to Validate
[list bug [warning] issues]

## Files Under Review
[file list]

## Instructions

For each warning, read the code with surrounding context and determine if the concern is valid.

FALSE POSITIVE if:
- The described scenario cannot actually occur given the code flow
- The code handles the case correctly when considering full context
- Language/framework guarantees prevent the issue
- The issue is in unchanged/pre-existing code

Return ONLY the true positive issues in the original format:
[warning] path/file.ext:LINE
Bug: ...
Impact: ...

If all issues are false positives, return: "No bug warnings."
```

**For bug errors** (if any), spawn:
```
Task: Validate critical bugs
Subagent type: general-purpose
Model: opus
Prompt:
You are validating critical bug findings. These require deep analysis - especially for race conditions, ownership issues, and memory safety.

## Critical Issues to Validate
[list bug [error] issues]

## Files Under Review
[file list]

## Instructions

For EACH error, perform thorough analysis:

1. Read the full file(s) involved, not just the flagged line
2. Trace data flow and control flow around the issue
3. For race conditions: analyze all access points, locks, synchronization
4. For ownership/lifetime issues: trace object creation, usage, and cleanup
5. For null/nil issues: check all paths that reach the flagged code
6. Consider framework/language guarantees that might prevent the issue

FALSE POSITIVE if:
- Deep analysis shows the issue cannot occur
- Synchronization/guards exist that the original review missed
- Language semantics guarantee safety
- The code path is unreachable

TRUE POSITIVE if:
- The bug can actually manifest under realistic conditions
- No existing guards prevent the issue

Return ONLY true positive issues in format:
[error] path/file.ext:LINE
Bug: ...
Impact: ...
Analysis: Brief explanation of why this is a real issue.

If all issues are false positives, return: "No critical bugs confirmed."
```

**For config issues** (if any), spawn:
```
Task: Validate config issues
Subagent type: general-purpose
Model: sonnet
Prompt:
You are validating configuration review findings for false positives.

## Config Issues to Validate
[list config issues]

## Files Under Review
[file list, including context files marked as unchanged]

## Instructions

For each issue, read the config with surrounding context and determine if it is valid.

FALSE POSITIVE if:
- The setting is defined elsewhere — a parent chart, a base overlay, a merged
  values file, a mutating admission controller
- The value is intentional and documented in the chart or repo
- The issue is in a file passed as unchanged context, not under review
- The template renders correctly despite looking suspicious — check the actual
  rendered output, not just the template source
- The permission or capability is genuinely required by the workload

TRUE POSITIVE if:
- The issue manifests at render time or deploy time
- Embedded code has a real defect with a concrete trigger
- A security setting is missing with no upstream source

Return ONLY the true positive issues in the original format:
[error|warning] path/file.yaml:LINE
Issue: ...
Impact: ...

If all issues are false positives, return: "No config issues."
```

### Merge Validated Results

Collect outputs from all validation agents. Combine the confirmed issues.

### Present Validated Report

```
## Validation Summary

### Critical Bugs (N)
[validated error-level bug issues]

### Bug Warnings (N)
[validated warning-level bug issues]

### Style Issues (N)
[validated style issues]

### Efficiency Issues (N)
[validated efficiency issues]

### Config Issues (N)
[validated config issues]

---
Initial findings: X issues
After validation: Y confirmed issues (Z filtered as false positives)
```
