---
id: "<ID>"
title: "<TITLE>"
status: draft
created: <DATE>
template: default
depends-on: []
references: []
description: "<ONE-LINE SUMMARY>"
agent-instructions: |
  1. Ask user to describe the problem they're solving and why it matters now.
  2. Clarify goals — what does success look like?
  3. Identify non-goals — what are we explicitly not doing?
  4. Discuss high-level approach before diving into details.
  5. For each section, draft content and confirm with user before moving on.
  6. Capture any unresolved questions in Open Questions.
---

# <TITLE>

## Problem Context

<!-- agent: Explain why this work is needed. Include business drivers,
technical debt, or user pain points. Keep to 2-3 paragraphs max. -->

## Goals

<!-- agent: List 2-5 concrete goals. Primary goal first, then secondary.
Each goal should be measurable or verifiable. -->

- <Primary goal>
- <Secondary goals>

## Non-Goals

<!-- agent: Explicitly state what this design does NOT cover. Helps scope
discussions and prevents feature creep. -->

- <What this design explicitly does NOT cover>

## Requirements

<!-- agent: Functional and non-functional requirements. Be specific enough
that stories can be derived from these. -->

## Constraints

<!-- agent: Testable assertions that must hold. Include performance limits,
security requirements, API contracts, business rules. Each constraint should
be verifiable — avoid vague statements. -->

- <Testable assertions>

## Technical Approach

<!-- agent: High-level strategy. Focus on the "how" at architecture level,
not implementation details. -->

### Entities

<!-- agent: List entities involved. For existing entities, note what changes.
For new entities, describe their purpose and key fields. -->

- **Existing**: <Entities touched by this design>
- **New**: <New entities introduced>

## Architecture Decisions

<!-- agent: Key technical decisions with rationale. Format as decision + why.
Include alternatives considered if relevant. -->

## Open Questions

<!-- agent: Unresolved issues to address before or during implementation.
Remove questions as they get answered. -->

- <Unresolved questions>

## References

<!-- agent: Links to external resources — RFCs, docs, Slack threads, related
designs, or prior art. -->

- <External links>
