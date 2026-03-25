#!/bin/bash
# Query designs and stories together with effective status
# Usage: query-work.sh [--actionable] [--search "term"] [--design-dir DIR] [--story-dir DIR]
# Output: TYPE | ID | FILENAME | TITLE | STATUS | BLOCKED-BY | DESCRIPTION (pipe-separated)
#
# STATUS is the effective status:
#   - "blocked" if story is ready/in-progress but has unresolved blockers
#   - otherwise the raw status from frontmatter
#
# --actionable: only show designs with actionable stories and the actionable stories themselves
#   A story is actionable if it is ready/in-progress with all blockers resolved (done).
#   A design is actionable if it has at least one actionable story.

SCRIPT_DIR="$(cd "$(dirname "$(readlink -f "$0")")" && pwd)"
source "$SCRIPT_DIR/resolve-dir.sh"
source "$SCRIPT_DIR/effective-status.sh"

DESIGN_DIR_OVERRIDE=""
STORY_DIR_OVERRIDE=""
ACTIONABLE=false
SEARCH_FILTER=""

while [[ $# -gt 0 ]]; do
    case $1 in
        --design-dir) DESIGN_DIR_OVERRIDE="--dir $2"; shift 2 ;;
        --story-dir) STORY_DIR_OVERRIDE="--dir $2"; shift 2 ;;
        --actionable) ACTIONABLE=true; shift ;;
        --search) SEARCH_FILTER="$2"; shift 2 ;;
        *) shift ;;
    esac
done

DESIGN_DIR=$(resolve_dir design $DESIGN_DIR_OVERRIDE) || exit 1
STORY_DIR=$(resolve_dir story $STORY_DIR_OVERRIDE) || exit 1

shopt -s nullglob

# --- Collect all story data and status map ---

STATUS_MAP=""
declare -a STORY_LINES

for f in "$STORY_DIR"/*.md; do
    FM=$(yq --front-matter=extract '.' "$f" 2>/dev/null) || continue

    ID=$(echo "$FM" | yq '.id // ""')
    TITLE=$(echo "$FM" | yq '.title // ""')
    STATUS=$(echo "$FM" | yq '.status // ""')
    DESIGN=$(echo "$FM" | yq '.design // ""')
    BLOCKED_BY=$(echo "$FM" | yq '.blocked-by // [] | join(",")')
    DESC=$(echo "$FM" | yq '.description // ""')

    [ -n "$ID" ] && STATUS_MAP="${STATUS_MAP}${ID}=${STATUS}"$'\n'

    STORY_LINES+=("$ID|$(basename "$f")|$TITLE|$STATUS|$DESIGN|$BLOCKED_BY|$DESC")
done

# --- Track which designs have actionable stories ---

declare -A DESIGN_HAS_ACTIONABLE

for line in "${STORY_LINES[@]}"; do
    IFS='|' read -r ID FNAME TITLE STATUS DESIGN BLOCKED_BY DESC <<< "$line"
    EFF=$(effective_status "$STATUS" "$BLOCKED_BY" "$STATUS_MAP")
    if [ -n "$DESIGN" ] && [[ "$EFF" == "ready" || "$EFF" == "in-progress" ]]; then
        DESIGN_HAS_ACTIONABLE["$DESIGN"]=true
    fi
done

# --- Output designs and track matched IDs ---

declare -A MATCHED_DESIGNS

for f in "$DESIGN_DIR"/*.md; do
    FM=$(yq --front-matter=extract '.' "$f" 2>/dev/null) || continue

    DID=$(echo "$FM" | yq '.id // ""')
    DTITLE=$(echo "$FM" | yq '.title // ""')
    DSTATUS=$(echo "$FM" | yq '.status // ""')
    DDESC=$(echo "$FM" | yq '.description // ""')

    if [ -n "$SEARCH_FILTER" ]; then
        echo "$DID $DTITLE $DDESC" | grep -qi "$SEARCH_FILTER" || continue
    fi

    if $ACTIONABLE; then
        case "$DSTATUS" in done|cancelled|archived) continue ;; esac
        [ "${DESIGN_HAS_ACTIONABLE[$DID]:-}" = "true" ] || continue
    fi

    MATCHED_DESIGNS["$DID"]=true
    echo "design | $DID | $(basename "$f") | $DTITLE | $DSTATUS | | $DDESC"
done

# --- Output stories ---

for line in "${STORY_LINES[@]}"; do
    IFS='|' read -r ID FNAME TITLE STATUS DESIGN BLOCKED_BY DESC <<< "$line"

    # Match if story matches search OR its parent design was matched
    if [ -n "$SEARCH_FILTER" ]; then
        STORY_MATCHES=false
        echo "$ID $TITLE $DESC" | grep -qi "$SEARCH_FILTER" && STORY_MATCHES=true
        [ -n "$DESIGN" ] && [ "${MATCHED_DESIGNS[$DESIGN]:-}" = "true" ] && STORY_MATCHES=true
        $STORY_MATCHES || continue
    fi

    EFF_STATUS=$(effective_status "$STATUS" "$BLOCKED_BY" "$STATUS_MAP")
    EFF_BLOCKED=$(unresolved_blockers "$BLOCKED_BY" "$STATUS_MAP")

    if $ACTIONABLE; then
        case "$EFF_STATUS" in ready|in-progress) ;; *) continue ;; esac
    fi

    echo "story | $ID | $FNAME | $TITLE | $EFF_STATUS | $EFF_BLOCKED | $DESC"
done
