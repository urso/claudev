---
name: discover-scan
description: Scan repository for structure and code ownership signals
user-invocable: true
disable-model-invocation: true
allowed-tools: [Bash, Read]
---

# discover-scan

Scan the current repository and write `.discovery/scan.json` with raw signals about structure and code ownership.

## Variables

- **SCAN_OUTPUT**: `${CLAUDE_PLUGIN_ROOT}/resources/scan-output.md`

## Process

### 1. Run repo scan

```bash
discover scan repo
```

This writes structural signals: extension counts, top-level directory inventory, marker files present.

### 2. Run churn scan (if applicable)

Check if this is a git repository with history:

```bash
git rev-parse --is-shallow-repository 2>/dev/null
```

- If not a git repo: skip churn scan, note "Churn analysis skipped: not a git repository"
- If shallow clone (`true`): skip churn scan, note "Churn analysis skipped: shallow clone lacks history"
- Otherwise run:

```bash
discover scan churn
```

This adds code ownership signals: per-file churn, per-author lines changed, per-directory ownership shares, recency splits.

### 3. Report results

Read SCAN_OUTPUT for structure reference. Then read `.discovery/scan.json` and synthesize a 3-5 line summary covering:
- Primary languages detected (by extension count)
- Repository shape (monorepo vs single project, key directories)
- Top contributors / ownership concentration
- Any notable signals (large commits, recent activity patterns)

Report which files were written.

## Notes

- **Idempotent**: Safe to re-run anytime. Each scan overwrites the previous `scan.json`.
- **Raw signals only**: The scan produces facts, not conclusions.
- **Schema reference**: Run `discover scan schema` to see the annotated JSON structure with jq examples.
