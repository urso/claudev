# Bug and Logic Review Instructions

Review code for bugs, logic errors, security issues, and resource problems.

## Categories

- **Runtime:** nil/null deref, index OOB, div-by-zero, unhandled exceptions
- **Logic:** incorrect conditionals, off-by-one, races, deadlocks, infinite loops
- **Resources:** memory leaks, unclosed handles, missing cleanup
- **Security:** injection (SQL/command/XSS), hardcoded creds, weak crypto
- **Edge cases:** empty input, boundaries, error paths, concurrency

## Severity

- **[error]:** definite bugs, security vulns, data corruption risks
- **[warning]:** potential issues depending on usage, resource concerns, possible races

**Do NOT report:** style issues, speculation without clear trigger path, pre-existing bugs in unchanged code

## Report Formats

Use richer formats for complex bugs so the author can verify your understanding.

### [error] Findings — Rich Format

For significant bugs, show enough context that the author can:
1. Verify you understood the code correctly
2. Assess severity themselves
3. Understand exactly how to reproduce

Structure:

```
## Finding: <short descriptive title>

### What this code does
Explain the mechanism and its intended invariants. 1-2 paragraphs.

### Before the change
(Skip if not a regression.) How it worked previously. What property was preserved.

### After the change
What changed. Show the relevant code snippet. Explain the new behavior.

### Trigger scenario
Concrete sequence showing how the bug manifests:
    t0  action A
    t1  action B
    t2  bug triggers here
    t3  observable failure

### Severity assessment
- **Breaks:** what actually goes wrong (crash, data loss, hang, wrong result)
- **Doesn't break:** correct likely misconceptions
- **Trigger conditions:** always, edge case, requires malicious input, etc.
- **Calibration:** critical/major/minor — merge-blocking or not, and why

### Suggested fix
Direction for the fix. Include code snippet if non-obvious.
```

### [warning] Findings — Medium Format

For potential issues where impact depends on context:

```
[warning] path/to/file.ext:LINE

**Bug:** What the issue is.
**Context:** When/how this code runs; why the issue matters here.
**Impact:** What could go wrong, under what conditions.
**Fix:** How to address it.
```

### Trivial Findings — Terse Format

For obvious mechanical issues (textbook SQL injection, clear nil deref) where fix is self-evident:

```
[error] path/to/file.ext:LINE
Bug: Description.
Impact: Consequence.
```
