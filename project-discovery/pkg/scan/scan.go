// Package scan produces raw repository signals for the discover CLI.
//
// Two entry points:
//
//   - Repo: single tree walk producing extension counts, top-level directory
//     inventory, and presence of well-known marker files. Pure filesystem,
//     no shell-out.
//   - Churn: parses `git log --numstat` over a window to produce per-file
//     churn, global contributor stats, and per-top-level-directory ownership.
//
// Output is intentionally raw. Classification (languages, frameworks, build
// systems, entry points) is left to the agent reading these signals.
//
// JSON shape emitted by the CLI:
//
//	scan repo  -> {"repo":  RepoSignals}
//	scan churn -> {"churn": ChurnSignals}
package scan
