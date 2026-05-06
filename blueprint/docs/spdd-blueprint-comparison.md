# OpenSPDD vs Blueprint: A Detailed Comparison

This document compares two approaches to structured AI-assisted development: [OpenSPDD](https://github.com/gszhangwei/open-spdd) (Structured Prompt-Driven Development) and Blueprint.

Both recognize the core problem: **AI needs more than just requirements — it needs intent, constraints, and boundaries**. They solve it differently.

## Philosophy

**OpenSPDD**: Spec upfront, contract-driven. Assumes you can and should specify enough detail before implementation that AI executes precisely. The REASONS Canvas is a contract.

**Blueprint**: Adaptive to problem certainty. When you know the problem space, spec fully upfront. When exploring, start lightweight and discover through implementation. The structure adapts to what you know.

## Document Structure

### OpenSPDD: Single Canvas

One REASONS Canvas per feature containing seven dimensions:

| Section | Purpose |
|---------|---------|
| **R**equirements | Business goals, scope, acceptance criteria |
| **E**ntities | Domain model (Mermaid diagrams) |
| **A**pproach | Solution strategy, trade-offs |
| **S**tructure | Architecture, dependencies |
| **O**perations | Precise implementation tasks with method signatures |
| **N**orms | Coding standards, patterns |
| **S**afeguards | Constraints, what NOT to do |

Everything lives in one document. Complete picture, but can get large.

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
/spdd-analysis        Strategic analysis, surface risks
    │
    ▼
/spdd-reasons-canvas  Generate REASONS Canvas
    │
    ▼
/spdd-generate        AI generates code per contract
    │
    ▼
Code review / refactor
    │
    ▼
/spdd-sync            Sync changes back to Canvas
```

Linear flow. Canvas is ground truth. Code syncs back to Canvas.

### Blueprint

```
Requirements
    │
    ▼
/analyze              Explore codebase, surface risks (optional)
    │
    ▼
/design-create        Create design through guided questions
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

## Key Differences

### Granularity

| | OpenSPDD | Blueprint |
|--|----------|-----------|
| **Spec unit** | REASONS Canvas (feature-level) | Story (task-level) |
| **Sync target** | Canvas ↔ code | Story ↔ code, Story → Design |

OpenSPDD's Canvas ≈ Blueprint's Design + all Stories in one document.

### Context Management

| | OpenSPDD | Blueprint |
|--|----------|-----------|
| **Context load** | Full Canvas always loaded | Load design OR story as needed |
| **Tradeoff** | Complete picture, but heavy | Focused context, navigate between docs |

Blueprint's split is intentional: keep agent context clean from details it doesn't currently need. For large features, this matters.

### Dependencies

| | OpenSPDD | Blueprint |
|--|----------|-----------|
| **Task ordering** | Flat Operations list | Story dependency graph (`blocked-by`) |
| **Cross-feature** | N/A | Design can depend on other designs |

Blueprint's dependency chains enable: add stories mid-flight, reorder work, handle blockers by inserting new stories.

### Handling Roadblocks

| | OpenSPDD | Blueprint |
|--|----------|-----------|
| **Blocker hit** | Rethink/rewrite Canvas | Add new story, reorder dependencies, continue |

SPDD's flat structure encourages starting over. Blueprint's graph structure encourages adaptation.

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
| **Code review** | `/spdd-code-review` (optional) | `/review-code` + reusable style rules |

Blueprint has explicit review skills with configurable criteria. Style rules are reusable across projects.

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

## When to Use Which

### OpenSPDD excels when:

- Requirements are well-understood upfront
- You want a single document with everything
- Team prefers prescriptive structure (7 sections, fill them out)
- Feature is self-contained, limited dependencies
- You want method-level specs before coding

### Blueprint excels when:

- Problem certainty varies (some known, some exploratory)
- Features have complex dependencies or sequencing
- You discover requirements through implementation
- You want focused context (not everything loaded at once)
- Multiple people work on different stories in parallel
- You want to track learnings and decisions over time

### Both work for:

- Design-first development
- Code review and style compliance
- Syncing code changes back to specs
- Pre-implementation analysis

## Verdict

**OpenSPDD** is a contract. Single document, seven dimensions, precise specs. Best when you know what you're building.

**Blueprint** is a framework that lets you choose your contract level. Fully rigid (detailed design + simple stories) or fully exploratory (standalone stories, discover as you go). Best when problem certainty varies.

The less you know about the problem space, the more Blueprint's flexibility helps. The more you know, the more SPDD's structure enforces completeness.

Neither is wrong — they optimize for different points on the certainty spectrum.
