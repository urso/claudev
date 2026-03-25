#!/bin/bash
# Query story documents with optional filters
# Usage: query-stories.sh [--dir DIR] [--design D] [--status X] [--actionable] [--search "term"]
# Output: ID | PATH | TITLE | STATUS | BLOCKED-BY | DESCRIPTION (pipe-separated)
#
# --actionable: show stories that are not done/cancelled AND have no unresolved blockers

SCRIPT_DIR="$(cd "$(dirname "$(readlink -f "$0")")" && pwd)"
source "$SCRIPT_DIR/resolve-dir.sh"
source "$SCRIPT_DIR/effective-status.sh"

DIR_OVERRIDE=""
DESIGN_FILTER=""
STATUS_FILTER=""
ACTIONABLE=false
SEARCH_FILTER=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --dir) DIR_OVERRIDE="--dir $2"; shift 2 ;;
        --design) DESIGN_FILTER="$2"; shift 2 ;;
        --status) STATUS_FILTER="$2"; shift 2 ;;
        --actionable) ACTIONABLE=true; shift ;;
        --search) SEARCH_FILTER="$2"; shift 2 ;;
        *) shift ;;
    esac
done

DIR=$(resolve_dir story $DIR_OVERRIDE) || exit 1

if [ ! -d "$DIR" ]; then
    exit 0
fi

shopt -s nullglob

# Collect all story statuses by ID (needed for effective status and --actionable)
STATUS_MAP=""
for f in "$DIR"/*.md; do
    FM=$(yq --front-matter=extract '.' "$f" 2>/dev/null) || continue
    SID=$(echo "$FM" | yq '.id // ""')
    SST=$(echo "$FM" | yq '.status // ""')
    [ -n "$SID" ] && STATUS_MAP="${STATUS_MAP}${SID}=${SST}"$'\n'
done

for f in "$DIR"/*.md; do
    FM=$(yq --front-matter=extract '.' "$f" 2>/dev/null) || continue

    ID=$(echo "$FM" | yq '.id // ""')
    TITLE=$(echo "$FM" | yq '.title // ""')
    STATUS=$(echo "$FM" | yq '.status // ""')
    DESIGN=$(echo "$FM" | yq '.design // ""')
    BLOCKED_BY=$(echo "$FM" | yq '.blocked-by // [] | join(",")')
    DESC=$(echo "$FM" | yq '.description // ""')

    if [ -n "$DESIGN_FILTER" ]; then
        [ "$DESIGN" -eq "$DESIGN_FILTER" ] 2>/dev/null || continue
    fi
    [ -n "$STATUS_FILTER" ] && [ "$STATUS" != "$STATUS_FILTER" ] && continue

    if [ -n "$SEARCH_FILTER" ]; then
        echo "$ID $TITLE $DESC" | grep -qi "$SEARCH_FILTER" || continue
    fi

    EFF_STATUS=$(effective_status "$STATUS" "$BLOCKED_BY" "$STATUS_MAP")
    EFF_BLOCKED=$(unresolved_blockers "$BLOCKED_BY" "$STATUS_MAP")

    if $ACTIONABLE; then
        case "$EFF_STATUS" in ready|in-progress) ;; *) continue ;; esac
    fi

    echo "$ID | $(basename "$f") | $TITLE | $EFF_STATUS | $EFF_BLOCKED | $DESC"
done
