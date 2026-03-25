#!/bin/bash
# Shared test helpers for bats tests

SCRIPTS_DIR="$(cd "$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")/.." && pwd)"

setup_fixtures() {
    FIXTURES="$(mktemp -d)"
    DESIGNS="$FIXTURES/designs"
    STORIES="$FIXTURES/stories"
    RULES="$FIXTURES/rules"
    WORKFLOWS="$FIXTURES/workflows"
    mkdir -p "$DESIGNS" "$STORIES"
}

teardown_fixtures() {
    rm -rf "$FIXTURES"
}

# Create a design doc
# Usage: create_design <id> <name> [status] [depends-on...]
create_design() {
    local id="$1" name="$2" status="${3:-draft}"
    shift; shift; shift || true
    local deps=""
    if [ $# -gt 0 ]; then
        deps="depends-on:"$'\n'
        for d in "$@"; do
            deps+="  - \"$d\""$'\n'
        done
    else
        deps="depends-on: []"
    fi
    cat > "$DESIGNS/$(printf '%04d' "$id")-${name}.md" <<EOF
---
id: "$(printf '%04d' "$id")"
title: "$name"
status: $status
created: 2026-01-01
description: "Test design: $name"
$deps
---
# $name

## Tasks
- [ ] Task one
- [ ] Task two
EOF
}

# Create a story doc
# Usage: create_story <id> <name> [status] [blocked-by...]
create_story() {
    local id="$1" name="$2" status="${3:-draft}"
    shift; shift; shift || true
    local blocked=""
    if [ $# -gt 0 ]; then
        blocked="blocked-by:"$'\n'
        for b in "$@"; do
            blocked+="  - \"$b\""$'\n'
        done
    else
        blocked="blocked-by: []"
    fi
    cat > "$STORIES/$(printf '%04d' "$id")-${name}.md" <<EOF
---
id: "$(printf '%04d' "$id")"
title: "$name"
status: $status
created: 2026-01-01
description: "Test story: $name"
$blocked
---
# $name

## Tasks
- [ ] Task one
- [ ] Task two

### Decision Log
Placeholder.

### Lessons Learned
Placeholder.
EOF
}

# Create a rule file
# Usage: create_rule <name> <applies-to-json> <tags-json> <description>
# Example: create_rule "go-style" '["go"]' '[style]' "Go style conventions"
create_rule() {
    local name="$1" applies="$2" tags="$3" desc="$4"
    mkdir -p "$FIXTURES/rules"
    # Convert name to title case for display name
    local display_name
    display_name="$(echo "$name" | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) substr($i,2)} 1') Rules"
    cat > "$FIXTURES/rules/${name}.md" <<EOF
---
name: $display_name
applies-to: $applies
tags: $tags
paths: []
description: $desc
---

# $display_name

- Example rule
EOF
}

# Create a workflow file
# Usage: create_workflow <name> <tags-json> <description>
# Example: create_workflow "design" '[design]' "Design workflow guidelines"
create_workflow() {
    local name="$1" tags="$2" desc="$3"
    mkdir -p "$FIXTURES/workflows"
    local display_name
    display_name="$(echo "$name" | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) substr($i,2)} 1') Workflow"
    cat > "$FIXTURES/workflows/${name}.md" <<EOF
---
name: $display_name
applies-to: ["*"]
tags: $tags
description: $desc
---

# $display_name

- Example guideline
EOF
}
