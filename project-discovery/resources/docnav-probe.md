# Docnav Probe

Determine if docnav is available for enhanced search.

## How to Check

Inspect your available-skills list (shown in `<system-reminder>` at conversation start) for a `doc-nav` skill entry.

## Search Tiers

**Tier 1 (always available)**: Built-in `discover ticket search` — per-field weighted TF-IDF over stemmed tokens.

**Tier 2 (when docnav present)**: Delegate to `doc-nav` skill with `--path .discovery/` for stemmed full-text and ranking.

## Reporting

Tell the user which tier is active:
- "Docnav detected — Tier 2 search available (full-text + ranking)"
- "Docnav not detected — using Tier 1 search (ticket fields only)"
