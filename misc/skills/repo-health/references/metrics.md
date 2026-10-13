# Metric definitions and reading guide

All git metrics exclude merges, revert chains (revert commits and the commits they revert) and noise files (lockfiles, generated code, vendored code, mocks, binaries, `linguist-generated`/`linguist-vendored`, files with a "generated" header). Blank lines are never counted: churn, rework, blame and code size all count non-blank lines. Comments are counted.

## Intent

One label per commit: `feature`, `fix`, `refactor`, `perf`, `test`, `docs`, `build`, `chore`. Claude assigns them in `intent.tsv`. Unlabelled commits fall back to conventional-commit prefixes and keywords. Keywords mislabel a lot ("shutdown on io error" is often a fix), so label before you interpret.

Maintenance view (Swanson / ISO 14764): corrective = fix, adaptive = build + chore, perfective = feature + perf, preventive = refactor + test.

## Rework

Each replaced line is blamed at the commit's parent. A line younger than 21 days is *recent*, otherwise *legacy*. Rework % = recent ÷ (new + recent + legacy).

| Rework % | Reading |
|---|---|
| < 5 | Mature or stable code |
| 5–12 | Normal for active development |
| 12–20 | Young code or a design still settling |
| > 20 | Code merged before it was done, or a churning design |

- Intent × rework separates the cases. A refactor with high *legacy* is healthy maintenance. A feature or fix with high *recent* means it shipped half-done.
- Spikes in single months usually come from one consolidation. Check the top commits before calling it a trend.

## Fix share by area

Fix-commit lines ÷ all lines in the area. Ranking areas by this shows where defects cluster. Compare with commit counts: one large fix inflates the line share.

## Hotspots

Score = (revisions ÷ max) × (whitespace complexity ÷ max). Whitespace complexity is the sum of indentation levels, a language-agnostic proxy for cyclomatic complexity (Tornhill). Chores are not counted as revisions.

- High score with many recent revisions: active risk. Candidate for refactoring or extra tests.
- High score with no recent revisions: stable but expensive to touch.
- Test files in the list are normal. Check if their production counterpart is there too.

## Change coupling

Pairs of files that change in the same commit (2–30 files, excluding chore/build commits). Degree = shared ÷ average revisions of both files. Reported when shared ≥ 5 and degree ≥ 0.3.

- Same-area or test↔code pairs are expected.
- Cross-area pairs that are not tests point to hidden dependencies or a leaky abstraction.
- Code ↔ docs pairs are good: docs are being kept current.
- Config ↔ main ↔ chart pairs mean each new setting touches several layers.

## Knowledge distribution

- **Blame owners:** surviving lines at HEAD per author (`git blame -w`). Revert chains, `chore` commits and the repo's `.git-blame-ignore-revs` are passed through to the previous author.
- **History owners:** lines added per author over time, excluding chores, halving every 180 days.
- **Reviewers:** PR reviewers other than the author and bots, weighted by the PR's changed lines per area, with the same decay. Commits map to PRs via the `(#123)` squash-merge suffix. Logins map to git names through PR author ↔ commit author.
- **Reviewed %:** share of an area's PR-merged churn reviewed by someone else. Commits without a PR (direct pushes, unmerged local work) are left out.
- **Truck factor:** smallest number of blame owners that together own over half the lines.
- **Dominant commit:** the single commit owning the most surviving lines in the area, shown when ≥25%.
- **Inactive %:** lines owned by authors without a commit in the last 90 days.

Reading:
- Truck factor 1 is the norm for small teams. It becomes a risk when the area is also a hotspot or has a high fix share.
- High inactive % means knowledge has already left.

Caveats — call them out in the reading instead of trusting the numbers:
- Ownership measures who wrote or last rewrote code, not who understands it. Talk to the team before acting on it.
- A refactor transfers blame ownership, even a mostly mechanical one. Refactors are not filtered, since scripts can't tell mechanical from substantive. When a dominant commit is a refactor by someone other than the history owner, say ownership may be overstated for that author.
- `--ignore-revs` is heuristic. Lines a chore actually rewrote (e.g. a lint sweep adding annotations) can stay attributed to it.
- Editors and reviewers are separate knowledge levels. Editors know the code deeply. Reviewers know it partially. An area with one editor but a regular second reviewer is a smaller risk than one with one editor and low reviewed %. Only recommend spreading knowledge when both editors and reviewers are concentrated, or reviewed % is low.
- Every review counts the same. A one-word approval weighs as much as a deep review. Review comment counts are not collected.
- Blame and history owners usually agree. When they disagree, a recent rewrite moved ownership. Mention both.


## Code survival

Lines added per month that still exist at HEAD. The first months are usually low, since early code gets replaced. Recent months are naturally high. Look at the shape: survival should rise with age as the design settles.

## Pull request flow

Merged PRs, grouped by merge month. Reviews exclude the PR author.

- Cycle time p50 = typical turnaround. p90 = long-tail PRs (stacked, blocked or abandoned and revived).
- First review p50 = review latency.
- Unreviewed % = merged with no review from someone else. High values mean self-merge, pairing or bot PRs.
- Size p90 > 1–2K lines makes thorough review unlikely.
