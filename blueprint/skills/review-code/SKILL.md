---
description: This skill should be used when the user asks to "review my code", "check for bugs", "check style", "review changes", "run code review", "check for issues", "check for performance issues", or needs code reviewed for style compliance, bugs, logic errors, or efficiency. Also use when the user mentions reviewing staged changes, checking code quality, or finding problems in code.
user-invocable: true
disable-model-invocation: false
argument-hint: "[files] [--story story-file] [--style-only] [--bugs-only] [--efficiency-only]"
allowed-tools: ["Read", "Grep", "Glob", "Bash(git:*)", "Task"]
---

# Code Review

Comprehensive code review checking style guide compliance, bugs/logic errors, and performance inefficiencies. Runs specialized reviews in parallel, then validates findings for false positives.

## Variables

- **LIST_RULES**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-rules.sh`

## User Input
```
$ARGUMENTS
```

Parse for:
- Specific files to review (optional, defaults to staged+unstaged changes)
- `--story <story-file>` for context (optional)
- `--style-only` to run only style review
- `--bugs-only` to run only bug review
- `--efficiency-only` to run only efficiency review
- Branch comparison: "for PR", "PR changes" → compare against main; "against <branch>", "vs <branch>" → compare against that branch

## Pre-computed Context

### Changed Files (unstaged)
!`git diff --name-only HEAD`

### Changed Files (staged)
!`git diff --name-only --cached`

## Process

### 1. Determine Files to Review

**If files specified:** Use those files.

**If branch comparison requested:** Run `git diff <branch>...HEAD --name-only` to get all files changed on this branch.

**If no files specified:** Use the changed files from the pre-computed context above.

**Filter out non-code files:**
- Skip: `*.md`, `*.json`, `*.yaml`, `*.yml`, `*.txt`, `*.lock`
- Skip: Files in `.git/`, `node_modules/`, `vendor/`, `dist/`, `build/`
- Skip: Generated files — `*.gen.*`, `*_gen.*`, `*.pb.go`, `*.pb.ts`, `*_generated.*`, `*.generated.*`, files with `// Code generated` or `# Generated` header comments
- Include: Source code files

If no code files to review, report that and exit.

### 2. Load Story Context (if provided)

If `--story` specified, read the story document for:
- Implementation context for other review agents
- Acceptance criteria from the tasks (for the acceptance criteria sub-agent)

### 3. Spawn Review Sub-agents

Spawn sub-agents in parallel. Each sub-agent reads its reference file directly — do NOT read the reference files yourself or paraphrase them into the prompt. The reference files contain detailed format instructions that must be followed exactly.

Unless a `--*-only` flag limits the scope, spawn all applicable agents.

**Style review** (skip if `--bugs-only` or `--efficiency-only`):
```
Task: Run style review
Model: sonnet

Read the style review instructions from references/style.md.
Review these files: [file list]
Story context: [if provided]

Follow the instructions to review each file and report issues.
```

**Bug review** (skip if `--style-only` or `--efficiency-only`):
```
Task: Run bug review
Model: opus

Read the bug review instructions from references/bugs.md.
Review these files: [file list]
Story context: [if provided]

Follow the instructions to review each file and report issues.
```

**Efficiency review** (skip if `--style-only` or `--bugs-only`):
```
Task: Run efficiency review
Model: sonnet

Read the efficiency review instructions from references/efficiency.md.
Review these files: [file list]
Story context: [if provided]

Follow the instructions to review each file and report issues.
```

**Test review** (skip if `--style-only` or `--bugs-only` or `--efficiency-only`, only if test files in file list):
```
Task: Run test review
Model: sonnet

Read the test review instructions from references/tests.md.
Review these test files: [test files from file list]
Story context: [if provided]

Follow the instructions to review test quality and report issues.
```

**Acceptance criteria review** (only if `--story` provided, skip if any `--*-only` flag):
```
Task: Verify acceptance criteria
Model: sonnet

Read the acceptance criteria instructions from references/acceptance.md.
Review these files: [file list]

Acceptance criteria to verify:
[extracted criteria from story tasks]

Follow the instructions to verify each criterion and report pass/fail.
```

### 4. Collect and Merge Results

Wait for all agents to complete, then merge findings grouped by severity:
- **Bug errors**: Bug review [error] items
- **Bug warnings**: Bug review [warning] items
- **Style warnings**: Style review [warning] items
- **Efficiency warnings**: Efficiency review [warning] items
- **Test warnings**: Test review [warning] items (if test files reviewed)
- **Acceptance criteria**: [pass]/[fail]/[unclear] items (if story provided)

### 5. Validate for False Positives

If issues were found, read `references/validate.md` for validation instructions. Spawn validation sub-agents to check each issue category with appropriate model tiers (haiku for style, sonnet for bug warnings, opus for critical bugs).

### 6. Present Combined Report

Assign typed IDs to each finding for easy reference in follow-up discussion:
- **B1, B2, ...** — Bug findings (both errors and warnings)
- **S1, S2, ...** — Style issues
- **E1, E2, ...** — Efficiency issues
- **T1, T2, ...** — Test issues
- **A1, A2, ...** — Acceptance criteria (pass/fail/unclear)

**Preserve detail for complex findings.** Sub-agents produce rich output for bugs, efficiency, and test coverage gaps. Include that detail in the report — don't flatten to one-liners. The goal is for the author to verify the reviewer understood the code.

```
## Code Review Summary

Reviewed X files (validated for false positives).

### Critical Bugs

#### B1. [error] <short title>

**What this code does:** ...
**Before:** ... (if regression)
**After:** ...
**Trigger:** t0 → t1 → t2 → failure
**Severity:** critical/major/minor — merge-blocking? why?
**Fix:** ...

#### B2. [error] ...

### Bug Warnings

#### B3. [warning] path/to/file.ext:LINE
**Bug:** ...
**Context:** ...
**Impact:** ...
**Fix:** ...

### Style Issues
S1. path/to/file.ext:LINE — Description
S2. ...

### Efficiency Issues

#### E1. [warning] path/to/file.ext:LINE
**Issue:** ...
**Context:** ...
**Impact:** ...
**Severity:** ...
**Fix:** ...

### Test Issues

#### T1. [warning] path/to/file_test.ext:LINE
**Claims to test:** ...
**Actually tests:** ...
**Gap:** ...
**Fix:** ...

(Use terse format for mechanical test issues: `T2. path:LINE — Description`)

### Clean Files
[files with no confirmed issues]

---
Initial findings: X issues
After validation: Y confirmed issues (Z filtered as false positives)
```

Users can reference findings by ID: "fix B3", "ignore S1-S3", "explain E2".

### 7. Report Acceptance Criteria (if story provided)

Include acceptance criteria results in the report with IDs:

```
### Acceptance Criteria

#### Task 1: <Name>
A1. [x] `Processor` interface exists in `pkg/core/processor.go` — found at line 42
A2. [ ] Has `Process(ctx, input) (output, error)` signature — FAIL: signature is `Process(input) error`
A3. [x] Unit test covers happy path — `processor_test.go:TestProcess`

#### Task 2: <Name>
A4. [x] ... — evidence
```

**After reporting**, update the story file:
- Mark passing criteria as `[x]`
- Leave failing criteria as `[ ]` with a comment explaining what's missing

### 8. Propose Style Guide Additions (if issues found)

If validated issues reveal recurring patterns not yet documented in style guides, propose additions:

```
## Style Guide Proposal

**Guide**: [which style guide, e.g., go.md, common.md]
**Rule**: [short rule name]
**Description**: [what the rule captures]
**Rationale**: [why this matters, based on issues found]

Add to the style guide? (Run `/style-update` to add)
```

Only propose if the issue represents a generalizable, actionable pattern not already covered.

## Reference Files

Detailed review instructions for each category:
- **`references/style.md`** — Style guide compliance checks, output format, review standards
- **`references/bugs.md`** — Bug/logic/security review categories, severity levels
- **`references/efficiency.md`** — Performance review categories, context-dependent filtering
- **`references/tests.md`** — Test code quality, coverage, maintainability
- **`references/acceptance.md`** — Acceptance criteria verification from story tasks
- **`references/validate.md`** — False positive validation with tiered model sub-agents

## Integration

- Part of `/review-fix` — automated review + fix loop
- Part of `/review-loop` — tidy, fix, and review-fix cycle
- Part of `/develop-loop` — full development workflow
- Use before committing to catch issues early

## Usage Examples

```
/review-code                       # Review all staged/unstaged code files
/review-code src/auth/             # Review specific directory
/review-code --story 0001-*-setup  # Review with story context
/review-code --style-only          # Only check style compliance (faster)
/review-code --bugs-only           # Only check for bugs (opus)
/review-code --efficiency-only     # Only check for inefficiencies
```
