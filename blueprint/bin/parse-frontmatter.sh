#!/bin/bash
# Shared frontmatter parsing function
# Usage: source parse-frontmatter.sh; parse_frontmatter <file>
# Returns: JSON representation of YAML frontmatter via stdout

parse_frontmatter() {
    local file="$1"
    # Extract content between first two --- lines only (stop after second ---)
    local raw
    raw=$(awk '
        /^---$/ {
            count++
            if (count == 2) exit
            next
        }
        count == 1 { print }
    ' "$file")
    # Return empty if no frontmatter content found
    [ -z "$raw" ] && return 1
    echo "$raw" | yq -o=json
}
