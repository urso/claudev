#!/bin/bash
# Outputs PR info (if exists) and changed files with chunking for review.
# Used by /pr-review-code skill.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

base="main"
pr_number=""
pr_title=""

# Check if PR exists for current branch
if pr_info=$(gh pr view --json number,title,baseRefName 2>/dev/null); then
  base=$(echo "$pr_info" | jq -r .baseRefName)
  pr_number=$(echo "$pr_info" | jq -r .number)
  pr_title=$(echo "$pr_info" | jq -r .title)
  echo "PR #$pr_number: $pr_title"
  echo "Base: $base"
else
  echo "No PR (comparing against $base)"
fi

# Delegate to diff-files.sh with computed base
"$SCRIPT_DIR/diff-files.sh" "$base"
