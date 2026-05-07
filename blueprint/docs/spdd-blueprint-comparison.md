# OpenSPDD vs Blueprint: A Detailed Comparison

This document compares two approaches to structured AI-assisted development: [OpenSPDD](https://github.com/gszhangwei/open-spdd) (Structured Prompt-Driven Development) and Blueprint.

Both recognize the core problem: **AI needs more than just requirements — it needs intent, constraints, and boundaries**. They solve it differently.

## Philosophy

**Spec-Driven Development (SDD)**: Specs are executable contracts that generate/constrain code. The spec is the thing being iterated — refined through AI dialogue until solid, then executed. "Specifications don't serve code — code serves specifications."

**OpenSPDD**: A specific SDD implementation. The REASONS Canvas captures seven dimensions (Requirements, Entities, Approach, Structure, Operations, Norms, Safeguards). Optionally, `/spdd-story` decomposes high-level requirements into INVEST-compliant stories before analysis. Refined through `/spdd-analysis` dialogue before implementation.

**Blueprint**: Also spec-driven, but with flexible structure. Designs capture intent; stories capture execution. Both are refined through conversation before implementation. Developer logs capture what actually happened — decisions, blockers, deviations, lessons — feeding forward to dependent work. Blueprint treats business analysis as a separate concern — it's encouraged to use external tools (BMAD, PRD generators, story mapping tools) for requirements gathering, then feed those artifacts into `/design-create`. Technical work can also start directly from requirements or ad-hoc.

## Document Structure

### OpenSPDD: Stories + Analysis + Canvas

OpenSPDD can use multiple document tiers:

**Stories** (optional, in `requirements/`):
- INVEST-compliant decomposition of high-level features
- Business-focused acceptance criteria (Given-When-Then)
- Generated via `/spdd-story`

**Analysis** (in `spdd/analysis/`):
- Strategic analysis: domain concepts, risks, gaps
- Generated via `/spdd-analysis`

**REASONS Canvas** (in `spdd/prompt/`):
- The implementation contract with seven dimensions:

| Section | Purpose |
|---------|---------|
| **R**equirements | Business goals, scope, acceptance criteria |
| **E**ntities | Domain model (Mermaid diagrams) |
| **A**pproach | Solution strategy, trade-offs |
| **S**tructure | Architecture, dependencies |
| **O**perations | Precise implementation tasks with method signatures |
| **N**orms | Coding standards, patterns |
| **S**afeguards | Constraints, what NOT to do |

The Canvas is the implementation contract. Stories and Analysis are optional but recommended for complex features.

### Blueprint: Design + Stories

Two-tier structure with separation of concerns:

**Design** (strategic):
- Problem Context
- Goals / Non-Goals
- Requirements / Constraints
- Technical Approach (including Entities)
- Architecture Decisions
- Open Questions

**Stories** (tactical):
- Tasks with checkboxes
- Technical Notes
- Developer Logs (decisions, blockers, deviations, lessons)

Designs can hold any level of detail — from high-level direction to full method specs. Stories handle execution tracking and capture learnings during implementation.

### Mapping REASONS to Blueprint

| REASONS | Blueprint Equivalent |
|---------|---------------------|
| **R**equirements | Requirements section in design |
| **E**ntities | Entities subsection in Technical Approach |
| **A**pproach | Technical Approach + Architecture Decisions |
| **S**tructure | Architecture Decisions |
| **O**perations | Stories (tasks + technical notes) |
| **N**orms | Style guides in `docs/ai/rules/` |
| **S**afeguards | Non-Goals + Constraints |

Blueprint splits Safeguards intentionally: "what we're not doing" (Non-Goals) vs "hard constraints" (Constraints) vs "coding rules" (style guides).

## Workflow

### OpenSPDD

```
Requirements
    │
    ▼
/spdd-story           Decompose into INVEST stories (optional)
    │
    ▼
/spdd-analysis        Strategic analysis, surface risks
    │
    ▼
/spdd-reasons-canvas  Generate REASONS Canvas
    │
    ▼
/spdd-generate        AI generates code per contract
    │
    ▼
/spdd-api-test        Generate curl-based API tests (optional)
    │
    ▼
Code review / refactor
    │
    ▼
/spdd-sync            Sync changes back to Canvas
    │
    ▼
/spdd-prompt-update   Update Canvas with new requirements (as needed)
```

Canvas is ground truth. Stories feed analysis. Code syncs back to Canvas. For simpler features, skip `/spdd-story` and go directly to `/spdd-analysis` or even `/spdd-reasons-canvas`.

### Blueprint

```
External input (PRDs, BMAD stories, specs, etc.) ─┐
                                                   │
Requirements ─────────────────────────────────────┼─→ /design-create
                                                   │
/analyze (explore codebase, surface risks) ───────┘
    │
    ▼
/design-review        Validate structure and clarity
    │
    ▼
/design-expand        Break into stories with dependencies
    │
    ▼
/develop-story        Implement tasks from a story
    │
    ▼
/story-update         Sync code changes to story, update dev logs
    │
    ▼
/design-update        Bubble learnings back to design
```

Two-tier sync: code → story → design. Stories can be added, removed, reordered mid-flight.

**Alternative flows**:

- **Story-first**: Skip design, create standalone stories for smaller work
- **Exploratory**: Lightweight design, discover through milestone stories, update design as you learn
- **Ad-hoc**: No tracking, just `/development` with style guide compliance
- **External input**: Feed PRDs, BMAD artifacts, or other specs directly into `/design-create`

## Key Differences

### Granularity

| | OpenSPDD | Blueprint |
|--|----------|-----------|
| **Spec unit** | REASONS Canvas (feature-level) | Story (task-level) |
| **Sync target** | Canvas ↔ code | Story ↔ code, Story → Design |
| **Pre-analysis** | `/spdd-story` (built-in, optional) | External tools (BMAD, PRDs, etc.) |

OpenSPDD's Canvas ≈ Blueprint's Design + all Stories in one document. Both support optional story decomposition, but Blueprint delegates that to external tooling.

### Context Management

| | OpenSPDD | Blueprint |
|--|----------|-----------|
| **Context load** | Full Canvas always loaded | Load design OR story as needed |
| **Tradeoff** | Complete picture, but heavy | Focused context, navigate between docs |

Blueprint's split is intentional: keep agent context clean from details it doesn't currently need. For large features, this matters.

### Dependencies and References

| | OpenSPDD | Blueprint |
|--|----------|-----------|
| **Work-item dependencies** | None — stories are INVEST "independent" | Explicit `blocked-by` graph between stories |
| **Cross-feature** | N/A | Design `depends-on` other designs |
| **Code dependencies** | Canvas "Structure" section (DI, component relationships) | Not tracked in docs |
| **Task ordering** | Operations list defines execution order | Dependency graph determines order |
| **File references** | `@file` in command input | Frontmatter `references:` + inline `@path` mentions |

SPDD tracks *code-level* dependencies within the Canvas. Blueprint tracks *work-item* dependencies between stories/designs.

Blueprint's reference system:
- **Frontmatter `references:`** — explicit file paths loaded when priming
- **Inline `@path/to/file`** — loose mentions anywhere in the doc, also picked up

Blueprint's dependency graph enables: add stories mid-flight, reorder work, handle blockers by inserting new stories. Dev logs from dependency stories are loaded into context when starting work.

### Handling Roadblocks

| | OpenSPDD | Blueprint |
|--|----------|-----------|
| **Blocker hit** | `/spdd-prompt-update` to revise Canvas | Add new story, reorder dependencies, continue |

OpenSPDD updates the Canvas in place. Blueprint's graph structure encourages adaptation by inserting/reordering stories.

### Drift Detection

| | OpenSPDD | Blueprint |
|--|----------|-----------|
| **Mechanism** | Manual `/spdd-sync`, AI diffs code vs Canvas | Git-sha based, automatic comparison |
| **Trigger** | User invokes sync | `/story-update` compares from `start-git-sha` |

Blueprint tracks when work started on a story, enabling precise drift detection.

### Analysis Phase

| | OpenSPDD | Blueprint |
|--|----------|-----------|
| **Command** | `/spdd-analysis` | `/analyze` |
| **Timing** | Before Canvas creation | Before design, before story, or on requirements |
| **Context** | Inline | Fresh sub-agent (no conversation bias) |

Both explore codebase against requirements to surface concepts, risks, gaps. Blueprint's `/analyze` works at both design and story level.

### Review

| | OpenSPDD | Blueprint |
|--|----------|-----------|
| **Spec review** | Implicit (Canvas is contract) | `/design-review`, `/story-review` with criteria files |
| **Code review** | `/spdd-code-review` (reviews against Canvas) | `/review-code` + reusable style rules |
| **API testing** | `/spdd-api-test` (generates curl scripts) | Not built-in |

Both have code review. OpenSPDD's reviews code against the Canvas for drift. Blueprint uses configurable style rules reusable across projects. OpenSPDD adds API test generation from acceptance criteria.

### Developer Logs

| | OpenSPDD | Blueprint |
|--|----------|-----------|
| **Learning capture** | Not structured | Developer Logs in every story |

Blueprint stories include:
- Decision Log
- Blockers Encountered
- Deviations from Design
- Lessons Learned

These feed forward to dependent stories and bubble up to designs via `/design-update`.

## Iteration Model

All three approaches iterate on specs through AI dialogue. The difference is structure and what gets captured.

### SDD / OpenSPDD

```
Requirements
    ↓
/spdd-story (optional: decompose into INVEST stories)
    ↓
/spdd-analysis (for each story or requirement)
    ↓
iterate: AI dialogue → refine Canvas → repeat
    ↓
Canvas is solid
    ↓
/spdd-generate (code from Canvas)
    ↓
/spdd-api-test (optional: generate test scripts)
    ↓
/spdd-sync (sync code back to Canvas)
```

Canvas is the contract. Stories and analysis are optional layers for complex features. `/spdd-prompt-update` handles evolving requirements.

### Blueprint

```
Requirements
    ↓
/design-create
    ↓
iterate: discuss → refine design → /design-review → repeat
    ↓
Design is solid
    ↓
/design-expand (break into stories)
    ↓
pick a story
    ↓
iterate: discuss impl → update story → clear context → reload → repeat
    ↓
Story is solid
    ↓
/develop-story
    ↓
/story-update (capture what actually happened)
    ↓
/design-update (bubble learnings back)
```

Two-tier structure. Iteration happens at both levels. Developer logs capture the journey, not just the destination.

**Why clear context and reload?** The story should be self-contained. If it only makes sense with conversation history, it's not ready.

## When to Use Which

### OpenSPDD excels when:

- You want an all-in-one methodology (story decomposition through code generation)
- Team prefers prescriptive structure (7 sections, fill them out)
- Business analysis and INVEST story decomposition are valuable
- You want method-level specs before coding
- API test generation from acceptance criteria is useful

### Blueprint excels when:

- You already have external tools for requirements/PRDs (BMAD, etc.)
- Features have complex dependencies or sequencing
- You want focused context (not everything loaded at once)
- Multiple people work on different stories in parallel
- You want to track learnings and decisions over time
- Technical work doesn't need deep business analysis

### Both work for:

- Iterative spec refinement through AI dialogue
- Code review and style compliance
- Syncing code changes back to specs
- Pre-implementation analysis

## Verdict

**OpenSPDD** is a comprehensive SDD implementation. Optional story decomposition, strategic analysis, seven-dimension Canvas, code generation, API testing, and bidirectional sync. All-in-one methodology from business requirements to tested code.

**Blueprint** is a composable SDD implementation. Design for intent, stories for execution, developer logs for learnings. Delegates business analysis to external tools (BMAD, PRDs, etc.) and focuses on the design-to-implementation loop. Flexible entry points for different levels of formality.

Both iterate specs through dialogue before implementation. OpenSPDD owns the full pipeline; Blueprint plugs into existing workflows.
