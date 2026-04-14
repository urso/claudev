---
description: >-
  Find and review stale documentation that may need updating. Use this skill when the user
  asks about outdated docs ("are any docs stale", "which docs need updating", "doc freshness",
  "check for outdated documentation", "what docs are affected by recent changes", "doc drift",
  "are my docs up to date", "documentation audit"), wants to know the doc impact of code
  changes ("I changed X, what docs need updating", "which docs should I review before merging",
  "did I forget to update any docs", "doc impact of my changes", "update docs for this PR"),
  or needs to identify documentation that has fallen behind the codebase. Also use when the
  user asks to review or verify docs after refactoring, renaming, or removing code. This skill
  wraps the docnav stale command — always prefer it over manual git-diff-and-grep when
  checking doc freshness.
user-invocable: true
argument-hint: "[range]"
allowed-tools: ["Bash", "Read", "Glob", "Agent"]
---

# doc-stale — Find and review stale documentation

Surface docs whose watched files have changed in recent commits, then guide review and revision.

## Variables

- **DOCNAV**: `${CLAUDE_PLUGIN_ROOT}/scripts/docnav.sh`

## Command

```bash
bash DOCNAV stale --since <commit> --path <root>
```

- `--since <commit>` — git commit to diff against. The CLI defaults to `HEAD~1`, but `HEAD~10` gives more useful coverage.
- `--path <dir>` — root directory (defaults to git root)
- `--json` — structured output

## Workflow

### 1. Find stale docs

Run the stale command. Default to `--since HEAD~10` for a useful range unless the user specifies something.

Interpret `$ARGUMENTS` as plain text and translate to a `--since` value:

| User says | `--since` value |
|-----------|----------------|
| *(empty)* | `HEAD~10` |
| `last 20 commits` | `HEAD~20` |
| `since v1.2.0` | `v1.2.0` |
| `abc123` | `abc123` |
| `last week`, `since Monday` | Use `git log --since="1 week ago" --format=%H \| tail -1` to find the commit |

```bash
bash DOCNAV stale --since <resolved-commit> --path <root>
```

### 2. Present results

Output is tab-separated: `path  title` followed by indented `updates-when:` triggers for each doc.

Present each stale doc with:
- File path and title
- The `updates-when` triggers that explain *why* this doc may be stale
- A brief note on which watched files changed (if useful context)

If no stale docs are found, say so and suggest trying a broader `--since` range.

### 3. Guide review

Present the list first and let the user choose which docs to review. Do not read or revise docs unprompted.

When the user selects a doc for review, delegate to a sub-agent to keep the main context lean:

```
Agent({
  description: "Review stale doc: <title>",
  model: "sonnet",
  prompt: "Review the doc at <path> for staleness.\n\nContext: the doc's `watches` field matched recent git changes. The `updates-when` triggers are:\n<triggers>\n\nRecent changes (git log):\n<relevant changes summary>\n\n1. Read the doc\n2. Check whether the triggers are actually relevant to the recent changes\n3. If the doc needs updating, propose specific edits with exact old/new text\n4. If the doc is still current despite the trigger, explain why no changes are needed\n\nReport: state whether the doc needs updates, and if so, list the specific edits."
})
```

After the sub-agent returns, present its findings to the user. If edits are proposed, confirm before applying them.
