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

@test "diff-files: includes deleted files" {
    echo "content" > file.go
    git add file.go
    git commit -m "add" -q
    git checkout -b feature -q
    rm file.go
    git add file.go
    git commit -m "delete" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"file.go"* ]]
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

@test "diff-files: separates code, config, docs, ignored sections" {
    git checkout -b feature -q
    echo "code" > main.go
    echo "config" > config.yaml
    echo "docs" > README.md
    echo "checksum" > go.sum
    git add .
    git commit -m "add" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"## Code"* ]]
    [[ "$output" == *"## Config"* ]]
    [[ "$output" == *"## Docs"* ]]
    [[ "$output" == *"## Ignored"* ]]
    [[ "$output" == *"go.sum"*"lockfile"* ]]
}

@test "diff-files: groups config files by chart root" {
    git checkout -b feature -q
    mkdir -p charts/api/templates
    echo "name: api" > charts/api/Chart.yaml
    echo "replicas: 1" > charts/api/values.yaml
    echo "kind: Deployment" > charts/api/templates/deploy.yaml
    git add .
    git commit -m "add" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"Group charts/api (chart)"* ]]
}

@test "diff-files: adds unchanged values.yaml as context" {
    mkdir -p charts/api/templates
    echo "name: api" > charts/api/Chart.yaml
    echo "replicas: 1" > charts/api/values.yaml
    git add .
    git commit -m "chart" -q

    git checkout -b feature -q
    echo "kind: Deployment" > charts/api/templates/deploy.yaml
    git add .
    git commit -m "template" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"context: charts/api/values.yaml"* ]]
    [[ "$output" == *"context: charts/api/Chart.yaml"* ]]
}

@test "diff-files: subchart wins over parent chart" {
    git checkout -b feature -q
    mkdir -p charts/parent/charts/child/templates
    echo "name: parent" > charts/parent/Chart.yaml
    echo "name: child" > charts/parent/charts/child/Chart.yaml
    echo "kind: Service" > charts/parent/charts/child/templates/svc.yaml
    git add .
    git commit -m "add" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"Group charts/parent/charts/child (chart)"* ]]
}

@test "diff-files: config without chart goes to loose group" {
    git checkout -b feature -q
    mkdir -p .github/workflows
    echo "on: push" > .github/workflows/ci.yml
    git add .
    git commit -m "add" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"Group loose"* ]]
    [[ "$output" == *".github/workflows/ci.yml"* ]]
}

@test "diff-files: detects generated files by header comment" {
    git checkout -b feature -q
    printf '// Code generated by protoc. DO NOT EDIT.\npackage main\n' > api.go
    git add .
    git commit -m "add" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"api.go"*"generated"* ]]
    [[ "$output" != *"## Code"* ]]
}

@test "diff-files: detects generated files by name pattern" {
    git checkout -b feature -q
    echo "package main" > types.pb.go
    git add .
    git commit -m "add" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"types.pb.go"*"generated"* ]]
}

@test "diff-files: tags deleted files" {
    echo "content" > file.go
    git add file.go
    git commit -m "add" -q
    git checkout -b feature -q
    git rm -q file.go
    git commit -m "delete" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"file.go"*"(deleted)"* ]]
}

@test "diff-files: tags added files" {
    git checkout -b feature -q
    echo "content" > new.go
    git add new.go
    git commit -m "add" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"new.go"*"(added)"* ]]
}

@test "diff-files: resolves renames to the new path" {
    echo "package main" > old.go
    git add old.go
    git commit -m "add" -q
    git checkout -b feature -q
    git mv old.go new.go
    git commit -m "rename" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"new.go"* ]]
    [[ "$output" != *"=>"* ]]
}

@test "diff-files: unchanged files are not tagged" {
    git checkout -b feature -q
    echo "a" > a.go
    git add a.go
    git commit -m "add" -q
    echo "b" >> a.go
    git add a.go
    git commit -m "modify" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" != *"(modified)"* ]]
}

@test "diff-files: chunk total sums add+delete per file" {
    for i in $(seq 1 100); do echo "line $i"; done > a.go
    git add a.go
    git commit -m "init" -q
    git checkout -b feature -q
    for i in $(seq 1 100); do echo "CHANGED $i"; done > a.go
    git add a.go
    git commit -m "rewrite" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    # 100 added + 100 deleted = 200
    [[ "$output" == *"200 lines total"* ]]
    [[ "$output" == *"Chunk 1 (200 lines)"* ]]
    [[ "$output" == *"+100 -100"* ]]
}

@test "diff-files: shows add/delete split for modified files" {
    printf 'a\nb\nc\n' > f.go
    git add f.go
    git commit -m "init" -q
    git checkout -b feature -q
    printf 'a\nB\nc\nd\n' > f.go
    git add f.go
    git commit -m "edit" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"+2 -1"* ]]
}

@test "diff-files: no add/delete split for added or deleted files" {
    git checkout -b feature -q
    echo "x" > new.go
    git add new.go
    git commit -m "add" -q

    run "$SCRIPTS_DIR/diff-files.sh" main
    [ "$status" -eq 0 ]
    [[ "$output" == *"(added)"* ]]
    [[ "$output" != *"+1 -0"* ]]
}
