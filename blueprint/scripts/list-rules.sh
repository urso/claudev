#!/bin/bash
# List rules and their metadata from frontmatter
# Usage: list-rules.sh [rules-dir] [applies-to-filter] [tag-filter]
# Examples:
#   list-rules.sh                           # List all rules
#   list-rules.sh docs/ai/rules go          # Rules that apply to 'go'
#   list-rules.sh docs/ai/rules go style    # Rules for 'go' with tag 'style'
#   list-rules.sh docs/ai/rules "" build    # All rules with tag 'build'
#
# Output format (grouped by directory):
#   <directory>:
#   full/path/to/file.md|name|applies-to|tags|paths|description

# Find project root (git repo root, or fallback to current dir)
PROJECT_ROOT=$(git rev-parse --show-toplevel 2>/dev/null || pwd)

DIR="${1:-$PROJECT_ROOT/docs/ai/rules}"
APPLIES_FILTER="$2"
TAG_FILTER="$3"

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
    # Handle both array and scalar values, join arrays with ", "
    APPLIES=$(echo "$JSON" | yq -r '."applies-to" // "" | @json' | tr -d '"[]' | tr ',' ', ' | sed 's/  /, /g')
    TAGS=$(echo "$JSON" | yq -r '.tags // "" | @json' | tr -d '"[]' | tr ',' ', ' | sed 's/  /, /g')
    PATHS=$(echo "$JSON" | yq -r '.paths // "" | @json' | tr -d '"[]' | tr ',' ', ' | sed 's/  /, /g')
    DESC=$(echo "$JSON" | yq -r '.description // ""' | tr '\n' ' ' | sed 's/  */ /g' | sed 's/ *$//')

    # Filter by applies-to if specified
    if [ -n "$APPLIES_FILTER" ]; then
        if [[ "$APPLIES" != *"*"* ]] && [[ "$APPLIES" != *"$APPLIES_FILTER"* ]]; then
            continue
        fi
    fi

    # Filter by tag if specified
    if [ -n "$TAG_FILTER" ]; then
        if [[ "$TAGS" != *"$TAG_FILTER"* ]]; then
            continue
        fi
    fi

    echo "$f|$NAME|$APPLIES|$TAGS|$PATHS|$DESC"
done
