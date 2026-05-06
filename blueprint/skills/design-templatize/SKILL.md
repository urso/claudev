---
description: Extract a reusable template from a finished design
user-invocable: true
disable-model-invocation: true
argument-hint: "[design-id or file]"
allowed-tools: ["Read", "Write", "Edit", "Glob", "Bash", "AskUserQuestion"]
---

# Templatize Design

Extract a reusable template from a successful design. The template can be used as a starting point for similar future designs.

## Variables

- **DISCOVERY_GUIDE**: `${CLAUDE_PLUGIN_ROOT}/resources/discovery.md`
- **USER_TEMPLATES**: `docs/ai/templates/designs/`

## User Input
```
$ARGUMENTS
```

Parse for design file path, name, or ID.

## Process

### 1. Find Design

Read DISCOVERY_GUIDE for available tools and strategy, then locate the design based on user input. Read the design document.

### 2. Verify Design Status

The design should be `done` or `in-progress` with learnings. Warn if extracting from a `draft` design — it may not have proven patterns yet.

### 3. Gather Template Details

Analyze the design's problem context, goals, and approach. Propose 1-2 template names based on the pattern (e.g., if it's a migration design, suggest `data-migration` or `system-migration`).

Present proposals to user:
- **Proposed name(s)** — short kebab-case identifiers with rationale
- **Proposed description** — one-line summary of when to use this template

Ask user to confirm, adjust, or provide their own. Also ask:
- **What made this design successful?** — patterns worth preserving

### 4. Extract Template Structure

Analyze the design and identify:

- **Sections to keep**: Which sections are specific to this template type vs. standard?
- **Content to generalize**: Replace specific details with placeholders or guidance
- **Patterns to preserve**: Structural patterns that worked well

### 5. Create agent-instructions

Based on the design's structure and user feedback, draft `agent-instructions` that guide future use:

1. What questions to ask the user upfront
2. What context to gather before writing
3. Key decisions that need to be made
4. Order of filling in sections

Present the draft instructions to user for review.

### 6. Create Section Guidance

For each section, create `<!-- agent: ... -->` comments explaining:

- What this section should contain for this template type
- Key considerations or patterns to follow
- Examples or hints from the original design (generalized)

### 7. Write Template

Create the template directory and files:

```
docs/ai/templates/designs/<template-name>/
  template.md
  review.md (optional)
```

**template.md**:
- Frontmatter with `template: <name>`, `description`, and `agent-instructions`
- Sections with placeholder content and `<!-- agent: -->` guidance
- Remove all specific content from original design

Ask user if they want a custom `review.md` with template-specific review criteria.

### 8. Confirm

- Report template location
- Suggest testing with `/design-create` to verify the template works

## Guidelines

- Templates should be general enough to apply to similar problems
- Preserve structural patterns, not specific content
- agent-instructions should guide, not dictate — leave room for judgment
- Keep section guidance concise — hints, not essays
