#!/bin/bash
# Wrapper for the docnav tool.
# Builds the Go binary on first run, then passes through to the CLI.
# Usage: docnav.sh <subcommand> [flags...]

set -euo pipefail

CLAUDE_PLUGIN_ROOT="${CLAUDE_PLUGIN_ROOT:-$(cd "$(dirname "$0")/.." && pwd)}"

if ! command -v go &>/dev/null; then
  echo "error: Go is not installed. docnav plugin requires Go." >&2
  exit 1
fi

# Build the binary if missing or stale
bin="${CLAUDE_PLUGIN_ROOT}/.bin/docnav"
needs_build=false
if [ ! -f "$bin" ]; then
  needs_build=true
elif [ -n "$(find "$CLAUDE_PLUGIN_ROOT" -name '*.go' -newer "$bin" -print -quit 2>/dev/null)" ] || [ "$CLAUDE_PLUGIN_ROOT/go.mod" -nt "$bin" ]; then
  needs_build=true
fi
if [ "$needs_build" = true ]; then
  mkdir -p "${CLAUDE_PLUGIN_ROOT}/.bin"
  tmp=$(mktemp "${bin}.XXXXXX")
  if go build -C "$CLAUDE_PLUGIN_ROOT" -o "$tmp" ./cmd/docnav/; then
    mv "$tmp" "$bin"
  else
    rm -f "$tmp"
    exit 1
  fi
fi

exec "$bin" "$@"
