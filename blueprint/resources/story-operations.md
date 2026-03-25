# Story Operations

How to create, modify, and manage stories. The template structure is in STORY_TEMPLATE within the plugin. Workflow guidelines can be loaded via `bash LIST_WORKFLOWS "" story`.

## Variables

- **RESOLVE_DIR**: `${CLAUDE_PLUGIN_ROOT}/scripts/resolve-dir.sh`
- **STORY_TEMPLATE**: `${CLAUDE_PLUGIN_ROOT}/templates/story-template.md`
- **NEXT_ID**: `${CLAUDE_PLUGIN_ROOT}/scripts/next-id.sh`
- **SET_STATUS**: `${CLAUDE_PLUGIN_ROOT}/scripts/set-status.sh`
- **SET_BLOCKED_BY**: `${CLAUDE_PLUGIN_ROOT}/scripts/set-blocked-by.sh`
- **REMOVE_BLOCKED_BY**: `${CLAUDE_PLUGIN_ROOT}/scripts/remove-blocked-by.sh`
- **QUERY_STORIES**: `${CLAUDE_PLUGIN_ROOT}/scripts/query-stories.sh`
- **LIST_WORKFLOWS**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-workflows.sh`

## Resolving the Stories Directory

Use `resolve-dir.sh` to get the absolute path (resolves via git root):
```bash
STORIES_DIR=$(bash RESOLVE_DIR story)
```

Always use this resolved path when writing files — never hardcode `docs/ai/stories/`.

## Creating a Story

1. Generate ID: `bash NEXT_ID story`
2. Read the template at STORY_TEMPLATE
3. Write to `$STORIES_DIR/<id>-<name>.md` (see naming convention below) using the template, filling in placeholders
4. Required frontmatter fields: `id`, `title`, `status`, `created`, `blocked-by`, `description`. Optional: `design` (parent design ID)
5. Must have at least one `## Tasks` or `## Phase N:` section with `- [ ]` checkboxes
6. Must have Developer Logs sections: Decision Log, Blockers Encountered, Deviations from Design, Lessons Learned
7. Set dependencies: `bash SET_BLOCKED_BY story <id> <blocker-ids...>`

## Modifying a Story

- **Status**: `bash SET_STATUS story <id> <status>`
- **Add blockers**: `bash SET_BLOCKED_BY story <id> <blocker-ids...>` (appends, deduplicates)
- **Remove blockers**: `bash REMOVE_BLOCKED_BY story <id> <blocker-ids...>`
- **Edit content**: Read the file, use Edit tool to modify sections
- **Relink to different design**: Edit the `design:` frontmatter field directly

## Cancelling a Story

1. Set status: `bash SET_STATUS story <id> cancelled`
2. Clean up references: remove this story's ID from other stories' `blocked-by` lists using `bash REMOVE_BLOCKED_BY`
3. Check if any stories were only blocked by this one — they may now be actionable

## Managing Dependencies

- Dependencies are directional: story A `blocked-by: [B]` means B must finish before A
- `bash QUERY_STORIES --actionable` returns stories that are `ready` or `in-progress` with all blockers resolved
- When adding stories between existing ones, update both the new story's `blocked-by` and downstream stories that should depend on it
- When removing a story from the chain, reconnect: if C was blocked by B, and B is removed, check if C should now be blocked by B's blockers

## Naming Convention

- With design: `<design-id>-<story-id>-<story-name>.md` (e.g. `0001-0003-auth-middleware.md`)
- Standalone: `<story-id>-<story-name>.md` (e.g. `0003-auth-middleware.md`)

Ask user for the kebab-case story name.
