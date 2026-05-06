---
description: Standard review criteria for design documents
---

# Design Review Criteria

## Structure Checks

- Has required sections: Problem Context, Goals, Non-Goals, Requirements, Constraints, Technical Approach, Architecture Decisions
- Frontmatter has required fields: id, title, status, created, description
- Uses concise bullet points, not long paragraphs
- No task lists or step-by-step implementation details (those belong in stories)

## Content Quality

### Problem Context
- Explains why this work is needed
- Clear business or technical driver
- Scoped appropriately (not too broad, not too narrow)

### Goals
- 2-5 concrete goals listed
- Goals are measurable or verifiable
- Primary goal is clear

### Non-Goals
- Explicitly states what is out of scope
- Non-goals don't contradict goals

### Requirements
- Specific enough to derive stories from
- Covers both functional and non-functional aspects
- No vague language ("fast", "secure" without criteria)

### Constraints
- Each constraint is testable/verifiable
- Includes relevant performance, security, or API requirements
- No vague assertions

### Technical Approach
- Addresses the requirements
- Architecture-level, not implementation details
- Entities section identifies existing and new entities

### Architecture Decisions
- Key decisions documented with rationale
- Alternatives considered where relevant

### Open Questions
- Unresolved issues captured
- No stale/already-resolved questions

## Consistency Checks

- Goals and requirements align
- Technical approach addresses all requirements
- No conflicting statements across sections
- Constraints are achievable given the approach
