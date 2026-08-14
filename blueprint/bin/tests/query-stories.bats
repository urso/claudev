#!/usr/bin/env bats

load test_helper

setup() { setup_fixtures; }
teardown() { teardown_fixtures; }

@test "query-stories: lists all stories" {
    create_story 1 "setup" "done"
    create_story 2 "auth" "draft"
    run "$SCRIPTS_DIR/query-stories.sh" --dir "$STORIES"
    [ "$status" -eq 0 ]
    [ "${#lines[@]}" -eq 2 ]
}

@test "query-stories: filters by status" {
    create_story 1 "setup" "done"
    create_story 2 "auth" "draft"
    run "$SCRIPTS_DIR/query-stories.sh" --dir "$STORIES" --status draft
    [ "$status" -eq 0 ]
    [ "${#lines[@]}" -eq 1 ]
    [[ "${lines[0]}" == *"auth"* ]]
}

@test "query-stories: filters by search" {
    create_story 1 "setup" "draft"
    create_story 2 "auth" "draft"
    run "$SCRIPTS_DIR/query-stories.sh" --dir "$STORIES" --search auth
    [ "$status" -eq 0 ]
    [ "${#lines[@]}" -eq 1 ]
    [[ "${lines[0]}" == *"auth"* ]]
}

@test "query-stories: actionable shows unblocked non-done stories" {
    create_story 1 "setup" "done"
    create_story 2 "auth" "ready" "0001"
    create_story 3 "api" "ready" "0002"
    run "$SCRIPTS_DIR/query-stories.sh" --dir "$STORIES" --actionable
    [ "$status" -eq 0 ]
    [ "${#lines[@]}" -eq 1 ]
    [[ "${lines[0]}" == *"auth"* ]]
}

@test "query-stories: actionable excludes done stories" {
    create_story 1 "setup" "done"
    run "$SCRIPTS_DIR/query-stories.sh" --dir "$STORIES" --actionable
    [ "$status" -eq 0 ]
    [ "$output" = "" ]
}

@test "query-stories: actionable shows stories with no blockers" {
    create_story 1 "standalone" "ready"
    run "$SCRIPTS_DIR/query-stories.sh" --dir "$STORIES" --actionable
    [ "$status" -eq 0 ]
    [ "${#lines[@]}" -eq 1 ]
    [[ "${lines[0]}" == *"standalone"* ]]
}

@test "query-stories: actionable excludes draft stories" {
    create_story 1 "draft-story" "draft"
    run "$SCRIPTS_DIR/query-stories.sh" --dir "$STORIES" --actionable
    [ "$status" -eq 0 ]
    [ "$output" = "" ]
}

@test "query-stories: actionable includes in-progress stories" {
    create_story 1 "wip" "in-progress"
    run "$SCRIPTS_DIR/query-stories.sh" --dir "$STORIES" --actionable
    [ "$status" -eq 0 ]
    [ "${#lines[@]}" -eq 1 ]
    [[ "${lines[0]}" == *"wip"* ]]
}

@test "query-stories: actionable excludes cancelled stories" {
    create_story 1 "nope" "cancelled"
    run "$SCRIPTS_DIR/query-stories.sh" --dir "$STORIES" --actionable
    [ "$status" -eq 0 ]
    [ "$output" = "" ]
}

@test "query-stories: empty directory produces no output" {
    run "$SCRIPTS_DIR/query-stories.sh" --dir "$STORIES"
    [ "$status" -eq 0 ]
    [ "$output" = "" ]
}
