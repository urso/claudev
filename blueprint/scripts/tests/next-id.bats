#!/usr/bin/env bats

load test_helper

setup() { setup_fixtures; }
teardown() { teardown_fixtures; }

@test "next-id: returns 0001 for empty directory" {
    run "$SCRIPTS_DIR/next-id.sh" story --dir "$STORIES"
    [ "$status" -eq 0 ]
    [ "$output" = "0001" ]
}

@test "next-id: returns 0001 for nonexistent directory" {
    run "$SCRIPTS_DIR/next-id.sh" story --dir "$FIXTURES/nonexistent"
    [ "$status" -eq 0 ]
    [ "$output" = "0001" ]
}

@test "next-id: returns next after highest id" {
    create_story 1 "first"
    create_story 5 "fifth"
    run "$SCRIPTS_DIR/next-id.sh" story --dir "$STORIES"
    [ "$status" -eq 0 ]
    [ "$output" = "0006" ]
}

@test "next-id: ignores files without id frontmatter" {
    create_story 3 "third"
    cat > "$STORIES/no-id.md" <<'EOF'
---
title: "No ID"
---
# No ID
EOF
    run "$SCRIPTS_DIR/next-id.sh" story --dir "$STORIES"
    [ "$status" -eq 0 ]
    [ "$output" = "0004" ]
}

@test "next-id: works for design type" {
    create_design 2 "mydesign"
    run "$SCRIPTS_DIR/next-id.sh" design --dir "$DESIGNS"
    [ "$status" -eq 0 ]
    [ "$output" = "0003" ]
}
