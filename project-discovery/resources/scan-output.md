# Scan Output

How to read and interpret `.discovery/scan.json`.

## Structure

The scan output has two top-level sections:

### repo

Structural signals from `discover scan repo`:

- `root`: absolute path to repository root
- `file_count`: total files scanned
- `extensions`: array of `{ext, count}` sorted by count descending
- `top_dirs`: array of `{dir, files}` for top-level directories
- `markers`: array of marker filenames found (Makefile, go.mod, package.json, etc.)

### churn

Code ownership signals from `discover scan churn` (only present if git history available):

- `by_file`: per-file churn stats
- `by_author`: per-author lines changed
- `by_dir`: per-top-level-directory ownership shares
- `large_commits`: notable commits (capped at 20)

## Example jq Queries

Top 3 languages by file count:
```bash
jq '.repo.extensions[:3]' .discovery/scan.json
```

List marker files:
```bash
jq -r '.repo.markers[]' .discovery/scan.json
```

Top contributors (if churn data present):
```bash
jq '.churn.by_author | sort_by(-.lines) | .[:5]' .discovery/scan.json
```

Directory with most files:
```bash
jq '.repo.top_dirs[0]' .discovery/scan.json
```
