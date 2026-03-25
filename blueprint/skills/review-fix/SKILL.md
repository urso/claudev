---
description: Review code and fix issues in a loop until clean (max 5 cycles)
user-invocable: true
disable-model-invocation: true
argument-hint: "[--story story-file] [--max-cycles N]"
allowed-tools: ["Read", "Write", "Edit", "Glob", "Grep", "Bash", "Task"]
---

# Review and Fix Loop

Run code review and fix issues iteratively until the code is clean or max cycles reached.

## Variables

- **DEFAULT_MAX_CYCLES**: 5

## User Input
```
$ARGUMENTS
```

Parse for:
- `--story <story-file>` for context (optional)
- `--max-cycles N` to override default (optional)

## Process

### 1. Initialize

Set cycle counter to 0. Determine max cycles (default 5 or from args).

### 2. Review-Fix Loop

```
while cycle < max_cycles:
    cycle++

    ## 2a. Run Review
    Spawn sub-agent to run /review with:
    - Current files (staged + unstaged)
    - --story if provided

    Wait for review results.

    ## 2b. Check Results
    If no errors and no warnings:
        → Exit loop, report success

    ## 2c. Fix Issues
    For each issue from review:
        - Read the file
        - Apply the fix
        - Follow style guides

    ## 2d. Verify Build
    Spawn sub-agent to run /develop-fix
    - Ensure fixes don't break build

    ## 2e. Continue
    Log: "Cycle {cycle}/{max_cycles} complete, {N} issues fixed"
```

### 3. Exit Conditions

**Success**: No issues found in review
```
## Review-Fix Complete

All checks pass after {cycle} cycle(s).

Files reviewed: X
Total issues fixed: Y
```

**Max cycles reached**: Still has issues
```
## Review-Fix Stopped

Reached max cycles ({max_cycles}).

Remaining issues:
[list of unfixed issues]

Run `/review-fix` again to continue, or fix manually.
```

## Sub-agent Prompts

### Review Sub-agent
```
Task: Run code review
Use the Skill tool to invoke "blueprint:review-code" with arguments: [--story if provided]
Return the full review output including all issues found.
```

### Fix Sub-agent
```
Task: Fix build/lint errors
Use the Skill tool to invoke "blueprint:develop-fix" with arguments: [--story if provided]
Return success/failure status.
```

## Key Guidelines

- **Fix only, never implement**: Only fix issues found in review - never add new functionality or implement story tasks
- **Minimal changes**: Apply the smallest fix needed to resolve each issue
- **Incremental fixes**: Fix all issues from one review before re-reviewing
- **Verify build**: Always run /develop-fix after applying review fixes
- **Track progress**: Log each cycle's results
- **Exit cleanly**: Report remaining issues if max cycles reached

## Important Constraints

**This command ONLY fixes review issues.** It must NOT:
- Implement new features or story tasks
- Add functionality beyond what's needed to fix the issue
- Refactor code unless the review explicitly flagged it
- Change code unrelated to review findings
