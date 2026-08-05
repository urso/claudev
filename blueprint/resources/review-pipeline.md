# Review Pipeline

Shared review process. The calling skill provides parsed chunks and optional flags.

## Read the Sections

The diff script classifies files and emits four sections. Do not re-classify —
use the groups as given.

- **`## Code`** — pre-chunked by diff size. Feed to the per-chunk reviews below.
- **`## Config`** — grouped by Helm chart, or `loose` for non-chart configs.
  Lines marked `context:` are unchanged files supplied for reference only.
- **`## Docs`** — not reviewed. Report them so the user can ask for a review.
- **`## Ignored`** — dropped, with a reason. Report them so a misclassification
  is visible and correctable.

Each file line is `path <total> +<added> -<deleted>`, where total is the sum that
drives chunking. A delete-heavy file is a removal to check for orphaned
references; an add-heavy one is new logic to audit.

Files carry a status tag when they are not a plain modification — `(added)`,
`(deleted)`, `(renamed)`, `(copied)`. The `+N -M` split is omitted for adds and
deletes, where it says nothing the tag doesn't. Pass these through to the review
agents:

- **`(deleted)`** — the file is gone. Do not try to read it. Review the deletion
  itself: dangling references, removed cleanup, a template whose values are now
  orphaned.
- **`(added)`** — no prior version. Regression comparisons do not apply.
- **`(renamed)`** — the path shown is the new one.

If both Code and Config are absent, report that and exit. If one is absent, run
the other and skip its sections.

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

### Config Review

Skip if no config files, or if any `--*-only` flag was given.

Spawn **one agent for all config files**. Config diffs are small, and the groups
are context for the agent — not a fan-out key.

```
Task: Run config review
Model: sonnet — use opus if any group contains embedded shell, RBAC, or securityContext

Read the config review instructions from ${CLAUDE_PLUGIN_ROOT}/skills/review-code/references/config.md.

Review these files, grouped by chart. Each group's context files are unchanged —
use them to resolve value references, but do NOT report issues in them.

[paste the `## Config` section verbatim — groups, context lines, and all]

Story context: [if provided]

Follow the instructions to review each file and report issues.
```

Config files also get a style pass — `style.md` picks up embedded languages, so a
ConfigMap holding shell is checked against the project's `shell` rules. Include
config files in the style review chunks.

## Collect and Merge Results

Wait for all agents to complete, then merge findings grouped by severity:
- **Bug errors**: Bug review [error] items
- **Bug warnings**: Bug review [warning] items
- **Style warnings**: Style review [warning] items
- **Efficiency warnings**: Efficiency review [warning] items
- **Test warnings**: Test review [warning] items
- **Config errors**: Config review [error] items
- **Config warnings**: Config review [warning] items
- **Acceptance criteria**: [pass]/[fail]/[unclear] items (if story provided)

## Validate for False Positives

If issues were found, read `${CLAUDE_PLUGIN_ROOT}/skills/review-code/references/validate.md` for validation instructions. Spawn validation sub-agents to check each issue category with appropriate model tiers (haiku for style, sonnet for bug warnings, opus for critical bugs).

## Present Combined Report

Assign typed IDs to each finding:
- **B1, B2, ...** — Bug findings
- **S1, S2, ...** — Style issues
- **E1, E2, ...** — Efficiency issues
- **T1, T2, ...** — Test issues
- **C1, C2, ...** — Config issues
- **A1, A2, ...** — Acceptance criteria

**Preserve detail for complex findings.** Sub-agents produce rich output for bugs, efficiency, and test coverage gaps. Include that detail — don't flatten to one-liners.

```
## Code Review Summary

Reviewed X code files in Y chunks and Z config files (validated for false positives).

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

### Config Issues

#### C1. [error] path/to/values.yaml:LINE
**Issue:** ...
**Context:** ...
**Impact:** ...
**Fix:** ...

### Clean Files
[files with no confirmed issues]

### Not Reviewed
docs/design.md — docs
go.sum — lockfile
internal/api/api.pb.go — generated

Ask to review any of these explicitly if the classification is wrong.

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
