# Review Pipeline

Shared review process. The calling skill provides parsed chunks and optional flags.

## Filter Non-Code Files

From each chunk, filter out:
- `*.md`, `*.json`, `*.yaml`, `*.yml`, `*.txt`, `*.lock`
- Files in `.git/`, `node_modules/`, `vendor/`, `dist/`, `build/`
- Generated files — `*.gen.*`, `*_gen.*`, `*.pb.go`, `*.pb.ts`, `*_generated.*`, `*.generated.*`, files with `// Code generated` or `# Generated` header comments

If no code files remain after filtering, report that and exit.

## Load Story Context (if provided)

If story file specified, read it for:
- Implementation context for review agents
- Acceptance criteria from the tasks (for acceptance criteria sub-agent)

## Spawn Review Sub-agents

Each sub-agent reads its reference file directly — do NOT read the reference files yourself or paraphrase them into the prompt.

Unless a `--*-only` flag limits the scope, spawn all applicable review types.

### Per-Chunk Reviews

Spawn one agent **per chunk** for each review type, all in parallel:

**Style review per chunk** (skip if `--bugs-only` or `--efficiency-only`):
```
Task: Run style review (chunk N of M)
Model: sonnet

Read the style review instructions from ${CLAUDE_PLUGIN_ROOT}/skills/review-code/references/style.md.
Review these files: [chunk N files]
Story context: [if provided]

Follow the instructions to review each file and report issues.
```

**Efficiency review per chunk** (skip if `--style-only` or `--bugs-only`):
```
Task: Run efficiency review (chunk N of M)
Model: sonnet

Read the efficiency review instructions from ${CLAUDE_PLUGIN_ROOT}/skills/review-code/references/efficiency.md.
Review these files: [chunk N files]
Story context: [if provided]

Follow the instructions to review each file and report issues.
```

**Bug review per chunk** (skip if `--style-only` or `--efficiency-only`):
```
Task: Run bug review (chunk N of M)
Model: opus

Read the bug review instructions from ${CLAUDE_PLUGIN_ROOT}/skills/review-code/references/bugs.md.
Review these files: [chunk N files]
Story context: [if provided]

Follow the instructions to review each file and report issues.
```

**Test review per chunk** (skip if any `--*-only` flag, only for chunks containing test files):
```
Task: Run test review (chunk N of M)
Model: sonnet

Read the test review instructions from ${CLAUDE_PLUGIN_ROOT}/skills/review-code/references/tests.md.
Review these test files: [test files in chunk N]
Story context: [if provided]

Follow the instructions to review test quality and report issues.
```

### Bug Synthesis Pass (if multiple chunks)

After chunk-level bug reviews complete, spawn a synthesis agent:

```
Task: Bug review synthesis
Model: opus

You are reviewing findings from chunk-level bug reviews and checking for cross-file issues.

Chunk findings:
[collected bug findings from all chunk agents]

All changed files:
[full file list across all chunks]

Tasks:
1. Deduplicate any findings that appear in multiple chunks
2. Check for cross-file bugs the chunk agents couldn't see:
   - Data flow issues across file boundaries
   - Inconsistent error handling patterns
   - Race conditions involving multiple files
   - API contract violations between caller and callee
3. If you need to examine cross-file relationships, read the relevant files

Output: Combined bug findings (chunk findings + any new cross-file issues found)
```

### Acceptance Criteria Review

Acceptance criteria span the full change, so always pass all files regardless of chunk count.

**Acceptance criteria review** (only if story provided, skip if any `--*-only` flag):
```
Task: Verify acceptance criteria
Model: sonnet

Read the acceptance criteria instructions from ${CLAUDE_PLUGIN_ROOT}/skills/review-code/references/acceptance.md.
Review these files: [all files across all chunks]

Acceptance criteria to verify:
[extracted criteria from story tasks]

Follow the instructions to verify each criterion and report pass/fail.
```

## Collect and Merge Results

Wait for all agents to complete, then merge findings grouped by severity:
- **Bug errors**: Bug review [error] items
- **Bug warnings**: Bug review [warning] items
- **Style warnings**: Style review [warning] items
- **Efficiency warnings**: Efficiency review [warning] items
- **Test warnings**: Test review [warning] items
- **Acceptance criteria**: [pass]/[fail]/[unclear] items (if story provided)

## Validate for False Positives

If issues were found, read `${CLAUDE_PLUGIN_ROOT}/skills/review-code/references/validate.md` for validation instructions. Spawn validation sub-agents to check each issue category with appropriate model tiers (haiku for style, sonnet for bug warnings, opus for critical bugs).

## Present Combined Report

Assign typed IDs to each finding:
- **B1, B2, ...** — Bug findings
- **S1, S2, ...** — Style issues
- **E1, E2, ...** — Efficiency issues
- **T1, T2, ...** — Test issues
- **A1, A2, ...** — Acceptance criteria

**Preserve detail for complex findings.** Sub-agents produce rich output for bugs, efficiency, and test coverage gaps. Include that detail — don't flatten to one-liners.

```
## Code Review Summary

Reviewed X files in Y chunks (validated for false positives).

### Critical Bugs

#### B1. [error] <short title>
**What this code does:** ...
**Trigger:** t0 → t1 → t2 → failure
**Severity:** critical/major/minor — merge-blocking? why?
**Fix:** ...

### Bug Warnings

#### B3. [warning] path/to/file.ext:LINE
**Bug:** ...
**Impact:** ...
**Fix:** ...

### Style Issues
S1. path/to/file.ext:LINE — Description

### Efficiency Issues

#### E1. [warning] path/to/file.ext:LINE
**Issue:** ...
**Impact:** ...
**Fix:** ...

### Test Issues

#### T1. [warning] path/to/file_test.ext:LINE
**Claims to test:** ...
**Actually tests:** ...
**Gap:** ...

### Clean Files
[files with no confirmed issues]

---
Initial findings: X issues
After validation: Y confirmed (Z filtered as false positives)
```

Users can reference findings by ID: "fix B3", "ignore S1-S3", "explain E2".

## Report Acceptance Criteria (if story provided)

```
### Acceptance Criteria

#### Task 1: <Name>
A1. [x] Criterion — evidence
A2. [ ] Criterion — FAIL: what's missing
```

**After reporting**, update the story file to mark passing criteria.

## Propose Style Guide Additions

If validated issues reveal recurring patterns not yet documented, propose additions:

```
## Style Guide Proposal

**Guide**: [which style guide]
**Rule**: [short rule name]
**Description**: [what the rule captures]
**Rationale**: [why this matters]

Add to the style guide? (Run `/style-update` to add)
```
