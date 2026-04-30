# Style Guide Review Instructions

Review code changes for compliance with project style guides.

## Variables

- **LIST_RULES**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-rules.sh`

## Process

### 1. Identify File Types

From the file list provided to you, note the file extensions (e.g., `.go`, `.ts`, `.py`).
These determine which rules apply.

### 2. Load Style Rules

For each file type, run LIST_RULES with the applies-to filter:

```bash
bash LIST_RULES "" "go" style    # for .go files
bash LIST_RULES "" "ts" style    # for .ts files
bash LIST_RULES "" "*" style     # universal rules (always load these)
```

Output format (pipe-separated): `filename|name|applies-to|tags|paths|description`

### 3. Read Rule Content

**Critical:** The LIST_RULES output only shows metadata. You MUST read each applicable rule file using the Read tool to get the actual requirements.

For each rule file path from step 2:
1. Read the full file content
2. Note the specific conventions, patterns, and requirements
3. Use these when reviewing code in step 4

Without reading the rule files, you cannot check compliance.

### 4. Review Each File

For each file:
1. Read the file content
2. Check against all applicable style guide rules (from the rule files you read)
3. Note any violations

### 5. Report Issues

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
