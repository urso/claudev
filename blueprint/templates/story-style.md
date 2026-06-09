---
name: Story Documents
applies-to: [stories]
tags: [style]
paths:
  - "docs/ai/stories/*.md"
  - "**/stories/*.md"
description: Guidelines for writing actionable story documents.
---

# Story Document Guidelines

## Core Principles

- **Be actionable** — every story is a unit of work with checkboxes that can be completed.
- **Track decisions** — the developer log captures what happened, not just what was planned.
- **Stay flexible** — use flat tasks or phased structure depending on complexity.
- **Metadata in frontmatter** — no overview, dependencies, or deliverables sections needed.

## Writing Style

### Do

- Write task items as imperative actions ("Implement X", "Add Y")
- Keep task items specific and completable in a single session
- Use sub-items (`  - [ ]`) for breaking down larger tasks
- Fill in developer log sections as you work, not all at the end

### Don't

- Add overview/summary sections (the title and description cover this)
- Add dependency sections (frontmatter `blocked-by` handles this)
- Add deliverables sections (completed tasks are the deliverables)
- Put checkboxes outside of Tasks/Phase sections
- Leave developer log sections as placeholders when marking status as `done`

## Acceptance Criteria

Each task should have an `#### Acceptance Criteria` section with verifiable outcomes, separate from subtasks.

### Writing Good Criteria

- **Code-based** — criteria should be verifiable by looking at the code ("interface X exists", "function handles error case Y")
- **Specific** — name files, types, functions, or behaviors ("has `Process(ctx, input) (output, error)` signature")
- **Testable assertions** — if tests are needed, describe what they verify ("unit test covers context cancellation")
- **Outcome-focused** — describe the end state, not the work steps

### What NOT to Include

- "Pass CI" — tests run automatically, this adds no value
- Vague criteria — "works correctly", "handles edge cases"
- Process steps — "write tests", "add docs" (those are subtasks)
- Redundant checks — if CI already enforces it, don't list it

### Who Marks Criteria Done

- **Subtasks**: marked by the development agent as work completes
- **Acceptance Criteria**: marked by the review process after verification — development agents must NOT mark these

## Structure Options

Choose based on complexity:
- **Flat**: single `## Tasks` section with checkboxes
- **Phased**: multiple `## Phase N: Name` sections, each with checkboxes

## Developer Log

Fill in these sub-sections during development (must be non-empty when status is `done`):
- **Decision Log**: technical decisions and rationale
- **Blockers Encountered**: issues and resolutions
- **Deviations from Plan**: changes from original design and why
- **Lessons Learned**: insights for future stories
