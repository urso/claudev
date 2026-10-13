# Finding Documents

Three tools for locating design and story documents.

## Variables

- **FIND_DOC**: `${CLAUDE_PLUGIN_ROOT}/bin/find-doc.sh`
- **QUERY_DESIGNS**: `${CLAUDE_PLUGIN_ROOT}/bin/query-designs.sh`
- **QUERY_STORIES**: `${CLAUDE_PLUGIN_ROOT}/bin/query-stories.sh`

## FIND_DOC <type> <id> [--dir DIR]

Direct lookup by numeric ID. Returns the absolute file path or exits 1.

## QUERY_DESIGNS [flags]

Search and filter designs. Returns: `ID | FILENAME | TITLE | STATUS | DESCRIPTION`

Flags:
- `--search "text"` — case-insensitive match against title and description
- `--status X` — filter by status (`draft`, `ready`, `in-progress`, `done`, `cancelled`, `archived`)

## QUERY_STORIES [flags]

Search and filter stories. Returns: `ID | FILENAME | TITLE | STATUS | BLOCKED-BY | DESCRIPTION`

Flags:
- `--search "text"` — case-insensitive match against title and description
- `--status X` — filter by status
- `--design D` — filter by parent design ID
- `--actionable` — only `ready` or `in-progress` stories with all blockers resolved

## Strategy

- Numeric input → `bash FIND_DOC` for direct lookup
- Descriptive text → `--search` to narrow results
- Combine filters when useful (e.g. `--status done --search "auth"`)
- Single result → use it; multiple → let user choose
- No results → broaden search or list all
