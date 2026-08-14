#!/bin/bash
# Validate a design or story document
# Usage: validate-doc.sh <type> <file> [--dir DIR]
# Exits 0 if valid, 1 if errors found. Prints all errors to stdout.

SCRIPT_DIR="$(cd "$(dirname "$(readlink -f "$0")")" && pwd)"
source "$SCRIPT_DIR/resolve-dir.sh"

TYPE="${1:?Usage: validate-doc.sh <type> <file> [--dir DIR]}"
FILE="${2:?Usage: validate-doc.sh <type> <file> [--dir DIR]}"
shift 2

DIR=$(resolve_dir "$TYPE" "$@") || exit 1

ERRORS=()
err() { ERRORS+=("$1"); }

if [ ! -f "$FILE" ]; then
    echo "ERROR: File not found: $FILE"
    exit 1
fi

FM=$(yq --front-matter=extract '.' "$FILE" 2>/dev/null)
if [ -z "$FM" ]; then
    echo "ERROR: No frontmatter found in $FILE"
    exit 1
fi

fm_val() { echo "$FM" | yq "$1 // \"\""; }

# --- Frontmatter validation ---

VALID_STATUSES="draft ready in-progress done cancelled archived"

ID=$(fm_val '.id')
TITLE=$(fm_val '.title')
STATUS=$(fm_val '.status')
CREATED=$(fm_val '.created')
DESC=$(fm_val '.description')

[ -z "$ID" ] && err "Missing required field: id"
[ -z "$TITLE" ] && err "Missing required field: title"
[ -z "$STATUS" ] && err "Missing required field: status"
[ -z "$CREATED" ] && err "Missing required field: created"
[ -z "$DESC" ] && err "Missing required field: description"

if [ -n "$STATUS" ]; then
    FOUND=false
    for S in $VALID_STATUSES; do
        [ "$STATUS" = "$S" ] && FOUND=true
    done
    $FOUND || err "Invalid status: '$STATUS' (expected: $VALID_STATUSES)"
fi

# Check for duplicate IDs
if [ -n "$ID" ] && [ -d "$DIR" ]; then
    for f in "$DIR"/*.md; do
        [ "$f" = "$FILE" ] && continue
        [ ! -f "$f" ] && continue
        OTHER_ID=$(yq --front-matter=extract '.id // ""' "$f" 2>/dev/null)
        if [ "$OTHER_ID" -eq "$ID" ] 2>/dev/null; then
            err "Duplicate id '$ID' — also found in $(basename "$f")"
            break
        fi
    done
fi

# Type-specific frontmatter
if [ "$TYPE" = "story" ]; then
    # blocked-by references must exist
    BLOCKED=$(echo "$FM" | yq '.blocked-by // [] | .[]' 2>/dev/null)
    for BID in $BLOCKED; do
        if ! "$SCRIPT_DIR/find-doc.sh" "$TYPE" "$BID" --dir "$DIR" >/dev/null 2>&1; then
            err "blocked-by references non-existent $TYPE: $BID"
        fi
    done

    # Cycle detection: walk blocked-by chains from this doc's ID
    if [ -n "$ID" ] && [ -n "$BLOCKED" ]; then
        VISITED="$ID"
        QUEUE="$BLOCKED"
        CYCLE_FOUND=false
        while [ -n "$QUEUE" ] && ! $CYCLE_FOUND; do
            NEXT_QUEUE=""
            for QID in $QUEUE; do
                if [ "$QID" -eq "$ID" ] 2>/dev/null; then
                    err "Cycle detected in blocked-by chain involving: $ID"
                    CYCLE_FOUND=true
                    break
                fi
                # Skip if already visited
                echo "$VISITED" | grep -qw "$QID" && continue
                VISITED="$VISITED $QID"
                # Get blocked-by for this node
                QFILE=$("$SCRIPT_DIR/find-doc.sh" "$TYPE" "$QID" --dir "$DIR" 2>/dev/null) || continue
                QBLOCKED=$(yq --front-matter=extract '.blocked-by // [] | .[]' "$QFILE" 2>/dev/null)
                NEXT_QUEUE="$NEXT_QUEUE $QBLOCKED"
            done
            QUEUE="$NEXT_QUEUE"
        done
    fi
elif [ "$TYPE" = "design" ]; then
    # depends-on references must exist
    DEPENDS=$(echo "$FM" | yq '.depends-on // [] | .[]' 2>/dev/null)
    for DID in $DEPENDS; do
        if ! "$SCRIPT_DIR/find-doc.sh" "$TYPE" "$DID" --dir "$DIR" >/dev/null 2>&1; then
            err "depends-on references non-existent $TYPE: $DID"
        fi
    done

    # references must point to existing files (relative to git root)
    FILE_DIR=$(cd "$(dirname "$FILE")" && pwd)
    GIT_ROOT=$(git -C "$FILE_DIR" rev-parse --show-toplevel 2>/dev/null)
    if [ -n "$GIT_ROOT" ]; then
        REFS=$(echo "$FM" | yq '.references // [] | .[]' 2>/dev/null)
        for REF in $REFS; do
            if [ ! -e "$GIT_ROOT/$REF" ]; then
                err "references non-existent file: $REF"
            fi
        done
    fi
fi

# --- Content validation ---

BODY=$(awk '/^---$/{if(n++)found=1;next} found{print}' "$FILE")

# Collect all ## headings
SECTIONS=$(echo "$BODY" | grep -E '^## ' | sed 's/^## //')

# Check for empty sections (heading followed by another heading or EOF with only whitespace between)
PREV_LINE=""
PREV_IS_HEADING=false
SECTION_HAS_CONTENT=false
CURRENT_SECTION=""
while IFS= read -r line; do
    if echo "$line" | grep -qE '^## '; then
        if $PREV_IS_HEADING || { [ -n "$CURRENT_SECTION" ] && ! $SECTION_HAS_CONTENT; }; then
            err "Empty section: ## $CURRENT_SECTION"
        fi
        CURRENT_SECTION=$(echo "$line" | sed 's/^## //')
        PREV_IS_HEADING=true
        SECTION_HAS_CONTENT=false
    elif echo "$line" | grep -qE '^[[:space:]]*$'; then
        : # skip blank lines
    else
        PREV_IS_HEADING=false
        SECTION_HAS_CONTENT=true
    fi
done <<< "$BODY"
# Check last section
if [ -n "$CURRENT_SECTION" ] && ! $SECTION_HAS_CONTENT; then
    err "Empty section: ## $CURRENT_SECTION"
fi

# Check checkboxes only appear in Tasks or Phase sections
CURRENT_SECTION=""
IN_TASK_SECTION=false
while IFS= read -r line; do
    if echo "$line" | grep -qE '^## |^### '; then
        CURRENT_SECTION=$(echo "$line" | sed 's/^#* //')
        if echo "$CURRENT_SECTION" | grep -qiE '^(Tasks|Task [0-9]|Phase [0-9])'; then
            IN_TASK_SECTION=true
        else
            IN_TASK_SECTION=false
        fi
    fi
    if echo "$line" | grep -qE '^\s*- \[(x| )\]' && ! $IN_TASK_SECTION; then
        err "Checkbox found outside Tasks/Phase section, in: $CURRENT_SECTION"
    fi
done <<< "$BODY"

# Require at least one section with checkboxes
HAS_CHECKBOXES=false
CURRENT_SECTION=""
while IFS= read -r line; do
    if echo "$line" | grep -qE '^## |^### '; then
        CURRENT_SECTION=$(echo "$line" | sed 's/^#* //')
    fi
    if echo "$line" | grep -qE '^\s*- \[(x| )\]'; then
        if echo "$CURRENT_SECTION" | grep -qiE '^(Tasks|Task [0-9]|Phase [0-9])'; then
            HAS_CHECKBOXES=true
        fi
    fi
done <<< "$BODY"
if [ "$TYPE" = "story" ]; then
    $HAS_CHECKBOXES || err "No checkboxes found in Tasks or Phase sections"
fi

# For done stories: developer log sections must not be empty
if [ "$TYPE" = "story" ] && [ "$STATUS" = "done" ]; then
    for LOG_SECTION in "Decision Log" "Lessons Learned"; do
        SECTION_CONTENT=""
        IN_SECTION=false
        while IFS= read -r line; do
            if echo "$line" | grep -qE "^### $LOG_SECTION"; then
                IN_SECTION=true
                continue
            fi
            if $IN_SECTION; then
                if echo "$line" | grep -qE '^##'; then
                    break
                fi
                TRIMMED=$(echo "$line" | xargs)
                if [ -n "$TRIMMED" ]; then
                    SECTION_CONTENT="yes"
                fi
            fi
        done <<< "$BODY"
        [ -z "$SECTION_CONTENT" ] && err "Developer log section empty for done story: $LOG_SECTION"
    done
fi

# --- Output ---

if [ ${#ERRORS[@]} -gt 0 ]; then
    for e in "${ERRORS[@]}"; do
        echo "ERROR: $e"
    done
    exit 1
fi

echo "OK: $FILE"
exit 0
