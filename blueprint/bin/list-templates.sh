#!/bin/bash
# List available design templates with metadata
# Usage: list-templates.sh [type]
# Examples:
#   list-templates.sh           # List design templates (default)
#   list-templates.sh design    # Explicit design templates
#   list-templates.sh story     # Story templates (future)
#
# Output format:
#   name|path|description
#
# Searches:
#   1. User templates: docs/ai/templates/<type>s/*/template.md
#   2. Default template: $CLAUDE_PLUGIN_ROOT/templates/<type>s/default/template.md

TYPE="${1:-design}"
SUBDIR="${TYPE}s"

PROJECT_ROOT=$(git rev-parse --show-toplevel 2>/dev/null || pwd)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PLUGIN_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

source "$SCRIPT_DIR/parse-frontmatter.sh"

list_template() {
    local f="$1"
    local name="$2"
    [ -f "$f" ] || return

    JSON=$(parse_frontmatter "$f" 2>/dev/null)
    DESC=""
    if [ -n "$JSON" ]; then
        DESC=$(echo "$JSON" | yq -r '.description // ""' | tr '\n' ' ' | sed 's/  */ /g' | sed 's/ *$//')
    fi

    echo "$name|$f|$DESC"
}

# User templates
USER_DIR="$PROJECT_ROOT/docs/ai/templates/$SUBDIR"
if [ -d "$USER_DIR" ]; then
    for d in "$USER_DIR"/*/; do
        [ -d "$d" ] || continue
        NAME=$(basename "$d")
        TEMPLATE="$d/template.md"
        list_template "$TEMPLATE" "$NAME"
    done
fi

# Default template (bundled)
DEFAULT="$PLUGIN_ROOT/templates/$SUBDIR/default/template.md"
list_template "$DEFAULT" "default"
