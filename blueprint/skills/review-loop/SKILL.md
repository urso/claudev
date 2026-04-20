---
description: Tidy, fix build, and review-fix code in a loop until clean (no implementation)
user-invocable: true
disable-model-invocation: true
argument-hint: "[--story story-file] [--max-cycles N] [--no-tidy]"
allowed-tools: ["Read", "Glob", "Bash", "Task"]
---

# Review Loop

Post-implementation review workflow: tidy code, fix build, then review+fix cycle until clean. Use this when code is already implemented and you only need the review pass.

## User Input
```
$ARGUMENTS
```

Parse for:
- `--story <story-file>` for context (optional)
- `--max-cycles N` to override review-fix max cycles (optional, default 5)
- `--no-tidy` to skip the tidy step (optional)

## Process

### 1. Tidy Code (Sub-agent)

Skip if `--no-tidy` was specified.

Spawn sub-agent to tidy the implemented code:
```
Task: Tidy code for clarity
Use the Skill tool to invoke "blueprint:tidy" with arguments: [--story if provided]
Return summary of changes made.
```

Wait for completion.

### 2. Fix Build After Tidy (Sub-agent)

Skip if step 1 was skipped.

Spawn sub-agent to fix any build/lint errors introduced by tidying:
```
Task: Fix build errors after tidy
Use the Skill tool to invoke "blueprint:develop-fix" with arguments: [--story if provided]
Return build status (pass/fail).
```

Wait for completion.

### 3. Review-Fix Cycle (Sub-agent)

Spawn sub-agent to run review and fix loop:
```
Task: Review and fix code issues
Use the Skill tool to invoke "blueprint:review-fix" with arguments: [--story if provided] [--max-cycles N if provided]
This only fixes review issues - it does NOT implement new functionality.
Return final status and any remaining issues.
```

Wait for completion (max cycles internally).

### 4. Present Results

Report:
```
## Review Loop Complete

### Tidy
[Summary of changes, or "Skipped"]

### Build Status
[Pass / Fail]

### Review Status
[Pass / N remaining issues]
```

If issues remain:
```
### Remaining Issues
[List from review-fix output]

### Next Steps
- Run `/review-loop` again to continue
- Or fix manually and run `/review-code` to verify
```

## Notes

- All steps run automatically without user interaction
- This command never implements new features — it only tidies, fixes, and reviews
- Use `--no-tidy` if you've already tidied or just want the review-fix cycle
