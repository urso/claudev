---
name: Design Documents
applies-to: [designs]
tags: [style]
paths:
  - "docs/ai/designs/*.md"
  - "**/designs/*.md"
description: Guidelines for writing lean design documents.
---

# Design Document Guidelines

## Core Principles

- **Be lean** — problem, goals, approach, decisions. No implementation tasks (those live in stories).
- **Be concise** — bullet points over paragraphs. Readers should grasp intent in 30 seconds.
- **Scope explicitly** — always include Non-Goals to prevent scope creep.
- **Decide up front** — capture key technical decisions and rationale in the design, not later.

## Writing Style

### Do

- Use bullet points for requirements and decisions
- Keep items to one line when possible
- Include code examples only for interfaces/contracts
- Reference other designs by ID number

### Don't

- Write implementation tasks (save for stories)
- Use long explanatory paragraphs
- Include full algorithm implementations
- Specify internal data structures in detail
- Add dependency/deliverable sections (frontmatter and stories handle those)

## Architecture Decisions Format

For each decision, include:
- **Decision**: what was decided
- **Rationale**: why this over alternatives
- **Alternatives considered**: what was rejected and why
