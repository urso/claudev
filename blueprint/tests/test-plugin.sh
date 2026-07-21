#!/usr/bin/env bash
set -euo pipefail

PLUGIN_DIR="$(cd "$(dirname "$0")/.." && pwd)"
TEST_DIR="${1:-/tmp/test-blueprint-$$}"

echo "=== Blueprint Plugin Test ==="
echo "Plugin: $PLUGIN_DIR"
echo "Test dir: $TEST_DIR"
echo

# Setup test directory
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

git init -q
git commit --allow-empty -m "init" -q

echo "=== Setup: Create test fixtures ==="
mkdir -p docs/ai/designs docs/ai/stories docs/ai/rules

cat > docs/ai/designs/0001-greeter.md << 'EOF'
---
id: "0001"
title: "Greeter CLI"
status: ready
created: 2026-07-21
depends-on: []
references: []
description: "Simple CLI that prints hello"
---

# Greeter CLI

## Problem
Need a CLI that greets users.

## Goals
- Print "Hello" when run

## Approach
Single Go binary with main function.
EOF

cat > docs/ai/stories/0001-0001-basic-greeting.md << 'EOF'
---
id: "0001"
title: "Basic Greeting"
status: ready
created: 2026-07-21
design: "0001"
blocked-by: []
description: "Implement basic greeting output"
---

# Basic Greeting

## Tasks

### Task 1: Create main.go
- [ ] Create main.go with hello output

## Technical Notes

Just `fmt.Println("Hello")`.

## Developer Logs

### Decision Log

### Blockers Encountered

### Deviations from Design
EOF

cat > docs/ai/stories/0001-0002-add-name.md << 'EOF'
---
id: "0002"
title: "Add Name Argument"
status: ready
created: 2026-07-21
design: "0001"
blocked-by: ["0001"]
description: "Add optional name argument"
---

# Add Name Argument

## Tasks

### Task 1: Parse name arg
- [ ] Add flag parsing for name

## Technical Notes

Use flag package.

## Developer Logs

### Decision Log

### Blockers Encountered

### Deviations from Design
EOF

cat > docs/ai/rules/build.md << 'EOF'
---
name: Build Configuration
applies-to: ["*"]
tags: [build]
description: Build commands
---

# Build

```bash
go build -o greeter .
```
EOF

echo "Created: design 0001, stories 0001 and 0002 (0002 blocked by 0001)"
git add -A && git commit -qm "add test fixtures"
echo

run_test() {
    local name="$1"
    local prompt="$2"
    echo "=== $name ==="
    echo "$prompt" | claude --plugin-dir "$PLUGIN_DIR" --dangerously-skip-permissions
    echo
    echo "---"
    echo
}

run_test "Test 1: /focus with no args" "/blueprint:focus"

echo "=== Check .claude-focus ==="
if [[ -f .claude-focus ]]; then
    echo "Focus file:"
    cat .claude-focus
else
    echo "No focus file yet"
fi
echo "---"
echo

run_test "Test 2: /focus design greeter" "/blueprint:focus design greeter"

echo "=== Check .claude-focus ==="
if [[ -f .claude-focus ]]; then
    echo "OK: Focus file created:"
    cat .claude-focus
else
    echo "FAIL: .claude-focus not created"
fi
echo "---"
echo

run_test "Test 3: /work (first story)" "/blueprint:work"

echo "=== Check .claude-focus after first /work ==="
if [[ -f .claude-focus ]]; then
    echo "Focus file:"
    cat .claude-focus
else
    echo "No focus file"
fi
echo "---"
echo

echo "=== Check story 0001 status ==="
grep -E "^status:" docs/ai/stories/0001-0001-basic-greeting.md || echo "No status found"
echo "---"
echo

run_test "Test 4: /work (should pick next story)" "/blueprint:work"

echo "=== Check .claude-focus after second /work ==="
if [[ -f .claude-focus ]]; then
    echo "Focus file:"
    cat .claude-focus
else
    echo "No focus file"
fi
echo "---"
echo

echo "=== Check story 0002 status ==="
grep -E "^status:" docs/ai/stories/0001-0002-add-name.md || echo "No status found"
echo "---"
echo

run_test "Test 5: /focus clear" "/blueprint:focus clear"

echo "=== Check cleared ==="
if [[ -f .claude-focus ]]; then
    echo "FAIL: Focus file still exists"
    cat .claude-focus
else
    echo "OK: Focus cleared"
fi
echo

echo "=== Test Complete ==="
echo "Test dir: $TEST_DIR"
