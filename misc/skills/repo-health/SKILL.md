---
description: Project health report from git history and GitHub PRs - churn by intent, rework rate, fix share, hotspots, change coupling, knowledge distribution, code survival, PR flow. Use when asked about code churn, project/repo health, hotspots, bus factor, or development stability.
user-invocable: true
argument-hint: "[--alias handle=Full Name]..."
allowed-tools: ["Read", "Write", "Bash(python3:*)", "Bash(git:*)", "Bash(gh:*)", "Bash(open:*)"]
---

# Repository health report

Zero-config. Run from inside the target repository. Data is cached in `<git-common-dir>/repo-health/`, so it is never committed and is shared across worktrees. Each run writes `runs/<date>/` with `metrics.json`, `report.html` and `report.md`.

`S=${CLAUDE_PLUGIN_ROOT}/skills/repo-health/scripts`

Read `references/metrics.md` before interpreting results.

## User input
```
$ARGUMENTS
```

## Steps

1. **Collect history.** `python3 -I $S/health.py collect`. Incremental: line-age (blame at parent) is cached per commit.
2. **Collect PRs.** `python3 -I $S/health.py prs`. Needs `gh` auth. On failure, tell the user and continue without PR flow.
3. **Classify intent.** `python3 -I $S/health.py pending > /tmp/repo-health-pending.tsv` lists unlabelled commits (revert chains excluded).
   - If there are any, label each one and append `<sha>\t<intent>` lines to `$(python3 -I $S/health.py dir)/intent.tsv`.
   - Intents: `feature`, `fix`, `refactor`, `perf`, `test`, `docs`, `build`, `chore`. See `references/metrics.md`.
   - Judge from subject, body and touched dirs, not keywords alone. Behavior corrections are `fix` even without the word. Dependency bumps, CI, Helm and Nix are `build`.
   - For more than ~50 commits, delegate to a subagent with the definitions above. Spot-check its output.
4. **Compute.** `python3 -I $S/health.py metrics [--alias "handle=Full Name"]...` prints the run dir.
   - Identities sharing a name or email merge automatically. `.mailmap` is honored.
   - If `summary`/`ownership` show the same person under two names, rerun with `--alias`.
5. **Read** `<run>/metrics.json`. Write `<run>/reading.md`:
   - 5–10 bullets, ordered by importance, each naming the area/file and the number.
   - Separate stages: early design churn vs ongoing instability. Weigh fixes and recent rework over feature volume.
   - Cross-check signals. A hotspot that also has a high fix share and truck factor 1 is the top risk.
   - Ownership: check `dominant_commit` per area. If one refactor decides the blame owner, say ownership may be overstated (see caveats in `references/metrics.md`). Ask the user when it matters for a conclusion.
   - End with 2–4 concrete follow-ups.
6. **Render.** `python3 -I $S/report.py <run>`. Then `open <run>/report.html` on macOS.
7. **Reply** with the report path and the reading. Don't restate the tables.

## Notes

- Areas are derived from the current tree: directories are split until each holds ≤15% of the code. Files in deleted directories map to their first two path components.
- To compare with an earlier run, diff `runs/*/report.md` or the `metrics.json` summaries.
- Commits rewritten by rebase get new SHAs. Their intent labels are lost and need relabelling.
