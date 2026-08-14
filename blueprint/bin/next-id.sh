#!/bin/bash
# Generate next document ID by scanning frontmatter 'id:' fields
# Usage: next-id.sh <type> [--dir DIR]
# Output: next available ID, zero-padded to 4 digits

SCRIPT_DIR="$(cd "$(dirname "$(readlink -f "$0")")" && pwd)"
source "$SCRIPT_DIR/resolve-dir.sh"

TYPE="${1:?Usage: next-id.sh <type> [--dir DIR]}"
shift
DIR=$(resolve_dir "$TYPE" "$@") || exit 1

if [ ! -d "$DIR" ]; then
    echo "0001"
    exit 0
fi

shopt -s nullglob
MAX=0
for f in "$DIR"/*.md; do
    ID=$(yq --front-matter=extract '.id // ""' "$f" 2>/dev/null)
    if [ -n "$ID" ] && [ "$ID" -gt "$MAX" ] 2>/dev/null; then
        MAX=$ID
    fi
done

printf "%04d\n" $(( 10#$MAX + 1 ))
