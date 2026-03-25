# Archive Operations

## Archiving a Design

`bash ARCHIVE_DOC design <id>` — sets status to `archived` and moves to the archive directory.

- `bash ARCHIVE_DOC design <id> --with-stories` — archive design and all its stories
- `bash ARCHIVE_DOC design <id> --stories-only` — archive only the stories for this design

## Variables

- **ARCHIVE_DOC**: `${CLAUDE_PLUGIN_ROOT}/scripts/archive-doc.sh`
