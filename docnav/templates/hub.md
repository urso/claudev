---
title: ""
type: hub
tags: []
scope: ""
watches: []
updates-when:
  - a doc is added or removed from this subsystem
  - the subsystem's scope or boundaries change
agent-instructions: |
  1. Ask the user which subsystem this hub covers (e.g. "xatastor-operator",
     "CSI driver", "storage backend").
  2. Scan the docs directory for existing docs in that subsystem. Read their
     frontmatter (title, scope) to build the index.
  3. Draft the frontmatter: title, scope, tags, watches (the docs directory
     this hub indexes), and updates-when.
  4. Write a brief Overview (2-3 sentences) explaining what the subsystem is
     and what the linked docs cover.
  5. Write the Doc Map — a table or list linking each doc with a one-line
     summary derived from its scope field. Group by category (design vs
     reference) if there are enough docs to warrant it.
  6. If there are cross-cutting concerns not covered by any single doc,
     note them in a brief section or suggest new docs to fill the gaps.
  7. Present the draft for review. Adjust as requested.
---

# {title}

## Overview

<!-- agent: 2-3 sentences explaining what this subsystem is, its purpose, and
what the linked docs cover as a set. -->

## Docs

<!-- agent: Table or grouped list linking each doc in this subsystem. For each:
- Link to the doc (relative path)
- One-line summary (use the doc's scope field)
Group by category (design / reference) if the set is large enough. -->

## Gaps

<!-- agent: Optional — remove if all areas are covered. Note any cross-cutting
concerns or subsystem areas that don't have a doc yet. Suggest what to write. -->
