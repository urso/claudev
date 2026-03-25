---
name: design-review
description: Review a design document against guidelines with fresh context. Use when a design has been created or updated and needs review for guideline compliance, internal consistency, and clarity.
allowed-tools: ["Skill", "Read", "Edit", "Glob", "Bash"]
hooks:
  PostToolUse:
    - matcher: "Edit"
      hooks:
        - type: command
          command: "${CLAUDE_PLUGIN_ROOT}/scripts/hook-validate-doc.sh"
---

Run `/design-review $ARGUMENTS` to review the design document.
