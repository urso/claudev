#!/bin/bash
# Archive a document by setting status to archived and moving to archive directory
# Usage: archive-doc.sh <type> <id> [--with-stories] [--stories-only] [--dir DIR]
#
# Modes:
#   archive-doc.sh design <id>                  Archive a single design
#   archive-doc.sh design <id> --with-stories   Archive design and all its stories
#   archive-doc.sh design <id> --stories-only   Archive only the stories for this design
#   archive-doc.sh story <id>                   Archive a single story

SCRIPT_DIR="$(cd "$(dirname "$(readlink -f "$0")")" && pwd)"
source "$SCRIPT_DIR/resolve-dir.sh"

TYPE="${1:?Usage: archive-doc.sh <type> <id> [--with-stories] [--stories-only] [--dir DIR]}"
ID="${2:?Usage: archive-doc.sh <type> <id> [--with-stories] [--stories-only] [--dir DIR]}"
shift 2

WITH_STORIES=false
STORIES_ONLY=false
DIR_ARGS=()
while [[ $# -gt 0 ]]; do
    case $1 in
        --with-stories) WITH_STORIES=true; shift ;;
        --stories-only) STORIES_ONLY=true; shift ;;
        --dir) DIR_ARGS=(--dir "$2"); shift 2 ;;
        *) shift ;;
    esac
done

ROOT=$(git rev-parse --show-toplevel 2>/dev/null || pwd)

archive_one() {
    local ATYPE="$1" AID="$2"
    local ADIR
    ADIR=$(resolve_dir "$ATYPE" "${DIR_ARGS[@]}") || return 1

    local AFILE
    AFILE=$("$SCRIPT_DIR/find-doc.sh" "$ATYPE" "$AID" --dir "$ADIR")
    if [ $? -ne 0 ]; then
        echo "Document not found: $ATYPE $AID" >&2
        return 1
    fi

    local DEST="$ROOT/docs/ai/archive"
    case "$ATYPE" in
        design) DEST="$DEST/designs" ;;
        story)  DEST="$DEST/stories" ;;
    esac
    mkdir -p "$DEST"

    "$SCRIPT_DIR/set-status.sh" "$ATYPE" "$AID" archived --dir "$ADIR" || return 1

    local BASENAME
    BASENAME=$(basename "$AFILE")
    mv "$AFILE" "$DEST/$BASENAME"
    echo "Archived $AFILE -> $DEST/$BASENAME"
}

# Archive stories for a design
archive_stories() {
    local DESIGN_ID="$1"
    local STORIES_DIR
    STORIES_DIR=$(resolve_dir story "${DIR_ARGS[@]}") || return 1

    local LINE
    while IFS= read -r LINE; do
        [ -z "$LINE" ] && continue
        local SID
        SID=$(echo "$LINE" | cut -d'|' -f1 | xargs)
        archive_one story "$SID"
    done < <("$SCRIPT_DIR/query-stories.sh" --design "$DESIGN_ID" --dir "$STORIES_DIR")
}

if [ "$TYPE" = "design" ] && { $WITH_STORIES || $STORIES_ONLY; }; then
    # Archive stories first
    archive_stories "$ID"

    # Archive the design itself unless stories-only
    if ! $STORIES_ONLY; then
        archive_one "$TYPE" "$ID"
    fi
else
    archive_one "$TYPE" "$ID"
fi
