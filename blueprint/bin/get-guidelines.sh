#!/bin/bash
# Get guidelines for a document type, creating from template if missing
# Usage: get-guidelines.sh <type> [workflows-dir]
# Types: design, story
# Output: Contents of the guidelines file

# Find project root (git repo root, or fallback to current dir)
PROJECT_ROOT=$(git rev-parse --show-toplevel 2>/dev/null || pwd)

TYPE="$1"
WORKFLOWS_DIR="${2:-$PROJECT_ROOT/docs/ai/workflows}"

if [ -z "$TYPE" ]; then
    echo "Usage: get-guidelines.sh <type> [workflows-dir]" >&2
    echo "Types: design, story" >&2
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GUIDELINES_FILE="$WORKFLOWS_DIR/$TYPE.md"
TEMPLATE="$SCRIPT_DIR/../templates/$TYPE-style.md"

# Create workflows dir if needed
if [ ! -d "$WORKFLOWS_DIR" ]; then
    mkdir -p "$WORKFLOWS_DIR"
fi

# Create from template if missing
if [ ! -f "$GUIDELINES_FILE" ]; then
    if [ -f "$TEMPLATE" ]; then
        cp "$TEMPLATE" "$GUIDELINES_FILE"
    else
        echo "Error: Template not found at $TEMPLATE" >&2
        exit 1
    fi
fi

cat "$GUIDELINES_FILE"
