#!/usr/bin/env bats

load test_helper

setup() { setup_fixtures; }
teardown() { teardown_fixtures; }

@test "resolve-dir: rules type returns docs/ai/rules" {
    run bash -c "cd '$FIXTURES' && git init -q && '$SCRIPTS_DIR/resolve-dir.sh' rules"
    [ "$status" -eq 0 ]
    [[ "$output" == *"/docs/ai/rules" ]]
}

@test "resolve-dir: workflows type returns docs/ai/workflows" {
    run bash -c "cd '$FIXTURES' && git init -q && '$SCRIPTS_DIR/resolve-dir.sh' workflows"
    [ "$status" -eq 0 ]
    [[ "$output" == *"/docs/ai/workflows" ]]
}

@test "resolve-dir: design type returns docs/ai/designs" {
    run bash -c "cd '$FIXTURES' && git init -q && '$SCRIPTS_DIR/resolve-dir.sh' design"
    [ "$status" -eq 0 ]
    [[ "$output" == *"/docs/ai/designs" ]]
}

@test "resolve-dir: story type returns docs/ai/stories" {
    run bash -c "cd '$FIXTURES' && git init -q && '$SCRIPTS_DIR/resolve-dir.sh' story"
    [ "$status" -eq 0 ]
    [[ "$output" == *"/docs/ai/stories" ]]
}

@test "resolve-dir: unknown type fails" {
    run "$SCRIPTS_DIR/resolve-dir.sh" unknown
    [ "$status" -eq 1 ]
    [[ "$output" == *"Unknown type"* ]]
}

@test "resolve-dir: --dir override takes precedence" {
    run "$SCRIPTS_DIR/resolve-dir.sh" rules --dir "/custom/path"
    [ "$status" -eq 0 ]
    [ "$output" = "/custom/path" ]
}
