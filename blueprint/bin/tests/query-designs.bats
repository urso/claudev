#!/usr/bin/env bats

load test_helper

setup() { setup_fixtures; }
teardown() { teardown_fixtures; }

@test "query-designs: lists all designs" {
    create_design 1 "auth" "draft"
    create_design 2 "api" "ready"
    run "$SCRIPTS_DIR/query-designs.sh" --dir "$DESIGNS"
    [ "$status" -eq 0 ]
    [ "${#lines[@]}" -eq 2 ]
    [[ "${lines[0]}" == *"auth"* ]]
    [[ "${lines[1]}" == *"api"* ]]
}

@test "query-designs: filters by status" {
    create_design 1 "auth" "draft"
    create_design 2 "api" "ready"
    run "$SCRIPTS_DIR/query-designs.sh" --dir "$DESIGNS" --status ready
    [ "$status" -eq 0 ]
    [ "${#lines[@]}" -eq 1 ]
    [[ "${lines[0]}" == *"api"* ]]
}

@test "query-designs: filters by search" {
    create_design 1 "auth" "draft"
    create_design 2 "api" "draft"
    run "$SCRIPTS_DIR/query-designs.sh" --dir "$DESIGNS" --search auth
    [ "$status" -eq 0 ]
    [ "${#lines[@]}" -eq 1 ]
    [[ "${lines[0]}" == *"auth"* ]]
}

@test "query-designs: empty directory produces no output" {
    run "$SCRIPTS_DIR/query-designs.sh" --dir "$DESIGNS"
    [ "$status" -eq 0 ]
    [ "$output" = "" ]
}

@test "query-designs: nonexistent directory produces no output" {
    run "$SCRIPTS_DIR/query-designs.sh" --dir "$FIXTURES/nonexistent"
    [ "$status" -eq 0 ]
    [ "$output" = "" ]
}
