# Design Operations

How to create, modify, and manage designs. The template structure is in DESIGN_TEMPLATE within the plugin. Workflow guidelines can be loaded via `bash LIST_WORKFLOWS "" design`.

## Variables

- **RESOLVE_DIR**: `${CLAUDE_PLUGIN_ROOT}/scripts/resolve-dir.sh`
- **DESIGN_TEMPLATE**: `${CLAUDE_PLUGIN_ROOT}/templates/design-template.md`
- **NEXT_ID**: `${CLAUDE_PLUGIN_ROOT}/scripts/next-id.sh`
- **SET_STATUS**: `${CLAUDE_PLUGIN_ROOT}/scripts/set-status.sh`
- **VALIDATE_DOC**: `${CLAUDE_PLUGIN_ROOT}/scripts/validate-doc.sh`
- **LIST_WORKFLOWS**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-workflows.sh`

## Resolving the Designs Directory

Use `resolve-dir.sh` to get the absolute path (resolves via git root):
```bash
DESIGNS_DIR=$(bash RESOLVE_DIR design)
```

Always use this resolved path when writing files — never hardcode `docs/ai/designs/`.

## Creating a Design

1. Generate ID: `bash NEXT_ID design`
2. Read the template at DESIGN_TEMPLATE
3. Write to `$DESIGNS_DIR/<id>-<name>.md` using the template, filling in placeholders
4. Required frontmatter fields: `id`, `title`, `status`, `created`, `description`. Optional: `depends-on` (list of design IDs), `references` (list of project-relative file paths for context)
5. Required sections: Problem Context, Goals, Non-Goals, Requirements, Technical Approach, Architecture Decisions
6. Optional sections: Open Questions, References

## Modifying a Design

- **Status**: `bash SET_STATUS design <id> <status>`
- **Edit content**: Read the file, use Edit tool to modify sections
- **Add dependency on another design**: Edit the `depends-on:` frontmatter field directly
- **Add references**: Edit the `references:` frontmatter field with project-relative file paths
- Keep designs lean — update approach and decisions, don't add implementation details or tasks

## Validating a Design

`bash VALIDATE_DOC design <file>` — checks frontmatter fields, valid status, existing `depends-on` references, `references` file paths, and empty sections.

## Status Lifecycle

`draft` → `ready` → `in-progress` → `done`

- `draft`: Design is being written, not ready for stories
- `ready`: Design is complete, stories can be created from it
- `in-progress`: Stories are being implemented
- `done`: All stories complete, design goals achieved
- `cancelled`: Design abandoned
- `archived`: Design and stories moved to archive

## Naming Convention

Files: `<id>-<kebab-case-name>.md` (e.g. `0001-user-auth.md`). Ask user for the name.
