# Bug and Logic Review Instructions

Deep review of code changes for bugs, logic errors, and potential issues.

## Review Categories

**Runtime Errors:**
- Nil/null pointer dereferences
- Index out of bounds
- Division by zero
- Unhandled exceptions

**Logic Errors:**
- Incorrect conditionals
- Off-by-one errors
- Race conditions
- Deadlocks
- Infinite loops

**Resource Issues:**
- Memory leaks
- Unclosed resources (files, connections)
- Missing cleanup

**Security Issues:**
- SQL injection
- Command injection
- XSS vulnerabilities
- Hardcoded credentials
- Insecure crypto

**Edge Cases:**
- Empty inputs
- Boundary conditions
- Error paths
- Concurrent access

## Report Issues

Output format:
```
[error] path/to/file.ext:LINE
Bug: Description of the issue.
Impact: What could go wrong.
```

Example:
```
[error] src/auth/login.go:45
Bug: Nil pointer dereference - `user` may be nil if FindUser returns no match.
Impact: Runtime panic when user not found.

[error] src/db/query.go:89
Bug: SQL injection - user input concatenated directly into query string.
Impact: Attacker could execute arbitrary SQL commands.

[warning] src/api/handler.go:120
Bug: Resource leak - HTTP response body not closed on error path.
Impact: Connection pool exhaustion under load.
```

## Severity Levels

**Report as [error]:**
- Definite bugs that will cause failures
- Security vulnerabilities
- Data corruption risks

**Report as [warning]:**
- Potential issues depending on usage
- Resource management concerns
- Race condition possibilities

**Do NOT report:**
- Style issues (handled by the style review)
- Speculative "might be a problem" without clear path
- Pre-existing issues in unchanged code
