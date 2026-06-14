---
name: design-reviewer
description: Review a design document against guidelines with fresh context. Use when a design has been created or updated and needs review for guideline compliance, internal consistency, and clarity.
allowed-tools: ["Read", "Edit", "Glob", "Bash"]
hooks:
  PostToolUse:
    - matcher: "Edit"
      hooks:
        - type: command
          command: "${CLAUDE_PLUGIN_ROOT}/scripts/hook-validate-doc.sh"
---

# Review Design

Validate a design against guidelines and check for internal consistency.

## Variables

- **DISCOVERY_GUIDE**: `${CLAUDE_PLUGIN_ROOT}/resources/discovery.md`
- **LIST_TEMPLATES**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-templates.sh`
- **DEFAULT_REVIEW**: `${CLAUDE_PLUGIN_ROOT}/templates/designs/default/review.md`

## User Input
```
$ARGUMENTS
```

Parse for design file path, name, or ID.

## Process

### 1. Find Design

Read DISCOVERY_GUIDE for available tools and strategy, then locate the design based on user input.

### 2. Validate Structure

```bash
bash VALIDATE_DOC design <design-file>
```

If validation errors, report them immediately.

### 3. Load Review Criteria

1. Read DEFAULT_REVIEW — base criteria for all designs
2. Check design's `template:` frontmatter field
3. If template specified, run `bash LIST_TEMPLATES design` to find it, then read its `review.md` (replace `template.md` with `review.md` in path)

### 4. Review Design

Apply the loaded review criteria. For each checklist item, verify the design passes.

For each issue found, explain:
- Which section has the issue
- What the problem is
- Suggested fix

### 5. Present Results

Display findings to the user.

If issues found:
- List each issue with suggested fixes
- Ask if user wants to fix them now
- If yes, apply the fixes directly

If passed:
- Suggest next steps: `/design-expand` to break into stories

## Guidelines

- Review focuses on guideline compliance and internal consistency
- Does not judge whether the design is a good idea, only whether it's well-formed
- Fixes are applied directly when user approves
