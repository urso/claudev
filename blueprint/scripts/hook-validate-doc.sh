#!/bin/bash
# PostToolUse hook: validate design/story documents after Write or Edit
# Reads hook JSON from stdin, extracts file_path, infers type from path
# Usage: hook-validate-doc.sh

SCRIPT_DIR="$(cd "$(dirname "$(readlink -f "$0")")" && pwd)"

INPUT=$(cat)
FILE_PATH=$(echo "$INPUT" | jq -r '.tool_input.file_path // empty')

[ -z "$FILE_PATH" ] && exit 0
[ -f "$FILE_PATH" ] || exit 0

# Infer type from path
case "$FILE_PATH" in
    */designs/*) TYPE="design" ;;
    */stories/*) TYPE="story" ;;
    *)           exit 0 ;;
esac

"$SCRIPT_DIR/validate-doc.sh" "$TYPE" "$FILE_PATH"
