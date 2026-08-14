#!/usr/bin/env bats

load test_helper

setup() { setup_fixtures; }
teardown() { teardown_fixtures; }

@test "find-doc: finds story by id" {
    create_story 1 "first"
    run "$SCRIPTS_DIR/find-doc.sh" story 0001 --dir "$STORIES"
    [ "$status" -eq 0 ]
    [[ "$output" == *"0001-first.md" ]]
}

@test "find-doc: finds design by id" {
    create_design 3 "mydesign"
    run "$SCRIPTS_DIR/find-doc.sh" design 0003 --dir "$DESIGNS"
    [ "$status" -eq 0 ]
    [[ "$output" == *"0003-mydesign.md" ]]
}

@test "find-doc: finds story by unpadded id" {
    create_story 14 "fourteenth"
    run "$SCRIPTS_DIR/find-doc.sh" story 14 --dir "$STORIES"
    [ "$status" -eq 0 ]
    [[ "$output" == *"0014-fourteenth.md" ]]
}

@test "find-doc: finds design by unpadded id" {
    create_design 3 "mydesign"
    run "$SCRIPTS_DIR/find-doc.sh" design 3 --dir "$DESIGNS"
    [ "$status" -eq 0 ]
    [[ "$output" == *"0003-mydesign.md" ]]
}

@test "find-doc: exits 1 for missing id" {
    run "$SCRIPTS_DIR/find-doc.sh" story 9999 --dir "$STORIES"
    [ "$status" -eq 1 ]
    [[ "$output" == *"No document found"* ]]
}

@test "find-doc: exits 1 for nonexistent directory" {
    run "$SCRIPTS_DIR/find-doc.sh" story 0001 --dir "$FIXTURES/nonexistent"
    [ "$status" -eq 1 ]
    [[ "$output" == *"Directory not found"* ]]
}
