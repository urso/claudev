---
description: Load development context (style guides and build config)
user-invocable: true
disable-model-invocation: true
argument-hint: "[tech...] [build|style]"
allowed-tools: ["Read", "Bash"]
---

# Load Development Context

Load rules (style guides and build configuration) for development work.

## Variables

- **LIST_RULES**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-rules.sh`

## User Input
```
$ARGUMENTS
```

## Process

### 1. Parse Arguments

Parse `$ARGUMENTS` to determine what to load:
- If contains "build" → filter by tag "build"
- If contains "style" → filter by tag "style"
- Any other tokens (e.g., "go", "rust", "typescript") → treat as applies-to filters
- No arguments → load everything

Examples:
- `go` → Load rules that apply to 'go' (any tag)
- `build` → Only load rules with tag 'build'
- `style` → Only load rules with tag 'style'
- `go style` → Load style rules for 'go'
- `go build` → Load build rules for 'go'
- (empty) → Load all rules

### 2. List Rules

List available rules using the script:
```bash
bash LIST_RULES "" [applies-to-filter] [tag-filter]
```

Note: first argument is the rules directory (empty string uses default from git root). Filters go in positions 2 and 3.

Output format (grouped by directory, pipe-separated columns):
```
<directory>:
full/path/to/file.md|name|applies-to|tags|paths|description
```

### 3. Load Applicable Rules

For each rule in the list:
- Read the file content
- Keep in context for later use

### 4. Summarize Loaded Context

Provide a brief summary of what was loaded:
- Number of rules loaded by type (style, build, etc.)
- Key commands available (from build rules)
- List rule names with their applies-to and tags

## Key Guidelines

- **Flexible parsing**: Handle natural language arguments gracefully
- **Filter logic**: applies-to filter AND tag filter (if both specified)
- **Concise output**: Summarize what's loaded, don't dump entire files
- **Ready for work**: Context is now loaded and ready for implementation tasks
