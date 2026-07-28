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

@test "diff-files: no changes shows no files" {
    run "$SCRIPTS_DIR/diff-files.sh"
    [ "$status" -eq 0 ]
    [[ "$output" == "No changed files" ]]
}

@test "diff-files: shows unstaged changes without base" {
    echo "content" > file.go
    git add file.go
    git commit -m "add" -q
    echo "modified" >> file.go

    run "$SCRIPTS_DIR/diff-files.sh"
    [ "$status" -eq 0 ]
    [[ "$output" == *"Files: 1"* ]]
    [[ "$output" == *"file.go"* ]]
}

@test "diff-files: shows branch changes with base" {
    git checkout -b feature -q
    echo "content" > file.go
    git add file.go
    git commit -m "add" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"Files: 1"* ]]
    [[ "$output" == *"file.go"* ]]
}

@test "diff-files: shows line counts" {
    git checkout -b feature -q
    for i in {1..50}; do echo "line $i"; done > file.go
    git add file.go
    git commit -m "add" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"50 lines"* ]]
    [[ "$output" == *"file.go"*"50"* ]]
}

@test "diff-files: chunks files exceeding budget" {
    git checkout -b feature -q

    # Create files that exceed 500 line budget
    for i in {1..300}; do echo "line $i"; done > a.go
    for i in {1..300}; do echo "line $i"; done > b.go
    for i in {1..300}; do echo "line $i"; done > c.go

    git add .
    git commit -m "add" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"900 lines total"* ]]
    [[ "$output" == *"Chunk 1"* ]]
    [[ "$output" == *"Chunk 2"* ]]
}

@test "diff-files: groups by language then directory" {
    git checkout -b feature -q
    mkdir -p src/auth src/api
    echo "go code" > src/auth/login.go
    echo "go code" > src/api/handler.go
    echo "ts code" > src/api/client.ts
    git add .
    git commit -m "add" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"Files: 3"* ]]
    # All should be in one chunk (small files)
    [[ "$output" == *"Chunk 1"* ]]
}

@test "diff-files: skips deleted files" {
    echo "content" > file.go
    git add file.go
    git commit -m "add" -q
    git checkout -b feature -q
    rm file.go
    git add file.go
    git commit -m "delete" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == "No changed files" ]]
}

@test "diff-files: handles nested directories" {
    git checkout -b feature -q
    mkdir -p src/pkg/auth
    echo "code" > src/pkg/auth/handler.go
    git add .
    git commit -m "add" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"src/pkg/auth/handler.go"* ]]
}
