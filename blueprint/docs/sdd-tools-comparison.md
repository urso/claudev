# SDD Tools Comparison

This document compares Blueprint with other spec-driven development tools: [Kiro](https://kiro.dev), [Spec-kit](https://github.com/github/spec-kit), and [Tessl](https://tessl.io).

For background on SDD levels and criticisms, see Birgitta Böckeler's [analysis of SDD tools](https://martinfowler.com/articles/exploring-gen-ai/sdd-3-tools.html).

## SDD Levels

Spec-driven development exists on a spectrum:

| Level | Spec lifecycle | Code ownership |
|-------|---------------|----------------|
| **Spec-first** | Written before code, discarded after | Human edits code |
| **Spec-anchored** | Maintained alongside code, evolves with features | Human edits both spec and code |
| **Spec-as-source** | Only maintained artifact | Code is generated, never edited directly |

## Tool Comparison

| Aspect | Kiro | Spec-kit | Tessl | Blueprint |
|--------|------|----------|-------|-----------|
| **SDD level** | Spec-first | Spec-first (aspires anchored) | Spec-anchored/as-source | Spec-anchored |
| **Spec structure** | 3 docs: Requirements → Design → Tasks | Constitution + multiple files per spec | 1:1 spec-to-code mapping | 2-tier: Design + Stories |
| **Workflow flexibility** | Single workflow | Single workflow | Single workflow | Three modes (design-driven, story-first, ad-hoc) |
| **Spec style** | User stories, Given/When/Then | Verbose, multiple markdown files | Tags (@generate, @test) | Lean, "scannable in 30 seconds" |
| **Code ownership** | Human edits code | Human edits code | Generated (DO NOT EDIT) | Human edits code |
| **Sync mechanism** | Manual | Branch per spec | `tessl build` regenerates | `/story-update`, `/design-update` |
| **Learning capture** | None | Constitution (global rules) | None | Developer logs per story |

## When Iteration Happens

| Tool | Iteration timing |
|------|------------------|
| Kiro | Before implementation |
| Spec-kit | Before implementation |
| Tessl | Before implementation (then regenerate) |
| Blueprint | Before *and* during — stories adapt as work progresses |

Blueprint's dependency graph (`blocked-by`) enables mid-flight adaptation:
- Hit a blocker → insert a new story
- Scope changed → remove a story
- Discovered complexity → split a story

Design stays stable (intent), stories flex (execution).

## Addressing SDD Criticisms

Böckeler's [article](https://martinfowler.com/articles/exploring-gen-ai/sdd-3-tools.html) identifies key problems with current SDD tools:

| Criticism | Kiro/Spec-kit/Tessl | Blueprint |
|-----------|---------------------|-----------|
| **"Sledgehammer for small tasks"** — rigid workflows don't scale down | Single workflow for all problem sizes | Three modes: ad-hoc for quick tasks, story-first for medium work, design-driven for features |
| **"Rather review code than markdown"** — verbose specs are tedious | Multiple verbose markdown files | Lean specs; stories are checklists, designs scannable in 30 seconds |
| **"AI ignores specs anyway"** — elaborate specs don't guarantee compliance | No mechanism to track deviations | Developer logs capture what actually happened |
| **"Unclear problem size scope"** — where does SDD fit? | Unclear positioning | Modes scale from ad-hoc to full design based on problem certainty |
| **"Upfront spec contradicts iteration"** — conflicts with iterative development | Specs finalized before implementation | Specs refined through dialogue; stories added/removed mid-flight |
| **"MDD inflexibility + LLM non-determinism"** — worst of both worlds | Tessl: code is generated artifact | Code is primary artifact, not generated |

## Blueprint's Position

Blueprint is a **spec-anchored** approach that addresses common SDD criticisms:

1. **Workflow flexibility** — Three modes scale from quick fixes to complex features
2. **Lean specs** — Designs answer what/why briefly; stories track execution
3. **Learning capture** — Developer logs record decisions, blockers, deviations, lessons
4. **Mid-flight adaptation** — Stories can be added, removed, reordered during implementation
5. **Two-tier separation** — Design (strategic intent) vs Stories (tactical execution)

The key insight: specs should constrain where possible and guide where uncertain. Blueprint supports both through iterative refinement before execution and adaptation during execution.
