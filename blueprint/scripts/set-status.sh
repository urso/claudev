#!/bin/bash
# Set status of a document by ID
# Usage: set-status.sh <type> <id> <status> [--dir DIR]

SCRIPT_DIR="$(cd "$(dirname "$(readlink -f "$0")")" && pwd)"
source "$SCRIPT_DIR/resolve-dir.sh"

TYPE="${1:?Usage: set-status.sh <type> <id> <status> [--dir DIR]}"
ID="${2:?Usage: set-status.sh <type> <id> <status> [--dir DIR]}"
STATUS="${3:?Usage: set-status.sh <type> <id> <status> [--dir DIR]}"
shift 3

FILE=$("$SCRIPT_DIR/find-doc.sh" "$TYPE" "$ID" "$@")
if [ $? -ne 0 ]; then
    echo "Document not found: $TYPE $ID" >&2
    exit 1
fi

yq --front-matter=process --inplace ".status = \"$STATUS\"" "$FILE"

# Record commit SHA when transitioning to in-progress
if [ "$STATUS" = "in-progress" ]; then
    CURRENT_SHA=$(git rev-parse HEAD 2>/dev/null || echo "")
    if [ -n "$CURRENT_SHA" ]; then
        yq --front-matter=process --inplace ".\"start-git-sha\" = \"$CURRENT_SHA\"" "$FILE"
        echo "Updated $FILE: status=$STATUS, start-git-sha=$CURRENT_SHA"
        exit 0
    fi
fi

echo "Updated $FILE: status=$STATUS"
