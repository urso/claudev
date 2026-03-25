#!/bin/bash
# Compute effective status for a story given its raw status and blockers
# Usage: source effective-status.sh; effective_status <raw-status> <blocked-by-csv> <status-map>
#
# status-map is a newline-separated string of "ID=STATUS" entries
# Returns: effective status on stdout
#   - "blocked" if ready/in-progress with unresolved blockers
#   - raw status otherwise
#
# A blocker is resolved if its status is done, cancelled, or archived.

effective_status() {
    local RAW_STATUS="$1"
    local BLOCKED_BY="$2"
    local STATUS_MAP="$3"

    case "$RAW_STATUS" in
        ready|in-progress)
            if [ -n "$BLOCKED_BY" ]; then
                IFS=',' read -ra BLOCKERS <<< "$BLOCKED_BY"
                for BID in "${BLOCKERS[@]}"; do
                    BID=$(echo "$BID" | xargs)
                    BSTATUS=""
                    while IFS='=' read -r MID MSTATUS; do
                        [ -z "$MID" ] && continue
                        if [ "$MID" -eq "$BID" ] 2>/dev/null; then
                            BSTATUS="$MSTATUS"
                            break
                        fi
                    done <<< "$STATUS_MAP"
                    case "$BSTATUS" in
                        done|cancelled|archived) ;; # resolved
                        *) echo "blocked"; return ;;
                    esac
                done
            fi
            ;;
    esac
    echo "$RAW_STATUS"
}

# Return only unresolved blocker IDs (comma-separated), empty if all resolved
unresolved_blockers() {
    local BLOCKED_BY="$1"
    local STATUS_MAP="$2"
    local RESULT=""

    if [ -n "$BLOCKED_BY" ]; then
        IFS=',' read -ra BLOCKERS <<< "$BLOCKED_BY"
        for BID in "${BLOCKERS[@]}"; do
            BID=$(echo "$BID" | xargs)
            BSTATUS=""
            while IFS='=' read -r MID MSTATUS; do
                [ -z "$MID" ] && continue
                if [ "$MID" -eq "$BID" ] 2>/dev/null; then
                    BSTATUS="$MSTATUS"
                    break
                fi
            done <<< "$STATUS_MAP"
            case "$BSTATUS" in
                done|cancelled|archived) ;; # resolved
                *) RESULT="${RESULT:+$RESULT,}$BID" ;;
            esac
        done
    fi
    echo "$RESULT"
}
