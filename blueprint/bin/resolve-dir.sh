#!/bin/bash
# Resolve document directory from type name
# Usage: resolve-dir.sh <type> [--dir override]
#   or:  source resolve-dir.sh; resolve_dir <type> [--dir override]
# Types: design, story, rules

resolve_dir() {
    local TYPE="$1"
    local DIR_OVERRIDE=""
    shift
    while [[ $# -gt 0 ]]; do
        case $1 in
            --dir) DIR_OVERRIDE="$2"; shift 2 ;;
            *) break ;;
        esac
    done

    if [ -n "$DIR_OVERRIDE" ]; then
        echo "$DIR_OVERRIDE"
        return
    fi

    local ROOT
    ROOT=$(git rev-parse --show-toplevel 2>/dev/null || pwd)

    case "$TYPE" in
        design) echo "$ROOT/docs/ai/designs" ;;
        story)  echo "$ROOT/docs/ai/stories" ;;
        rules)  echo "$ROOT/docs/ai/rules" ;;
        workflows) echo "$ROOT/docs/ai/workflows" ;;
        *)
            echo "Unknown type: $TYPE (expected: design, story, rules)" >&2
            return 1
            ;;
    esac
}

# When executed directly (not sourced), run with args
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    resolve_dir "$@"
fi
