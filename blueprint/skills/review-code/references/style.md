# Style Guide Review Instructions

Review code changes for compliance with project style guides.

## Variables

- **LIST_RULES**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-rules.sh`

## Process

### Load Style Rules

```bash
bash LIST_RULES "" "" style
```

Output format (pipe-separated): `filename|name|applies-to|tags|paths|description`

Read all applicable style rules (those with `applies-to: ["*"]` or matching file types).

### Review Each File

For each file:
1. Read the file content
2. Check against all applicable style guide rules
3. Note any violations

### Report Issues

Output format:
```
[warning] path/to/file.ext:LINE
Style: Description of the violation.
Rule: Reference to style guide rule if applicable.
```

Example:
```
[warning] src/auth/login.go:72
Style: Use structured logging with slog instead of fmt.Printf.
Rule: common.md - "Use structured logging"

[warning] src/api/handler.go:15
Style: Error messages should start with lowercase.
Rule: go.md - "Error strings should not be capitalized"
```

## Review Standards

**Report:**
- Clear violations of documented style guide rules
- Inconsistencies with established patterns

**Do NOT report:**
- Subjective preferences not in style guides
- Pre-existing issues in unchanged code
- Speculative "might be better" suggestions
