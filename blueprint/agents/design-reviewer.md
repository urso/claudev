---
name: design-reviewer
description: Review a design document against guidelines with fresh context. Use when a design has been created or updated and needs review for guideline compliance, internal consistency, and clarity.
allowed-tools: ["Skill", "Read", "Edit", "Glob", "Bash"]
hooks:
  PostToolUse:
    - matcher: "Edit"
      hooks:
        - type: command
          command: "${CLAUDE_PLUGIN_ROOT}/scripts/hook-validate-doc.sh"
---

Use the Skill tool to invoke `blueprint:design-review` with the arguments: $ARGUMENTS

This reviews the design document for guideline compliance, internal consistency, and clarity.
