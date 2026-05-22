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
bash "${CLAUDE_PLUGIN_ROOT}/scripts/discover.sh" scan repo
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
bash "${CLAUDE_PLUGIN_ROOT}/scripts/discover.sh" scan churn
```

This adds code ownership signals: per-file churn, per-author lines changed, per-directory ownership shares, recency splits.

### 3. Probe discovered artifacts

Read `.discovery/scan.json` and examine the `artifacts` array. For each artifact found:

**Documentation (docs, readme)**:
- Read the first 50 lines of key docs (README.md, ARCHITECTURE.md)
- Note what topics they cover

**Helm/K8s**:
- List charts or manifests present
- Note if there are multiple environments (dev, staging, prod)

**Playbooks/Runbooks**:
- List available playbooks
- Note their purpose from filenames or headers

**Scripts/Tools**:
- List key scripts
- Note if there's a Makefile with common targets

### 3b. Probe unfamiliar directories

List top-level directories:

```bash
ls -d */ | head -20
```

For each directory not already categorized by scan.json artifacts:
- Check for README, index, or config files
- Look at file types present
- Make a judgment: docs? configs? scripts? data? tests? vendor? unknown?

Add findings to the summary — the CLI catches common patterns, you catch project-specific ones.

### 4. Report results

Read SCAN_OUTPUT for structure reference. Synthesize a summary covering:
- Primary languages detected (by extension count)
- Repository shape (monorepo vs single project, key directories)
- Top contributors / ownership concentration
- **Documentation landscape** (what docs exist, what's covered)
- **Operational artifacts** (helm charts, k8s manifests, playbooks found)
- Any notable signals (large commits, recent activity patterns)

Report which files were written.

## Notes

- **Idempotent**: Safe to re-run anytime. Each scan overwrites the previous `scan.json`.
- **Raw signals only**: The scan produces facts, not conclusions.
- **Schema reference**: Run `bash "${CLAUDE_PLUGIN_ROOT}/scripts/discover.sh" scan schema` to see the annotated JSON structure with jq examples.
