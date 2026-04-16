#!/usr/bin/env bash
# Emit conscise mode rules as session context

SKILL_PATH="${CLAUDE_PLUGIN_ROOT}/skills/conscise/SKILL.md"

if [[ -f "$SKILL_PATH" ]]; then
  # Strip YAML frontmatter, emit body
  sed '1,/^---$/d; 1,/^---$/d' "$SKILL_PATH"
else
  echo "CONSCISE MODE ACTIVE"
  echo ""
  echo "Cut fluff. Keep grammar. Stay accurate. Short sentences. One thought per line."
fi
