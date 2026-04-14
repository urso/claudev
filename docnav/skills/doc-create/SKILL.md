---
description: >-
  Create a new markdown document from a template. Use this skill when the user asks to
  create a doc ("create a doc", "new doc", "write a runbook", "start a note", "add a design doc",
  "create documentation for X", "document X", "new runbook for Y", "write up X",
  "scaffold a doc", "draft a doc"),
  wants to document something ("I need to document X", "let's write docs for Y"),
  or mentions creating specific doc types by name (note, runbook, etc.). This skill
  lists available templates, guides the user through filling sections using the template's
  embedded agent instructions, and writes the final doc with instructions stripped.
user-invocable: true
argument-hint: "[template-type] [topic]"
allowed-tools: ["Bash", "Read", "Write", "Glob"]
---

# doc-create — Create a new document from a template

Guide the user through creating a markdown document using project templates with embedded agent instructions.

## Variables

- **DOCNAV**: `${CLAUDE_PLUGIN_ROOT}/scripts/docnav.sh`
- **BUNDLED_TEMPLATES**: `${CLAUDE_PLUGIN_ROOT}/templates`

## Step 1: Resolve template directories

Determine which template directories exist:

1. **Bundled**: `BUNDLED_TEMPLATES` (always available)
2. **Project-local**: `docs/ai/wiki/templates/` relative to git root

Check if the project-local directory exists with Glob. Use both dirs if present (project-local templates override bundled ones with the same name).

## Step 2: Select a template

If `$ARGUMENTS` names a known template type (e.g., "runbook", "note"), match it directly — skip listing.

Otherwise, list available templates:

```bash
bash DOCNAV templates --templates-dir BUNDLED_TEMPLATES
bash DOCNAV templates --templates-dir <project-local-dir>
```

Present the combined list (deduplicated, project-local wins) and ask the user to pick one.

If `$ARGUMENTS` contains both a template type and a topic (e.g., "runbook for deploying X"), extract both — use the type for template selection and carry the topic into the creation workflow.

## Step 3: Read the template

Use the `path` from the templates output (or construct it from the matched directory + `<name>.md`) and read the full template file with the Read tool.

Extract:
- **`agent-instructions`** frontmatter field — drives the overall workflow
- **`<!-- agent: ... -->` inline comments** — per-section guidance

## Step 4: Follow agent-instructions

Execute the workflow described in the template's `agent-instructions` field. This typically includes:

- Asking the user what they want to document (skip if topic was provided in `$ARGUMENTS`)
- Checking for duplicate docs (use `bash DOCNAV search "<topic>"` when instructed)
- Suggesting a title and tags
- Filling in frontmatter fields (`title`, `tags`, `scope`, etc.)

Adapt based on what the template instructions say — each template drives its own workflow.

## Step 5: Fill sections

Work through each section of the template:

- Follow `<!-- agent: ... -->` comments for per-section guidance (what to ask, whether required/optional)
- Ask the user for input on each section, or draft content based on context and confirm
- Remove sections marked optional if not applicable

## Step 6: Write the final document

Before writing, strip all agent instructions from the content:

1. Remove the `agent-instructions` field from frontmatter
2. Remove all `<!-- agent: ... -->` comments (single-line and multi-line)
3. Clean up any resulting blank lines (collapse multiple blank lines to one)

Ask the user for the file path, or suggest one based on the doc type and title (e.g., `docs/<title-slug>.md`).

Write the final cleaned document using the Write tool.
