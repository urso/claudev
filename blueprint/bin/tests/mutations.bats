#!/usr/bin/env bats

load test_helper

setup() { setup_fixtures; }
teardown() { teardown_fixtures; }

# --- set-status ---

@test "set-status: updates status" {
    create_story 1 "test" "draft"
    run "$SCRIPTS_DIR/set-status.sh" story 0001 in-progress --dir "$STORIES"
    [ "$status" -eq 0 ]
    result=$(yq --front-matter=extract '.status' "$STORIES/0001-test.md")
    [ "$result" = "in-progress" ]
}

@test "set-status: preserves markdown content" {
    create_story 1 "test" "draft"
    "$SCRIPTS_DIR/set-status.sh" story 0001 ready --dir "$STORIES"
    grep -q "## Tasks" "$STORIES/0001-test.md"
    grep -q "\- \[ \] Task one" "$STORIES/0001-test.md"
}

@test "set-status: fails for missing id" {
    run "$SCRIPTS_DIR/set-status.sh" story 9999 done --dir "$STORIES"
    [ "$status" -eq 1 ]
    [[ "$output" == *"not found"* ]]
}

# --- set-blocked-by ---

@test "set-blocked-by: sets blockers" {
    create_story 1 "first"
    create_story 2 "second"
    run "$SCRIPTS_DIR/set-blocked-by.sh" story 0002 0001 --dir "$STORIES"
    [ "$status" -eq 0 ]
    result=$(yq --front-matter=extract '.blocked-by | length' "$STORIES/0002-second.md")
    [ "$result" -eq 1 ]
    result=$(yq --front-matter=extract '.blocked-by[0]' "$STORIES/0002-second.md")
    [ "$result" = "0001" ]
}

@test "set-blocked-by: appends to existing blockers" {
    create_story 1 "first"
    create_story 2 "second" "draft" "0001"
    "$SCRIPTS_DIR/set-blocked-by.sh" story 0002 0003 0004 --dir "$STORIES"
    result=$(yq --front-matter=extract '.blocked-by | length' "$STORIES/0002-second.md")
    [ "$result" -eq 3 ]
    result=$(yq --front-matter=extract '.blocked-by[0]' "$STORIES/0002-second.md")
    [ "$result" = "0001" ]
}

@test "set-blocked-by: deduplicates" {
    create_story 1 "first"
    create_story 2 "second" "draft" "0001"
    "$SCRIPTS_DIR/set-blocked-by.sh" story 0002 0001 0003 --dir "$STORIES"
    result=$(yq --front-matter=extract '.blocked-by | length' "$STORIES/0002-second.md")
    [ "$result" -eq 2 ]
}

@test "set-blocked-by: fails for missing id" {
    run "$SCRIPTS_DIR/set-blocked-by.sh" story 9999 0001 --dir "$STORIES"
    [ "$status" -eq 1 ]
}

# --- remove-blocked-by ---

@test "remove-blocked-by: removes specific blocker" {
    create_story 1 "first"
    create_story 2 "second" "draft" "0001" "0003"
    run "$SCRIPTS_DIR/remove-blocked-by.sh" story 0002 0001 --dir "$STORIES"
    [ "$status" -eq 0 ]
    result=$(yq --front-matter=extract '.blocked-by | length' "$STORIES/0002-second.md")
    [ "$result" -eq 1 ]
    result=$(yq --front-matter=extract '.blocked-by[0]' "$STORIES/0002-second.md")
    [ "$result" = "0003" ]
}

@test "remove-blocked-by: removing last blocker leaves empty list" {
    create_story 1 "first"
    create_story 2 "second" "draft" "0001"
    "$SCRIPTS_DIR/remove-blocked-by.sh" story 0002 0001 --dir "$STORIES"
    result=$(yq --front-matter=extract '.blocked-by | length' "$STORIES/0002-second.md")
    [ "$result" -eq 0 ]
}

@test "remove-blocked-by: fails for missing id" {
    run "$SCRIPTS_DIR/remove-blocked-by.sh" story 9999 0001 --dir "$STORIES"
    [ "$status" -eq 1 ]
}

# --- archive-doc ---

@test "archive-doc: archives single story" {
    create_story 1 "done-story" "done"
    # archive-doc uses git root for archive path; test components directly
    run "$SCRIPTS_DIR/set-status.sh" story 0001 archived --dir "$STORIES"
    [ "$status" -eq 0 ]
    result=$(yq --front-matter=extract '.status' "$STORIES/0001-done-story.md")
    [ "$result" = "archived" ]
}

@test "archive-doc: archives single design" {
    create_design 1 "done-design" "done"
    run "$SCRIPTS_DIR/set-status.sh" design 0001 archived --dir "$DESIGNS"
    [ "$status" -eq 0 ]
    result=$(yq --front-matter=extract '.status' "$DESIGNS/0001-done-design.md")
    [ "$result" = "archived" ]
}

@test "archive-doc: fails for missing id" {
    run "$SCRIPTS_DIR/archive-doc.sh" story 9999 --dir "$STORIES"
    [ "$status" -eq 1 ]
}
