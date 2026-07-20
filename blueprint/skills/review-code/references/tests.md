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

Use richer format for coverage gaps and weak assertions. Terse format for mechanical issues.

### Coverage Gaps / Weak Assertions — Medium Format

For issues where the test claims to verify something but doesn't actually catch failures:

```
[warning] path/to/file_test.ext:LINE

**Claims to test:** What the test name/structure suggests it verifies.
**Actually tests:** What the test actually exercises.
**Gap:** What failure mode slips through.
**Fix:** How to strengthen it.
```

### Mechanical Issues — Terse Format

For naming, hygiene, and straightforward anti-patterns:

```
[warning] path/to/file_test.ext:LINE
Issue: Description.
Suggestion: How to fix.
```

### Examples

Coverage gap (medium):
```
[warning] src/auth/login_test.go:45

**Claims to test:** User login with invalid password returns error.
**Actually tests:** Calls Login() and checks error is non-nil.
**Gap:** Doesn't verify the error type or message. A timeout, DB error, or rate-limit would also pass.
**Fix:** Assert specific error type: `assert.ErrorIs(err, ErrInvalidPassword)`.
```

Mechanical (terse):
```
[warning] src/api/handler_test.go:120
Issue: Hardcoded timestamp `2024-01-15` will cause flaky behavior over time.
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
