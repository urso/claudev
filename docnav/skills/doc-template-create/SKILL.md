---
description: >-
  Create or customize a project-local document template. Use this skill when the user asks to
  create a template ("create a template", "new template", "add a doc template", "make a template
  for X", "custom template", "scaffold a template"), wants to customize an existing bundled
  template ("customize the runbook template", "modify the note template", "override the default
  template"), wants to derive a template from existing docs ("make a template based on this doc",
  "turn this into a template", "extract a template from docs/foo.md", "create a template like
  our existing runbooks"), or needs to set up project-specific doc templates ("add project templates",
  "set up templates for this repo", "create a template that does X"). This skill guides
  template creation with frontmatter presets, agent-instructions, and per-section guidance,
  writing the result to the project-local templates directory.
user-invocable: true
argument-hint: "[template-name or 'from <bundled-name>' or 'from <doc-path>']"
allowed-tools: ["Bash", "Read", "Write", "Glob"]
---

# doc-template-create — Create or customize project-local templates

Guide the user through creating a markdown document template with embedded agent instructions.

## Variables

- **DOCNAV**: `${CLAUDE_PLUGIN_ROOT}/scripts/docnav.sh`
- **BUNDLED_TEMPLATES**: `${CLAUDE_PLUGIN_ROOT}/templates`

## Step 1: Determine starting point

There are three starting points:

1. **From a bundled template** — if `$ARGUMENTS` mentions a bundled template (e.g., "from runbook", "customize note"), list bundled templates and read the selected one:

   ```bash
   bash DOCNAV templates --templates-dir BUNDLED_TEMPLATES
   ```

2. **From existing docs** — if `$ARGUMENTS` points to a doc or docs (e.g., "based on docs/runbooks/deploy.md", "like our existing runbooks"), read them with the Read tool. Analyze the structure: extract the common section pattern, frontmatter fields, and writing style. Infer what the agent workflow should be — how would an agent guide someone to create another doc like this? Different docs imply different flows (e.g., a runbook needs step-by-step verification, an ADR needs context/decision/consequences). If multiple docs are referenced, identify the common structure across them.

3. **From scratch** — if `$ARGUMENTS` names a new template type (e.g., "incident", "adr") or is empty.

## Step 2: Gather template requirements

Ask the user:

1. **Template name** — the type name (e.g., `adr`, `incident`, `playbook`). This becomes the filename.
2. **Purpose** — what kind of document this template is for. This becomes the `scope` field.
3. **Sections** — what sections the template should have. Suggest defaults based on the purpose.
4. **Agent workflow** — how the agent should guide doc creation. This becomes `agent-instructions`. For simple templates, default to an interview-style workflow and skip detailed questions. For complex templates, ask:
   - Should the agent interview the user or draft content directly?
   - Should it check for duplicates first?
   - Any special steps (e.g., "suggest reviewers", "link to related docs")?

If starting from a bundled template, present its current structure and ask what to change.

If derived from existing docs, present the inferred structure (sections, frontmatter, proposed agent workflow) and ask the user to confirm or adjust. Pay attention to the flow — a template for incident reports needs a different agent workflow than one for design docs. Propose `agent-instructions` that match how the doc type is actually written: what questions to ask, what order, what to verify.

## Step 3: Build the template

Construct the template markdown with:

### Frontmatter

```yaml
---
title: ""
type: <template-name>
tags: []
scope: "<purpose from step 2>"
agent-instructions: |
  <workflow from step 2>
---
```

Include `watches: []` and `updates-when: []` fields only if relevant to the template type. These fields are used by the doc-stale skill: `watches` lists file paths that trigger staleness checks, and `updates-when` lists prose triggers the agent interprets against recent changes.

### Sections

For each section, add a heading and an `<!-- agent: ... -->` comment with per-section guidance:

```markdown
## Section Name

<!-- agent: Guidance for what to write in this section. -->
```

Mark optional sections with `<!-- agent: Optional — remove if not applicable. -->`.

## Step 4: Review and write

Present the complete template to the user for review. Make adjustments as requested.

Determine the project-local templates directory:

```bash
git rev-parse --show-toplevel
```

The target directory is `docs/ai/wiki/templates/` relative to git root. Create it if it doesn't exist.

Write the template to `<templates-dir>/<name>.md` using the Write tool.

Confirm by listing templates to verify it appears:

```bash
bash DOCNAV templates --templates-dir <project-local-dir>
```
