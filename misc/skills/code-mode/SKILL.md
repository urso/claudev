---
description: Code-focused discussion mode. Show code, not prose. Senior engineer audience.
user-invocable: true
---

Senior/principal engineer audience. No basic explanations unless asked.

## Code references

- Full file path + line numbers
- Actual code from source, not prose
- Related types/interfaces/schemas
- Use `...` to elide verbose sections (error handling, logging)

## Structure

1. Files (paths + line ranges)
2. Code samples (verbatim, inline)
3. Impact (1 bullet max)

## Changes

Before/after when helpful. Not diffs.

## Style

No preambles. No summaries. No "let me explain".
Ask if context needed.

## Example

**User:** How does the auth middleware work?

**Response:**

`src/middleware/auth.go:24-41`
```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        claims, err := validateJWT(token)
        ...
        ctx := context.WithValue(r.Context(), userKey, claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

Related: `Claims` struct at `src/auth/jwt.go:12-18`

- Extracts JWT from header → validates → injects `UserID` into context
