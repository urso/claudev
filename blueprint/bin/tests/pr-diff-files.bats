#!/usr/bin/env bats

load test_helper

setup() {
    FIXTURES="$(mktemp -d)"
    cd "$FIXTURES"
    git init -q
    git config user.email "test@test.com"
    git config user.name "Test"
    git commit --allow-empty -m "init" -q
}

teardown() {
    cd /
    rm -rf "$FIXTURES"
}

@test "pr-diff-files: no PR, shows main as base" {
    git checkout -b feature-branch -q
    echo "content" > file.go
    git add file.go
    git commit -m "add file" -q

    run "$SCRIPTS_DIR/pr-diff-files.sh"
    [ "$status" -eq 0 ]
    [[ "$output" == *"No PR (comparing against main)"* ]]
    [[ "$output" == *"Files:"* ]]
    [[ "$output" == *"file.go"* ]]
}

@test "pr-diff-files: includes line counts and chunks" {
    git checkout -b feature -q
    for i in {1..100}; do echo "line $i"; done > file.go
    git add file.go
    git commit -m "add" -q

    run "$SCRIPTS_DIR/pr-diff-files.sh"
    [ "$status" -eq 0 ]
    [[ "$output" == *"100 lines"* ]]
    [[ "$output" == *"Chunk 1"* ]]
}

@test "pr-diff-files: lists multiple changed files" {
    git checkout -b multi-file -q
    echo "a" > a.go
    echo "b" > b.go
    echo "c" > c.go
    git add .
    git commit -m "add files" -q

    run "$SCRIPTS_DIR/pr-diff-files.sh"
    [ "$status" -eq 0 ]
    [[ "$output" == *"Files: 3"* ]]
    [[ "$output" == *"a.go"* ]]
    [[ "$output" == *"b.go"* ]]
    [[ "$output" == *"c.go"* ]]
}

@test "pr-diff-files: no changes shows no files" {
    run "$SCRIPTS_DIR/pr-diff-files.sh"
    [ "$status" -eq 0 ]
    [[ "$output" == *"No PR (comparing against main)"* ]]
    [[ "$output" == *"No changed files"* ]]
}
