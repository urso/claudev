#!/bin/bash
# Query design documents with optional filters
# Usage: query-designs.sh [--dir DIR] [--status X] [--search "term"]
# Output: ID | PATH | TITLE | STATUS | DESCRIPTION (pipe-separated)

SCRIPT_DIR="$(cd "$(dirname "$(readlink -f "$0")")" && pwd)"
source "$SCRIPT_DIR/resolve-dir.sh"

DIR_OVERRIDE=""
STATUS_FILTER=""
SEARCH_FILTER=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --dir) DIR_OVERRIDE="--dir $2"; shift 2 ;;
        --status) STATUS_FILTER="$2"; shift 2 ;;
        --search) SEARCH_FILTER="$2"; shift 2 ;;
        *) shift ;;
    esac
done

DIR=$(resolve_dir design $DIR_OVERRIDE) || exit 1

if [ ! -d "$DIR" ]; then
    exit 0
fi

shopt -s nullglob
for f in "$DIR"/*.md; do
    FM=$(yq --front-matter=extract '.' "$f" 2>/dev/null) || continue

    ID=$(echo "$FM" | yq '.id // ""')
    TITLE=$(echo "$FM" | yq '.title // ""')
    STATUS=$(echo "$FM" | yq '.status // ""')
    DESC=$(echo "$FM" | yq '.description // ""')

    [ -n "$STATUS_FILTER" ] && [ "$STATUS" != "$STATUS_FILTER" ] && continue

    if [ -n "$SEARCH_FILTER" ]; then
        echo "$ID $TITLE $DESC" | grep -qi "$SEARCH_FILTER" || continue
    fi

    echo "$ID | $(basename "$f") | $TITLE | $STATUS | $DESC"
done
