#!/bin/bash
# Remove entries from blocked-by list for a story by ID
# Usage: remove-blocked-by.sh <type> <id> <blocker-id...> [--dir DIR]

SCRIPT_DIR="$(cd "$(dirname "$(readlink -f "$0")")" && pwd)"
source "$SCRIPT_DIR/resolve-dir.sh"

TYPE="${1:?Usage: remove-blocked-by.sh <type> <id> <blocker-id...> [--dir DIR]}"
ID="${2:?Usage: remove-blocked-by.sh <type> <id> <blocker-id...> [--dir DIR]}"
shift 2

TO_REMOVE=()
DIR_ARGS=()
while [[ $# -gt 0 ]]; do
    case $1 in
        --dir) DIR_ARGS=(--dir "$2"); shift 2 ;;
        *) TO_REMOVE+=("$1"); shift ;;
    esac
done

if [ ${#TO_REMOVE[@]} -eq 0 ]; then
    echo "Usage: remove-blocked-by.sh <type> <id> <blocker-id...> [--dir DIR]" >&2
    exit 1
fi

FILE=$("$SCRIPT_DIR/find-doc.sh" "$TYPE" "$ID" "${DIR_ARGS[@]}")
if [ $? -ne 0 ]; then
    echo "Document not found: $TYPE $ID" >&2
    exit 1
fi

# Build yq expression to remove each blocker
YQ_EXPR=".\"blocked-by\" |= (. // [])"
for BID in "${TO_REMOVE[@]}"; do
    YQ_EXPR+=" | .\"blocked-by\" -= [\"$BID\"]"
done

yq --front-matter=process --inplace "$YQ_EXPR" "$FILE"
echo "Updated $FILE: removed [${TO_REMOVE[*]}] from blocked-by"
