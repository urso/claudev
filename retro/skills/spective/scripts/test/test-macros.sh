#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FIXTURE="$SCRIPT_DIR/fixture.jsonl"
MACROS="$SCRIPT_DIR/../macros.sql"
DB_FILE="/tmp/retro-test-$$.duckdb"
rm -f "$DB_FILE"

trap "rm -f $DB_FILE" EXIT

echo "=== Setting up test database ==="

duckdb "$DB_FILE" <<SQL
CREATE TABLE msgs AS
SELECT
  * EXCLUDE (timestamp),
  timestamp::TIMESTAMP as timestamp
FROM read_json_auto('$FIXTURE', union_by_name = true, sample_size = -1)
WHERE type IN ('user', 'assistant') AND isSidechain = false;

SELECT COUNT(*) AS msg_count FROM msgs;
SQL

echo ""
echo "=== Loading macros ==="
duckdb "$DB_FILE" < "$MACROS"

echo ""
echo "=== Test 1: edit_sequences should find auth.rs (3 edits, user messages between) ==="
RESULT=$(duckdb "$DB_FILE" -json -c "SELECT * FROM edit_sequences(3) WHERE file_path LIKE '%auth.rs';")
if echo "$RESULT" | grep -q "auth.rs"; then
  echo "PASS: Found auth.rs edit sequence"
else
  echo "FAIL: auth.rs not found in edit_sequences"
  echo "Result: $RESULT"
  exit 1
fi

echo ""
echo "=== Test 2: edit_sequences should NOT find handler.rs (no user messages between edits) ==="
RESULT=$(duckdb "$DB_FILE" -json -c "SELECT * FROM edit_sequences(3) WHERE file_path LIKE '%handler.rs';")
if echo "$RESULT" | grep -q "handler.rs"; then
  echo "FAIL: handler.rs should not appear (no user messages between consecutive edits)"
  exit 1
else
  echo "PASS: handler.rs correctly excluded"
fi

echo ""
echo "=== Test 3: thread_back should walk parent chain ==="
RESULT=$(duckdb "$DB_FILE" -c "SELECT COUNT(*) as cnt FROM thread_back('aaaaaaaa-0001-0001-0001-000000000007'::uuid, 10);")
if echo "$RESULT" | grep -q "7"; then
  echo "PASS: thread_back found 7 messages in chain"
else
  echo "FAIL: thread_back wrong count"
  echo "Result: $RESULT"
  exit 1
fi

echo ""
echo "=== Test 4: thread_forward should walk child chain ==="
RESULT=$(duckdb "$DB_FILE" -c "SELECT COUNT(*) as cnt FROM thread_forward('aaaaaaaa-0001-0001-0001-000000000001'::uuid, 10);")
if echo "$RESULT" | grep -q "7"; then
  echo "PASS: thread_forward found 7 messages in chain"
else
  echo "FAIL: thread_forward wrong count"
  echo "Result: $RESULT"
  exit 1
fi

echo ""
echo "=== All tests passed ==="
