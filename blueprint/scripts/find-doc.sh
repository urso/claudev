#!/bin/bash
# Find document file path by ID
# Usage: find-doc.sh <type> <id> [--dir DIR]
# Output: absolute file path, or exits 1 if not found

SCRIPT_DIR="$(cd "$(dirname "$(readlink -f "$0")")" && pwd)"
source "$SCRIPT_DIR/resolve-dir.sh"

TYPE="${1:?Usage: find-doc.sh <type> <id> [--dir DIR]}"
ID="${2:?Usage: find-doc.sh <type> <id> [--dir DIR]}"
shift 2
DIR=$(resolve_dir "$TYPE" "$@") || exit 1

if [ ! -d "$DIR" ]; then
    echo "Directory not found: $DIR" >&2
    exit 1
fi

shopt -s nullglob
for f in "$DIR"/*.md; do
    FID=$(yq --front-matter=extract '.id // ""' "$f" 2>/dev/null)
    if [ "$FID" -eq "$ID" ] 2>/dev/null; then
        echo "$f"
        exit 0
    fi
done

echo "No document found with id: $ID" >&2
exit 1
