# Acceptance Criteria Verification Instructions

Verify that acceptance criteria from the story are satisfied by the code changes.

## Input

You receive:
- A list of acceptance criteria from the story's tasks
- The files changed in this review
- Access to the full codebase for verification

## Verification Process

For each criterion:

1. **Understand what to check** — criteria are code-based assertions (interface exists, function has signature, test covers case)
2. **Locate evidence** — read files, grep for symbols, check test coverage
3. **Determine pass/fail** — be strict; partial implementations fail

## Report Format

Output format:
```
[pass] Task: <task name>
Criterion: <criterion text>
Evidence: <file:line or explanation>

[fail] Task: <task name>
Criterion: <criterion text>
Missing: <what's missing or wrong>
```

Example:
```
[pass] Task: Implement Processor interface
Criterion: `Processor` interface exists in `pkg/core/processor.go`
Evidence: pkg/core/processor.go:42-48

[fail] Task: Implement Processor interface
Criterion: Has `Process(ctx context.Context, input Input) (Output, error)` signature
Missing: Signature is `Process(input Input) error` — missing context parameter and Output return

[pass] Task: Implement Processor interface
Criterion: Unit test covers happy path
Evidence: pkg/core/processor_test.go:TestProcess_Success
```

## Guidelines

**Be strict:**
- Criterion says "interface X exists" → verify the exact name and location
- Criterion says "has signature Y" → match exactly, including parameter names if specified
- Criterion says "test covers Z" → verify the test actually tests that scenario

**Do NOT:**
- Assume something works without checking
- Pass criteria for "close enough" implementations
- Infer intent beyond what the criterion states

**Edge cases:**
- If a criterion is ambiguous, report as `[unclear]` with explanation
- If a criterion references code outside the diff, still verify it in the codebase
