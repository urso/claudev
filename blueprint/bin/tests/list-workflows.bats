#!/usr/bin/env bats

load test_helper

setup() { setup_fixtures; }
teardown() { teardown_fixtures; }

@test "list-workflows: empty directory produces header only" {
    mkdir -p "$FIXTURES/workflows"
    run "$SCRIPTS_DIR/list-workflows.sh" "$FIXTURES/workflows"
    [ "$status" -eq 0 ]
    [ "${lines[0]}" = "$FIXTURES/workflows:" ]
    [ "${#lines[@]}" -eq 1 ]
}

@test "list-workflows: nonexistent directory produces no output" {
    run "$SCRIPTS_DIR/list-workflows.sh" "$FIXTURES/nonexistent"
    [ "$status" -eq 0 ]
    [ "$output" = "" ]
}

@test "list-workflows: lists workflow with full path" {
    create_workflow "design" '[design]' "Design workflow guidelines"
    run "$SCRIPTS_DIR/list-workflows.sh" "$FIXTURES/workflows"
    [ "$status" -eq 0 ]
    [ "${lines[0]}" = "$FIXTURES/workflows:" ]
    [[ "${lines[1]}" == "$FIXTURES/workflows/design.md|"* ]]
}

@test "list-workflows: outputs name and description" {
    create_workflow "design" '[design]' "Design workflow guidelines"
    run "$SCRIPTS_DIR/list-workflows.sh" "$FIXTURES/workflows"
    [ "$status" -eq 0 ]
    [[ "${lines[1]}" == *"|Design Workflow|"* ]]
    [[ "${lines[1]}" == *"|Design workflow guidelines" ]]
}

@test "list-workflows: filters by type tag" {
    create_workflow "design" '[design]' "Design guidelines"
    create_workflow "story" '[story]' "Story guidelines"
    run "$SCRIPTS_DIR/list-workflows.sh" "$FIXTURES/workflows" "design"
    [ "$status" -eq 0 ]
    [ "${#lines[@]}" -eq 2 ]  # header + 1 match
    [[ "${lines[1]}" == *"design.md|"* ]]
}

@test "list-workflows: no filter lists all" {
    create_workflow "design" '[design]' "Design guidelines"
    create_workflow "story" '[story]' "Story guidelines"
    run "$SCRIPTS_DIR/list-workflows.sh" "$FIXTURES/workflows"
    [ "$status" -eq 0 ]
    [ "${#lines[@]}" -eq 3 ]  # header + 2 workflows
}

@test "list-workflows: skips files without frontmatter" {
    create_workflow "design" '[design]' "Valid"
    echo "# No frontmatter" > "$FIXTURES/workflows/nofm.md"
    run "$SCRIPTS_DIR/list-workflows.sh" "$FIXTURES/workflows"
    [ "$status" -eq 0 ]
    [ "${#lines[@]}" -eq 2 ]  # header + 1 valid workflow
}
