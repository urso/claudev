#!/bin/bash
set -e

RETRO_DIR="$HOME/.claude/retro"
WATERMARK_FILE="$RETRO_DIR/last_run"
DB_FILE="$RETRO_DIR/retro.duckdb"

mkdir -p "$RETRO_DIR"

# Parse arguments
SINCE_TS=""
while [[ $# -gt 0 ]]; do
  case $1 in
    --all)
      SINCE_TS="1970-01-01T00:00:00Z"
      echo "Mode: full history (--all)"
      shift
      ;;
    --since)
      SINCE_TS="$2"
      echo "Mode: since $SINCE_TS (--since)"
      shift 2
      ;;
    *)
      shift
      ;;
  esac
done

# Default: use watermark or full history
if [ -z "$SINCE_TS" ]; then
  if [ -f "$WATERMARK_FILE" ]; then
    SINCE_TS=$(cat "$WATERMARK_FILE")
    echo "Mode: since last run ($SINCE_TS)"
  else
    SINCE_TS="1970-01-01T00:00:00Z"
    echo "Mode: full history (no prior watermark)"
  fi
fi

echo "Loading messages since: $SINCE_TS"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Initialize DuckDB
duckdb "$DB_FILE" <<SQL
DROP TABLE IF EXISTS msgs;

CREATE TABLE msgs AS
SELECT *
FROM read_json_auto(
  '$HOME/.claude/projects/**/*.jsonl',
  union_by_name = true,
  sample_size = -1,
  filename = true
)
WHERE type IN ('user', 'assistant')
  AND isSidechain = false
  AND timestamp >= '${SINCE_TS}';

SELECT COUNT(*) AS loaded_messages FROM msgs;
SQL

# Load macros from separate file
duckdb "$DB_FILE" < "$SCRIPT_DIR/macros.sql"

# Update watermark
date -u +"%Y-%m-%dT%H:%M:%SZ" > "$WATERMARK_FILE"
echo "Watermark updated."
