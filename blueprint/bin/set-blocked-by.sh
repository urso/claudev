#!/bin/bash
# Append to blocked-by list for a story by ID (deduplicates)
# Usage: set-blocked-by.sh <type> <id> <blocker-id...> [--dir DIR]

SCRIPT_DIR="$(cd "$(dirname "$(readlink -f "$0")")" && pwd)"
source "$SCRIPT_DIR/resolve-dir.sh"

TYPE="${1:?Usage: set-blocked-by.sh <type> <id> <blocker-id...> [--dir DIR]}"
ID="${2:?Usage: set-blocked-by.sh <type> <id> <blocker-id...> [--dir DIR]}"
shift 2

BLOCKERS=()
DIR_ARGS=()
while [[ $# -gt 0 ]]; do
    case $1 in
        --dir) DIR_ARGS=(--dir "$2"); shift 2 ;;
        *) BLOCKERS+=("$1"); shift ;;
    esac
done

if [ ${#BLOCKERS[@]} -eq 0 ]; then
    echo "Usage: set-blocked-by.sh <type> <id> <blocker-id...> [--dir DIR]" >&2
    exit 1
fi

FILE=$("$SCRIPT_DIR/find-doc.sh" "$TYPE" "$ID" "${DIR_ARGS[@]}")
if [ $? -ne 0 ]; then
    echo "Document not found: $TYPE $ID" >&2
    exit 1
fi

# Append each blocker, deduplicating
for B in "${BLOCKERS[@]}"; do
    yq --front-matter=process --inplace \
        ".\"blocked-by\" = ((.\"blocked-by\" // []) + [\"$B\"] | unique)" "$FILE"
done

# Read back the final list for confirmation
FINAL=$(yq --front-matter=extract '.blocked-by // [] | join(",")' "$FILE")
echo "Updated $FILE: blocked-by=[$FINAL]"
