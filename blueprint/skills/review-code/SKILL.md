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

## Process

### 1. Determine Files to Review

**If files specified:** Use those files.

**If no files specified:** Get staged and unstaged changes:
```bash
git diff --name-only HEAD
git diff --name-only --cached
```

**Filter out non-code files:**
- Skip: `*.md`, `*.json`, `*.yaml`, `*.yml`, `*.txt`, `*.lock`
- Skip: Files in `.git/`, `node_modules/`, `vendor/`, `dist/`, `build/`
- Include: Source code files

If no code files to review, report that and exit.

### 2. Load Story Context (if provided)

If `--story` specified, read the story document for implementation context.

### 3. Spawn Review Sub-agents

Read the relevant reference files for review instructions, then spawn sub-agents in parallel. Each agent receives the review instructions from the reference file, the file list, and story context if provided.

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

### 4. Collect and Merge Results

Wait for all agents to complete, then merge findings grouped by severity:
- **Bug errors**: Bug review [error] items
- **Bug warnings**: Bug review [warning] items
- **Style warnings**: Style review [warning] items
- **Efficiency warnings**: Efficiency review [warning] items

### 5. Validate for False Positives

If issues were found, read `references/validate.md` for validation instructions. Spawn validation sub-agents to check each issue category with appropriate model tiers (haiku for style, sonnet for bug warnings, opus for critical bugs).

### 6. Present Combined Report

```
## Code Review Summary

Reviewed X files (validated for false positives).

### Critical Bugs (N)
[validated error-level bug issues]

### Bug Warnings (N)
[validated warning-level bug issues]

### Style Issues (N)
[validated style issues]

### Efficiency Issues (N)
[validated efficiency issues]

### Clean Files
[files with no confirmed issues]

---
Initial findings: X issues
After validation: Y confirmed issues (Z filtered as false positives)
```

### 7. Propose Style Guide Additions (if issues found)

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
- **`references/validate.md`** — False positive validation with tiered model sub-agents

## Integration

- Part of `/review-fix` — automated review + fix loop
- Part of `/review-loop` — simplify, fix, and review-fix cycle
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
