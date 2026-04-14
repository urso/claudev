---
title: ""
type: reference-doc
tags: []
scope: ""
watches: []
updates-when: []
agent-instructions: |
  1. Ask the user which subsystem or component this reference covers (e.g.
     "XVol CRD", "reconciler set", "CSI driver RPCs"). Use their answer to
     scope the code scan.
  2. Scan the relevant source files — read type definitions, field tags,
     validation rules, interfaces, and constants. Extract the exhaustive list
     of items to document.
  3. Determine scope: if the component surface is large, propose splitting into
     multiple focused reference docs with a central hub file that links them.
     Ask the user to confirm before proceeding.
  4. For each doc (or the single doc if no split):
     a. Draft the frontmatter: title, scope, tags, watches (source files/dirs
        this doc tracks), and updates-when (events that should trigger a review).
     b. Write a brief Overview (1-2 paragraphs) stating what this reference
        covers and how to use it.
     c. Write the Components / Fields section — one subsection per item.
        Use tables for field listings (name, type, description, constraints).
        Include Go interface definitions where relevant.
        Use Mermaid diagrams for state machines, decision trees, or flow.
     d. If multiple components interact, write a Relationships section
        explaining how they fit together (with a diagram if helpful).
  5. Populate watches with the source file globs that back this doc.
     Populate tags with relevant domain terms.
  6. Present the draft for review. Adjust as requested.
---

# {title}

## Overview

<!-- agent: Brief intro (1-2 paragraphs) stating what this reference covers,
its scope, and how to read it. -->

## Components

<!-- agent: One subsection per component, field group, or API surface. For each:
- Summary of purpose
- Table of fields/parameters (name, type, description, constraints)
- Go interface or type definition if applicable
- Mermaid diagram if the component has a state machine or decision tree
Repeat this subsection pattern for each item. -->

## Relationships

<!-- agent: Optional — remove if not applicable. Explain how the documented
components interact. Use a Mermaid diagram showing event flow, ownership, or
call relationships. -->
