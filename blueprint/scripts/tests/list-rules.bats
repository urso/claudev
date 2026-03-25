#!/usr/bin/env bats

load test_helper

setup() { setup_fixtures; }
teardown() { teardown_fixtures; }

@test "list-rules: empty directory produces header only" {
    mkdir -p "$FIXTURES/rules"
    run "$SCRIPTS_DIR/list-rules.sh" "$FIXTURES/rules"
    [ "$status" -eq 0 ]
    [ "${lines[0]}" = "$FIXTURES/rules:" ]
    [ "${#lines[@]}" -eq 1 ]
}

@test "list-rules: nonexistent directory produces no output" {
    run "$SCRIPTS_DIR/list-rules.sh" "$FIXTURES/nonexistent"
    [ "$status" -eq 0 ]
    [ "$output" = "" ]
}

@test "list-rules: lists rule with full path" {
    create_rule "common" '["*"]' '[style]' "Universal style"
    run "$SCRIPTS_DIR/list-rules.sh" "$FIXTURES/rules"
    [ "$status" -eq 0 ]
    [ "${lines[0]}" = "$FIXTURES/rules:" ]
    [[ "${lines[1]}" == "$FIXTURES/rules/common.md|"* ]]
}

@test "list-rules: outputs name and description" {
    create_rule "common" '["*"]' '[style]' "Universal style"
    run "$SCRIPTS_DIR/list-rules.sh" "$FIXTURES/rules"
    [ "$status" -eq 0 ]
    [[ "${lines[1]}" == *"|Common Rules|"* ]]
    [[ "${lines[1]}" == *"|Universal style" ]]
}

@test "list-rules: filters by applies-to" {
    create_rule "go-rules" '["go"]' '[style]' "Go style"
    create_rule "py-rules" '["python"]' '[style]' "Python style"
    run "$SCRIPTS_DIR/list-rules.sh" "$FIXTURES/rules" "go"
    [ "$status" -eq 0 ]
    [ "${#lines[@]}" -eq 2 ]  # header + 1 match
    [[ "${lines[1]}" == *"go-rules.md|"* ]]
}

@test "list-rules: wildcard applies-to matches any filter" {
    create_rule "common" '["*"]' '[style]' "Universal"
    create_rule "go-rules" '["go"]' '[style]' "Go only"
    run "$SCRIPTS_DIR/list-rules.sh" "$FIXTURES/rules" "go"
    [ "$status" -eq 0 ]
    [ "${#lines[@]}" -eq 3 ]  # header + 2 matches
}

@test "list-rules: filters by tag" {
    create_rule "style-rules" '["*"]' '[style]' "Style"
    create_rule "build-rules" '["*"]' '[build]' "Build"
    run "$SCRIPTS_DIR/list-rules.sh" "$FIXTURES/rules" "" "style"
    [ "$status" -eq 0 ]
    [ "${#lines[@]}" -eq 2 ]  # header + 1 match
    [[ "${lines[1]}" == *"style-rules.md|"* ]]
}

@test "list-rules: filters by both applies-to and tag" {
    create_rule "go-style" '["go"]' '[style]' "Go style"
    create_rule "go-build" '["go"]' '[build]' "Go build"
    create_rule "py-style" '["python"]' '[style]' "Python style"
    run "$SCRIPTS_DIR/list-rules.sh" "$FIXTURES/rules" "go" "style"
    [ "$status" -eq 0 ]
    [ "${#lines[@]}" -eq 2 ]  # header + 1 match
    [[ "${lines[1]}" == *"go-style.md|"* ]]
}

@test "list-rules: skips files without frontmatter" {
    create_rule "valid" '["*"]' '[style]' "Valid"
    echo "# No frontmatter" > "$FIXTURES/rules/nofm.md"
    run "$SCRIPTS_DIR/list-rules.sh" "$FIXTURES/rules"
    [ "$status" -eq 0 ]
    [ "${#lines[@]}" -eq 2 ]  # header + 1 valid rule
}
