#!/usr/bin/env bats

load test_helper

setup() { setup_fixtures; }
teardown() { teardown_fixtures; }

@test "validate-doc: valid story passes" {
    create_story 1 "good" "ready"
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/0001-good.md" --dir "$STORIES"
    [ "$status" -eq 0 ]
    [[ "$output" == "OK:"* ]]
}

@test "validate-doc: valid design passes" {
    create_design 1 "good" "draft"
    run "$SCRIPTS_DIR/validate-doc.sh" design "$DESIGNS/0001-good.md" --dir "$DESIGNS"
    [ "$status" -eq 0 ]
    [[ "$output" == "OK:"* ]]
}

@test "validate-doc: missing id field" {
    cat > "$STORIES/bad.md" <<'EOF'
---
title: "No ID"
status: draft
created: 2026-01-01
description: "Missing id"
---
# Bad

## Tasks
- [ ] Something
EOF
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/bad.md" --dir "$STORIES"
    [ "$status" -eq 1 ]
    [[ "$output" == *"Missing required field: id"* ]]
}

@test "validate-doc: missing multiple fields" {
    cat > "$STORIES/bad.md" <<'EOF'
---
id: "0099"
---
# Bad

## Tasks
- [ ] Something
EOF
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/bad.md" --dir "$STORIES"
    [ "$status" -eq 1 ]
    [[ "$output" == *"Missing required field: title"* ]]
    [[ "$output" == *"Missing required field: status"* ]]
    [[ "$output" == *"Missing required field: created"* ]]
    [[ "$output" == *"Missing required field: description"* ]]
}

@test "validate-doc: invalid status" {
    cat > "$STORIES/bad.md" <<'EOF'
---
id: "0099"
title: "Bad"
status: banana
created: 2026-01-01
description: "Invalid status"
---
# Bad

## Tasks
- [ ] Something
EOF
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/bad.md" --dir "$STORIES"
    [ "$status" -eq 1 ]
    [[ "$output" == *"Invalid status"* ]]
}

@test "validate-doc: empty section detected" {
    cat > "$STORIES/bad.md" <<'EOF'
---
id: "0099"
title: "Bad"
status: draft
created: 2026-01-01
description: "Empty section"
---
# Bad

## Tasks
- [ ] Something

## Notes

## Other
Content here.
EOF
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/bad.md" --dir "$STORIES"
    [ "$status" -eq 1 ]
    [[ "$output" == *"Empty section: ## Notes"* ]]
}

@test "validate-doc: checkbox outside Tasks/Phase section" {
    cat > "$STORIES/bad.md" <<'EOF'
---
id: "0099"
title: "Bad"
status: draft
created: 2026-01-01
description: "Bad checkbox"
---
# Bad

## Overview
- [ ] Shouldn't be here

## Tasks
- [ ] Fine here
EOF
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/bad.md" --dir "$STORIES"
    [ "$status" -eq 1 ]
    [[ "$output" == *"Checkbox found outside Tasks/Phase section"* ]]
}

@test "validate-doc: no checkboxes at all" {
    cat > "$STORIES/bad.md" <<'EOF'
---
id: "0099"
title: "Bad"
status: draft
created: 2026-01-01
description: "No checkboxes"
---
# Bad

## Tasks
Just some text, no checkboxes.
EOF
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/bad.md" --dir "$STORIES"
    [ "$status" -eq 1 ]
    [[ "$output" == *"No checkboxes found"* ]]
}

@test "validate-doc: Phase sections with checkboxes are valid" {
    cat > "$STORIES/good.md" <<'EOF'
---
id: "0099"
title: "Phased"
status: draft
created: 2026-01-01
description: "Uses phases"
---
# Phased

## Phase 1: Setup
- [ ] Do setup

## Phase 2: Build
- [ ] Do build
EOF
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/good.md" --dir "$STORIES"
    [ "$status" -eq 0 ]
}

@test "validate-doc: done story needs non-empty developer logs" {
    cat > "$STORIES/bad.md" <<'EOF'
---
id: "0099"
title: "Done"
status: done
created: 2026-01-01
description: "Done but empty logs"
---
# Done

## Tasks
- [x] Finished

### Decision Log

### Lessons Learned
EOF
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/bad.md" --dir "$STORIES"
    [ "$status" -eq 1 ]
    [[ "$output" == *"Developer log section empty"*"Decision Log"* ]]
    [[ "$output" == *"Developer log section empty"*"Lessons Learned"* ]]
}

@test "validate-doc: done story with filled logs passes" {
    create_story 1 "done-good" "done"
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/0001-done-good.md" --dir "$STORIES"
    [ "$status" -eq 0 ]
}

@test "validate-doc: blocked-by references non-existent story" {
    create_story 1 "first" "draft" "9999"
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/0001-first.md" --dir "$STORIES"
    [ "$status" -eq 1 ]
    [[ "$output" == *"blocked-by references non-existent"* ]]
}

@test "validate-doc: blocked-by references existing story passes" {
    create_story 1 "first" "draft"
    create_story 2 "second" "draft" "0001"
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/0002-second.md" --dir "$STORIES"
    [ "$status" -eq 0 ]
}

@test "validate-doc: design depends-on references non-existent design" {
    create_design 1 "first" "draft" "9999"
    run "$SCRIPTS_DIR/validate-doc.sh" design "$DESIGNS/0001-first.md" --dir "$DESIGNS"
    [ "$status" -eq 1 ]
    [[ "$output" == *"depends-on references non-existent"* ]]
}

@test "validate-doc: cancelled status is valid" {
    create_story 1 "cancelled-story" "cancelled"
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/0001-cancelled-story.md" --dir "$STORIES"
    [ "$status" -eq 0 ]
}

@test "validate-doc: archived status is valid" {
    create_story 1 "archived-story" "archived"
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/0001-archived-story.md" --dir "$STORIES"
    [ "$status" -eq 0 ]
}

@test "validate-doc: detects cycle in blocked-by" {
    create_story 1 "first" "draft" "0002"
    create_story 2 "second" "draft" "0001"
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/0001-first.md" --dir "$STORIES"
    [ "$status" -eq 1 ]
    [[ "$output" == *"Cycle detected"* ]]
}

@test "validate-doc: no cycle in valid blocked-by chain" {
    create_story 1 "first" "draft"
    create_story 2 "second" "draft" "0001"
    create_story 3 "third" "draft" "0002"
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/0003-third.md" --dir "$STORIES"
    [ "$status" -eq 0 ]
}

@test "validate-doc: duplicate story id detected" {
    create_story 1 "first" "draft"
    create_story 1 "duplicate" "draft"
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/0001-first.md" --dir "$STORIES"
    [ "$status" -eq 1 ]
    [[ "$output" == *"Duplicate id"* ]]
}

@test "validate-doc: duplicate design id detected" {
    create_design 1 "first" "draft"
    create_design 1 "duplicate" "draft"
    run "$SCRIPTS_DIR/validate-doc.sh" design "$DESIGNS/0001-first.md" --dir "$DESIGNS"
    [ "$status" -eq 1 ]
    [[ "$output" == *"Duplicate id"* ]]
}

@test "validate-doc: unique ids pass" {
    create_story 1 "first" "draft"
    create_story 2 "second" "draft"
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/0001-first.md" --dir "$STORIES"
    [ "$status" -eq 0 ]
}

@test "validate-doc: nonexistent file" {
    run "$SCRIPTS_DIR/validate-doc.sh" story "$STORIES/nope.md" --dir "$STORIES"
    [ "$status" -eq 1 ]
    [[ "$output" == *"File not found"* ]]
}

@test "validate-doc: design with valid references passes" {
    git -C "$FIXTURES" init -q
    mkdir -p "$FIXTURES/src"
    echo "code" > "$FIXTURES/src/main.go"
    cat > "$DESIGNS/0001-refs.md" <<'EOF'
---
id: "0001"
title: "Refs"
status: draft
created: 2026-01-01
depends-on: []
references:
  - "src/main.go"
description: "Test references"
---
# Refs

## Tasks
- [ ] Task one
EOF
    run "$SCRIPTS_DIR/validate-doc.sh" design "$DESIGNS/0001-refs.md" --dir "$DESIGNS"
    [ "$status" -eq 0 ]
    [[ "$output" == "OK:"* ]]
}

@test "validate-doc: design with non-existent reference fails" {
    git -C "$FIXTURES" init -q
    cat > "$DESIGNS/0001-badrefs.md" <<'EOF'
---
id: "0001"
title: "Bad Refs"
status: draft
created: 2026-01-01
depends-on: []
references:
  - "does/not/exist.go"
description: "Bad references"
---
# Bad Refs

## Tasks
- [ ] Task one
EOF
    run "$SCRIPTS_DIR/validate-doc.sh" design "$DESIGNS/0001-badrefs.md" --dir "$DESIGNS"
    [ "$status" -eq 1 ]
    [[ "$output" == *"references non-existent file: does/not/exist.go"* ]]
}

@test "validate-doc: design with empty references passes" {
    git -C "$FIXTURES" init -q
    cat > "$DESIGNS/0001-norefs.md" <<'EOF'
---
id: "0001"
title: "No Refs"
status: draft
created: 2026-01-01
depends-on: []
references: []
description: "Empty references"
---
# No Refs

## Tasks
- [ ] Task one
EOF
    run "$SCRIPTS_DIR/validate-doc.sh" design "$DESIGNS/0001-norefs.md" --dir "$DESIGNS"
    [ "$status" -eq 0 ]
    [[ "$output" == "OK:"* ]]
}
