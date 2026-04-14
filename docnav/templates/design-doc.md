---
title: ""
type: design-doc
tags: []
scope: ""
watches: []
updates-when: []
agent-instructions: |
  1. Ask the user which subsystem or area this doc covers (e.g. "xvol reconcilers",
     "CSI slot mode", "NVMe connection lifecycle"). Use their answer to scope the
     code scan.
  2. Scan the relevant source files — read package structure, key types, and
     function signatures. Build a mental model of the layers, data flow, and
     decision points before writing anything.
  3. Determine scope: if the subsystem is large, propose splitting into multiple
     focused docs with a central hub file that links them (like the operator docs
     split architecture / lifecycle / reconcilers / CRD). Ask the user to confirm
     the split before proceeding.
  4. For each doc (or the single doc if no split):
     a. Draft the frontmatter: title, scope, tags, watches (source files/dirs
        this doc tracks), and updates-when (events that should trigger a review).
     b. Write an Overview explaining what this area does and why it exists.
     c. Write the Design / How It Works section with narrative explanation.
        Use Mermaid diagrams (flowchart, sequence, state) to illustrate
        architecture, data flow, or state transitions.
     d. Write Key Decisions — each decision as a heading with the rationale
        and trade-offs. Only include decisions that aren't obvious from the code.
     e. If relevant, write Failure Modes / Edge Cases covering error paths,
        race conditions, and recovery mechanisms.
  5. Populate watches with the source file globs that back this doc.
     Populate tags with relevant domain terms.
  6. Present the draft for review. Adjust as requested.
---

# {title}

## Overview

<!-- agent: Explain what this area of the system does and why it exists. Motivate
the design — what problem does it solve, what did it replace, what constraints
shaped it. Keep it to 2-4 paragraphs. -->

## Design

<!-- agent: Narrative explanation of how the system works. Use Mermaid diagrams
(flowchart, sequence, stateDiagram) to illustrate architecture, data flow, or
state transitions. Walk through the layers or components and explain how they
compose. -->

## Key Decisions

<!-- agent: One subsection per non-obvious design decision. For each: state the
decision, the alternatives considered, and why this approach was chosen. Skip
decisions that are obvious from the code. -->

## Failure Modes

<!-- agent: Optional — remove if not applicable. Cover error paths, race
conditions, and recovery mechanisms. For each scenario: what triggers it, what
happens, and how the system recovers (or doesn't). -->
