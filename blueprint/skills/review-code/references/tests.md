# Test Code Review Instructions

Review test code for quality, coverage, and maintainability.

## Review Categories

**Coverage & Completeness:**
- Happy path tested
- Key error cases covered (prioritize likely/dangerous failures, not exhaustive)
- Boundary conditions for critical logic
- Don't flag missing tests for trivial error paths

**Test Quality:**
- Tests actually assert something meaningful
- Assertions match what the test name claims
- No tautological assertions (always true)
- Tests are independent (no order dependency)

**Test Hygiene:**
- Descriptive test names
- Good readability (table-driven/parameterized tests are fine)
- Mocks/stubs used appropriately (not over-mocked)

**Maintainability:**
- Prefer black-box testing — test via public interfaces
- Avoid testing implementation details that change frequently
- If testing internals: treat them as self-contained units with clear boundaries
- Tests shouldn't need updates when refactoring internals
- No hardcoded values that'll rot (timestamps, IDs)

**Anti-patterns:**
- Flaky tests (timing, randomness, external deps)
- Tests that test the framework, not the code
- Overly broad tests (testing too many things)
- Tests tightly coupled to implementation details

## Report Format

Output format:
```
[warning] path/to/file_test.ext:LINE
Issue: Description of the problem.
Suggestion: How to improve.
```

Example:
```
[warning] src/auth/login_test.go:45
Issue: Test name `TestLogin` doesn't describe what scenario is tested.
Suggestion: Rename to `TestLogin_InvalidPassword_ReturnsError` or use table-driven tests.

[warning] src/db/query_test.go:89
Issue: Test mocks the entire database layer, then asserts the mock was called — tests nothing real.
Suggestion: Test the actual query logic or use an integration test with a real database.

[warning] src/api/handler_test.go:120
Issue: Hardcoded timestamp `2024-01-15` will cause test to behave differently over time.
Suggestion: Use relative time or inject a clock.
```

## Severity

**Report as [warning]:**
- All test issues are warnings — tests don't break production

**Do NOT report:**
- Style issues unrelated to test quality (handled by style review)
- Missing tests for trivial getters/setters
- Missing tests for simple delegation methods
- Pre-existing issues in unchanged test code
