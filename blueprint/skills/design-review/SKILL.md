---
description: Review a design document against guidelines
user-invocable: true
disable-model-invocation: true
argument-hint: "[design-id or file]"
context: fork
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
- **LIST_WORKFLOWS**: `${CLAUDE_PLUGIN_ROOT}/scripts/list-workflows.sh`

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

### 3. Load Design Guidelines

```bash
bash LIST_WORKFLOWS "" design
```

Read the listed workflow files for design guidelines.

### 4. Review Design

Review the design against the loaded guidelines. Check for:

1. **Guideline compliance**:
   - Has required sections (Problem Context, Goals, Non-Goals, Requirements, Technical Approach)
   - Uses concise bullet points, not long paragraphs
   - Is lean — no task lists or step-by-step implementation details
   - Has proper frontmatter fields (id, title, status, created, description)

2. **Internal consistency**:
   - Goals and requirements align
   - Non-goals don't contradict goals
   - Technical approach addresses the requirements
   - No conflicting statements

3. **Clarity**:
   - Requirements are unambiguous
   - Approach is concrete enough to derive stories from
   - No vague language without concrete meaning

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
