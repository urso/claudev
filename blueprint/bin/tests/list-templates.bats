#!/usr/bin/env bats

load test_helper

setup() {
    setup_fixtures
    git -C "$FIXTURES" init -q
}
teardown() { teardown_fixtures; }

@test "list-templates: default template always listed" {
    run bash -c "cd '$FIXTURES' && '$SCRIPTS_DIR/list-templates.sh' design"
    [ "$status" -eq 0 ]
    [[ "$output" == *"default|"* ]]
}

@test "list-templates: outputs name|path|description format" {
    run bash -c "cd '$FIXTURES' && '$SCRIPTS_DIR/list-templates.sh' design"
    [ "$status" -eq 0 ]
    # Should have exactly 2 pipes per line
    local pipes
    pipes=$(echo "${lines[0]}" | tr -cd '|' | wc -c)
    [ "$pipes" -eq 2 ]
}

@test "list-templates: lists user templates" {
    create_template "migration" "Template for data migrations"
    run bash -c "cd '$FIXTURES' && '$SCRIPTS_DIR/list-templates.sh' design"
    [ "$status" -eq 0 ]
    [[ "$output" == *"migration|"* ]]
    [[ "$output" == *"Template for data migrations"* ]]
}

@test "list-templates: user templates listed before default" {
    create_template "migration" "Migration template"
    run bash -c "cd '$FIXTURES' && '$SCRIPTS_DIR/list-templates.sh' design"
    [ "$status" -eq 0 ]
    # migration should appear before default
    local migration_line default_line
    migration_line=$(echo "$output" | grep -n "migration|" | cut -d: -f1)
    default_line=$(echo "$output" | grep -n "default|" | cut -d: -f1)
    [ "$migration_line" -lt "$default_line" ]
}

@test "list-templates: multiple user templates" {
    create_template "migration" "Migration template"
    create_template "refactor" "Refactor template"
    run bash -c "cd '$FIXTURES' && '$SCRIPTS_DIR/list-templates.sh' design"
    [ "$status" -eq 0 ]
    [[ "$output" == *"migration|"* ]]
    [[ "$output" == *"refactor|"* ]]
    [[ "$output" == *"default|"* ]]
}

@test "list-templates: skips directories without template.md" {
    mkdir -p "$FIXTURES/docs/ai/templates/designs/empty"
    echo "not a template" > "$FIXTURES/docs/ai/templates/designs/empty/readme.md"
    run bash -c "cd '$FIXTURES' && '$SCRIPTS_DIR/list-templates.sh' design"
    [ "$status" -eq 0 ]
    [[ "$output" != *"empty|"* ]]
}
