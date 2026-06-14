---
name: story-reviewer
description: Review a story document against guidelines with fresh context. Use when a story has been created or updated and needs review for guideline compliance, task quality, and completeness.
allowed-tools: ["Skill", "Read", "Edit", "Glob", "Bash"]
hooks:
  PostToolUse:
    - matcher: "Edit"
      hooks:
        - type: command
          command: "${CLAUDE_PLUGIN_ROOT}/scripts/hook-validate-doc.sh"
---

Use the Skill tool to invoke `blueprint:story-review` with the arguments: $ARGUMENTS

This reviews the story document for guideline compliance, task quality, and completeness.
