#!/bin/bash
# List workflow guides and their metadata from frontmatter
# Usage: list-workflows.sh [workflows-dir] [type-filter]
# Examples:
#   list-workflows.sh                              # List all workflows
#   list-workflows.sh docs/ai/workflows design     # Only design workflows
#   list-workflows.sh docs/ai/workflows story      # Only story workflows
#
# Output format (grouped by directory):
#   <directory>:
#   full/path/to/file.md|name|applies-to|tags|description

# Find project root (git repo root, or fallback to current dir)
PROJECT_ROOT=$(git rev-parse --show-toplevel 2>/dev/null || pwd)

DIR="${1:-$PROJECT_ROOT/docs/ai/workflows}"
TYPE_FILTER="$2"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/parse-frontmatter.sh"

if [ ! -d "$DIR" ]; then
    exit 0
fi

echo "$DIR:"

shopt -s nullglob
for f in "$DIR"/*.md; do
    [ -f "$f" ] || continue

    # Parse frontmatter to JSON
    JSON=$(parse_frontmatter "$f" 2>/dev/null)
    [ -z "$JSON" ] && continue

    NAME=$(echo "$JSON" | yq -r '.name // ""')
    APPLIES=$(echo "$JSON" | yq -r '."applies-to" // "" | @json' | tr -d '"[]' | tr ',' ', ' | sed 's/  /, /g')
    TAGS=$(echo "$JSON" | yq -r '.tags // "" | @json' | tr -d '"[]' | tr ',' ', ' | sed 's/  /, /g')
    DESC=$(echo "$JSON" | yq -r '.description // ""' | tr '\n' ' ' | sed 's/  */ /g' | sed 's/ *$//')

    # Filter by type (tag) if specified
    if [ -n "$TYPE_FILTER" ]; then
        if [[ "$TAGS" != *"$TYPE_FILTER"* ]]; then
            continue
        fi
    fi

    echo "$f|$NAME|$APPLIES|$TAGS|$DESC"
done
