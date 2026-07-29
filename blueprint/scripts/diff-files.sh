#!/bin/bash
# Output changed files with diff line counts, grouped into ~500 line chunks.
# Usage: diff-files.sh [base]
#   base: branch/commit to compare against (e.g., main, origin/main)
#         if omitted, shows staged + unstaged changes (git diff HEAD)

set -euo pipefail

BASE="${1:-}"
CHUNK_BUDGET=500

# Get changed files with diff stats (additions + deletions)
if [[ -n "$BASE" ]]; then
    numstat=$(git diff "$BASE"...HEAD --numstat 2>/dev/null || true)
else
    # Staged + unstaged changes
    numstat=$(git diff HEAD --numstat 2>/dev/null || true)
fi

if [[ -z "$numstat" ]]; then
    echo "No changed files"
    exit 0
fi

# Build file list with diff line counts, language, directory
# Format: lines|language|directory|filepath
declare -a file_data=()
total_lines=0

while IFS=$'\t' read -r added deleted filepath; do
    [[ -z "$filepath" ]] && continue

    # Handle binary files (shown as "-" in numstat)
    if [[ "$added" == "-" ]]; then
        added=0
        deleted=0
    fi

    # Diff size = additions + deletions
    lines=$((added + deleted))
    total_lines=$((total_lines + lines))

    # Extract extension for language grouping
    ext="${filepath##*.}"
    if [[ "$filepath" == "$ext" ]]; then
        ext="none"
    fi

    # Extract directory
    dir=$(dirname "$filepath")

    file_data+=("$lines|$ext|$dir|$filepath")
done <<< "$numstat"

if [[ ${#file_data[@]} -eq 0 ]]; then
    echo "No changed files"
    exit 0
fi

# Sort by language, then directory (for grouping)
IFS=$'\n' sorted=($(sort -t'|' -k2,2 -k3,3 <<< "${file_data[*]}")); unset IFS

# Pack into chunks
declare -a chunks=()
current_chunk=""
current_lines=0
current_label=""

for entry in "${sorted[@]}"; do
    IFS='|' read -r lines ext dir filepath <<< "$entry"

    label="$ext:$dir"

    # Start new chunk if:
    # 1. Adding this file would exceed budget AND chunk is not empty
    # 2. OR this single file exceeds budget (gets its own chunk)
    if [[ $current_lines -gt 0 && $((current_lines + lines)) -gt $CHUNK_BUDGET ]]; then
        chunks+=("$current_chunk")
        current_chunk=""
        current_lines=0
        current_label=""
    fi

    # Add file to current chunk
    if [[ -n "$current_chunk" ]]; then
        current_chunk+=$'\n'
    fi
    current_chunk+="  $filepath	$lines"
    current_lines=$((current_lines + lines))

    # Track label for first file in chunk
    if [[ -z "$current_label" ]]; then
        current_label="$label"
    fi
done

# Don't forget last chunk
if [[ -n "$current_chunk" ]]; then
    chunks+=("$current_chunk")
fi

# Output
echo "Files: ${#file_data[@]} ($total_lines lines total)"
echo "---"

chunk_num=1
for chunk in "${chunks[@]}"; do
    # Calculate chunk total
    chunk_lines=0
    while IFS= read -r line; do
        [[ -z "$line" ]] && continue
        l=$(echo "$line" | awk '{print $NF}')
        chunk_lines=$((chunk_lines + l))
    done <<< "$chunk"

    echo "Chunk $chunk_num ($chunk_lines lines):"
    echo "$chunk"
    echo ""
    chunk_num=$((chunk_num + 1))
done
