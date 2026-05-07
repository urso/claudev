# AI Workflow Comparison

This document compares Blueprint with other Claude Code workflow methodologies: [BMAD](https://github.com/bmad-code-org/BMAD-METHOD) and [Superpowers](https://github.com/obra/superpowers).

All three are spec-anchored approaches built for AI-assisted development. They differ in scope, structure, and how they handle learning.

## Overview

| Aspect | BMAD | Superpowers | Blueprint |
|--------|------|-------------|-----------|
| **Focus** | Full agile methodology | Execution with TDD | Lean spec-anchored workflow |
| **Lifecycle** | Analysis → Plan → Solution → Implementation → Retro | Brainstorm → Design → Plan → Execute → Finish | Design → Stories → Implement → Learn |
| **Agent model** | 12+ specialized personas | Subagent per task | Single agent, fresh context per skill |
| **Spec structure** | PRD → Epics → Stories | Design doc → Plan doc | Design → Stories (two-tier) |

## Artifacts

| Artifact type | BMAD | Superpowers | Blueprint |
|---------------|------|-------------|-----------|
| **Requirements/PRD** | Yes — full PRD with user stories | Design doc from brainstorming | Design doc (problem, goals, approach) |
| **Architecture** | Yes — dedicated docs | Part of design doc | Part of design (Technical Approach) |
| **Task breakdown** | Epics → Stories → Tasks → Subtasks | Plan → Tasks with Do/Verify sections | Design → Stories → Tasks → Subtasks |
| **Decision logs** | Dev agent record, retrospectives | Subagent self-review (not persisted) | Developer logs per story |
| **Sprint tracking** | Yes — YAML status files | No | Story status + dependency graph |

## Task Granularity

| Tool | Structure | Format |
|------|-----------|--------|
| **BMAD** | Story → Tasks → Subtasks | Checkboxes, linked to acceptance criteria (`- [ ] Task 1 (AC: #)`) |
| **Superpowers** | Plan → Tasks | Each task has **Do:** (actions) and **Verify:** (checks) sections |
| **Blueprint** | Story → Tasks → Subtasks | Checkboxes, grouped by phase or area |

BMAD and Blueprint use similar checkbox-based task/subtask structure. Superpowers uses structured Do/Verify sections — more prescriptive about verification.

## Dependency Tracking

| Tool | Dependency model |
|------|------------------|
| **BMAD** | Implicit — epic ordering + story sequence within epic. Stories within an epic must not depend on future stories. No explicit cross-epic story dependencies. |
| **Superpowers** | Implicit — sequential task list in plan |
| **Blueprint** | Explicit — `blocked-by` in story frontmatter. Cross-story and cross-design dependencies supported. |

Blueprint's explicit dependencies enable:
- Non-linear story ordering
- Inserting stories mid-flight without rewriting the plan
- Visualizing the dependency graph
- Cross-design dependencies for multi-feature work

## Spec Maintenance

| Tool | Specs maintained? | Sync mechanism | Learning capture |
|------|-------------------|----------------|------------------|
| **BMAD** | Yes | Stories updated during impl, `bmad-correct-course` for changes, retrospectives bubble up | Dev agent record, retrospective action items |
| **Superpowers** | Partially | Specs committed but no update mechanism post-impl | Subagent reviews (session only) |
| **Blueprint** | Yes | `/story-update`, `/design-update` | Developer logs (decisions, blockers, deviations, lessons) |

## When Decisions Are Captured

All three have agent-written logs — the developer confirms but doesn't write manually.

| Tool | During discussion | During implementation | Post-implementation |
|------|-------------------|----------------------|---------------------|
| **BMAD** | PRD/epic creation | Dev agent record (debug log, completion notes) | Retrospectives |
| **Superpowers** | Design doc | Subagent self-review (session only, not persisted) | No |
| **Blueprint** | Story iteration (before code) | Developer logs (decisions, blockers, deviations) | `/design-update` bubbles learnings |

**Key Blueprint distinction:** Decisions are captured during the *discussion phase*, not just implementation. The iterate → update story → clear context → reload loop means reasoning is recorded before code exists — why this approach, what was considered, what was rejected.

## Mid-Flight Adaptation

| Tool | How changes are handled |
|------|------------------------|
| **BMAD** | `bmad-correct-course` — analyzes impact across PRD, epics, architecture; produces Sprint Change Proposals |
| **Superpowers** | Not explicit — plan is consumed by execution |
| **Blueprint** | Add/remove/reorder stories via dependency graph; design stays stable, stories flex |

## Document Reviews

All three have intermediate review phases for generated documents — not just code review.

| Tool | Document review | Customizable |
|------|-----------------|--------------|
| **BMAD** | PRD validation, architecture review, story review before implementation | Checklists are built-in per skill |
| **Superpowers** | Design self-review (placeholder scan, consistency check, ambiguity check) before user approval | Prompt templates are built-in |
| **Blueprint** | `/design-review` and `/story-review` — check structure, clarity, consistency; usable after creation and after updates | Review rules overridable via `docs/ai/workflows/` + per-template review criteria |

Blueprint's document reviews catch:
- Inconsistencies between sections
- Ambiguous requirements
- Missing constraints or non-goals
- Tasks that don't match the technical approach

Reviews are repeatable — run after any update to ensure quality.

## Customization

| Aspect | BMAD | Superpowers | Blueprint |
|--------|------|-------------|-----------|
| **Review criteria** | Checklists per skill (built-in) | Prompt templates (built-in) | Rules in `docs/ai/workflows/`, overridable per project |
| **Skill/workflow customization** | 3-layer TOML merge: defaults → team → user overrides | User preferences override defaults | Rules loaded based on `applies-to` tags and paths |
| **Document templates** | Fixed templates per artifact type | Fixed templates | Custom design templates, each with own review rules |
| **Derive new templates** | No | No | Yes — derive templates from existing designs |

BMAD's customization is powerful — the 3-layer TOML merge allows team-wide and personal overrides for agents, menus, and persistent facts.

Blueprint focuses on design template customization — teams can create project-specific design templates with tailored review criteria, then derive new templates from existing designs as patterns emerge.

## Code Reviews

| Aspect | BMAD | Superpowers | Blueprint |
|--------|------|-------------|-----------|
| **Approach** | Parallel adversarial layers | Single reviewer, structured template | Parallel specialized reviewers |
| **Reviewers** | Blind Hunter (no context), Edge Case Hunter (project access), Acceptance Auditor (spec + context) | One reviewer checking: plan alignment, code quality, architecture, testing, production readiness | Style reviewer, Bug reviewer (opus), Efficiency reviewer |
| **Validation** | Triage step categorizes findings | Reviewer categorizes by severity | False positive validation pass with tiered models (haiku/sonnet/opus) |
| **Spec integration** | Acceptance Auditor checks against spec | Reviewer checks plan alignment | `--story` flag for story context |
| **Output** | Findings by reviewer layer | Strengths + Issues (Critical/Important/Minor) + Assessment | Merged report by severity + style guide proposals |
| **Customizable** | 3-layer TOML | Template-based | Reference files + style guides |

BMAD's adversarial approach is unique — Blind Hunter gets no project context at all, forcing it to review code purely on its merits. Blueprint's false positive validation pass reduces noise. Superpowers focuses on single comprehensive review with clear verdict.

## TDD

| Aspect | BMAD | Superpowers | Blueprint |
|--------|------|-------------|-----------|
| **Approach** | Red-green-refactor built into dev-story workflow | Dedicated TDD skill with strict enforcement | Not mandated — developer decides per story |
| **Enforcement** | Default practice in implementation step | Strict: "NO PRODUCTION CODE WITHOUT A FAILING TEST FIRST" — delete code and restart if violated | Optional |
| **Exceptions** | Customizable via 3-layer TOML | Only with human partner approval (throwaway prototypes, generated code, config files) | Developer discretion when crafting stories |
| **Philosophy** | Integrated default | Dogmatic with extensive rationalization handling ("All of these mean: Delete code. Start over with TDD.") | Pragmatic — TDD useful sometimes, overkill other times |

Superpowers has the most opinionated TDD stance — a full skill dedicated to enforcing and defending TDD against "rationalizations." BMAD integrates TDD as default practice. Blueprint leaves it to developer judgment, allowing TDD tasks in stories when appropriate.

## Workflow Flexibility

| Aspect | BMAD | Superpowers | Blueprint |
|--------|------|-------------|-----------|
| **Modes** | Full methodology + Quick Dev for small changes | Single workflow (brainstorm → plan → execute → finish) | Three modes: design-driven, story-first, ad-hoc |
| **Lightweight path** | `bmad-quick-dev` — intent in, code out, minimal checkpoints | None — brainstorming required before implementation | Ad-hoc mode — describe task, auto-load rules, implement |
| **Entry points** | Any agent (PM, Architect, Dev), Quick Dev, or full method | Always starts with brainstorming unless subagent | `/design-create`, `/story-create`, or just describe the task |
| **Skip planning** | Quick Dev routes small changes straight to implementation | Not supported — brainstorming is mandatory gate | Ad-hoc mode skips design/story entirely |
| **Scale adaptation** | Quick Dev auto-routes: small changes → direct impl, larger → planning | Single workflow for all sizes | User chooses mode based on task size/certainty |

**BMAD Quick Dev** is notable — it compresses intent first, routes to smallest safe path (direct impl vs full planning), then runs longer with less supervision. Review findings are triaged to diagnose failures at the right layer (intent, spec, or implementation).

**Superpowers** has a hard gate: "NO implementation until design approved." Brainstorming must happen first, even for small changes. The philosophy is that skipping planning is rationalization.

**Blueprint** lets the user choose their entry point based on what they know:
- **Design-driven**: Full design → stories → implement (high certainty, complex features)
- **Story-first**: Skip design, create story directly (medium certainty, focused work)
- **Ad-hoc**: Just describe the task, auto-load rules, implement (quick fixes, small changes)

## Combining Workflows

These tools can complement each other:

**BMAD → Blueprint:** Use BMAD for ideation, PRDs, and architecture. Use Blueprint for execution with lean stories and learning capture. BMAD's brainstorming and PRD generation work well as input for Blueprint designs — the agent can help define required designs to establish initial macro-milestones before breaking down into designs and stories.

**All three work for greenfield and brownfield.** The difference is how much structure they provide upfront:

| Scenario | BMAD | Superpowers | Blueprint |
|----------|------|-------------|-----------|
| **Greenfield** | Full methodology guides from zero — PRD, architecture, epics, stories | Brainstorming + TDD gives strong foundation | Works well but requires developer guidance on initial project shape (tooling, structure, conventions) |
| **Brownfield** | Quick Dev for targeted changes, full method for larger features | Same workflow applies | Mid-flight adaptation and learning capture help when requirements evolve |

**Developer guidance in Blueprint:** For greenfield, the developer typically has ideas about project structure, tooling choices, and conventions before diving into features. Blueprint expects this context — either through discussion during design iteration, or captured in style guides. BMAD's methodology guides these decisions more explicitly through its agent personas (Architect, UX, etc.).

## Planning Overhead

The frameworks differ significantly in how much planning is required before implementation begins.

### BMAD Full Method

Before implementation, BMAD can involve:

1. **Analysis phase:** Document project, product brief, PRFAQ, research
2. **Plan phase:** Create PRD (multi-step facilitated workflow), validate PRD, create UX design
3. **Solutioning phase:** Create architecture, create epics and stories, check implementation readiness, generate project context

Each step is a facilitated workflow with multiple checkpoints and human approval gates. The full method provides comprehensive guidance but requires significant upfront investment.

**Quick Dev** bypasses most of this — compresses intent, routes small changes directly to implementation.

### Superpowers

Before implementation:

1. **Brainstorming:** Explore context → clarifying questions (one at a time) → propose 2-3 approaches → present design sections → write design doc → self-review → user reviews spec
2. **Writing plans:** Scope check → file structure mapping → bite-sized tasks (2-5 min each, full code in every step) → self-review

Brainstorming is mandatory — no implementation without approved design, even for "simple" projects. Plans are detailed with complete code in every step.

### Blueprint

Before implementation:

- **Design-driven:** `/design-create` (guided questions) → iterate with discussion → `/design-review` → `/design-expand` (generates stories)
- **Story-first:** `/story-create` → iterate with discussion → `/story-review`
- **Ad-hoc:** Describe task → auto-load rules → create plan → get approval → implement

Blueprint minimizes required planning and allows deferring it. Ad-hoc mode goes straight to implementation with just a plan approval. Design and story iteration can be as light or detailed as the developer chooses.

### Comparison

| Aspect | BMAD | Superpowers | Blueprint |
|--------|------|-------------|-----------|
| **Minimum path to code** | Quick Dev (intent → impl) | Brainstorming + plan (mandatory) | Ad-hoc (task → plan → impl) |
| **Full planning steps** | ~10+ facilitated workflows | 2 workflows, many substeps | 2-3 commands with iteration |
| **Mandatory gates** | Approval at each phase | Design approval, plan approval | Plan approval only |
| **Plan detail level** | Epics → Stories → Tasks with ACs | Tasks with full code, 2-5 min each | Stories → Tasks → Subtasks (checkboxes) |
| **Deferrable** | Via Quick Dev | No — brainstorming required | Yes — start with ad-hoc, add structure later |

### Discovery-Driven Development

Sometimes you can't plan ahead — especially in greenfield projects where the real design emerges during development. Plans fall apart when requirements are uncertain.

| Aspect | BMAD | Superpowers | Blueprint |
|--------|------|-------------|-----------|
| **Exploration support** | Brainstorming workflow for ideation; Quick Dev for small experiments | "Need to explore first? Fine. Throw away exploration, start with TDD." | Built-in — defer planning, discover through implementation, flow learnings back to design |
| **Prototype-then-rewrite** | Not explicit workflow | Throwaway prototypes listed as TDD exception | Supported — use design to throw away code, create new stories with learnings |
| **Learning integration** | Retrospectives after epics | Not persisted | Core workflow — developer logs capture discoveries, bubble to design via `/design-update` |
| **Branching exploration** | Course correction for mid-flight changes | Plan consumed linearly | Add exploratory stories, reorder via dependencies, keep or discard branches |

**Blueprint's approach:** Reduce planning to favor discovery. Learn through implementation, capture discoveries in developer logs, flow learnings back into the design. The design becomes a living document that reflects what you've learned, not just what you planned. This supports:

- **Prototyping:** Build to learn, throw away code, start fresh with knowledge
- **Exploratory branching:** Try different approaches in parallel stories, keep what works
- **Greenfield uncertainty:** Start with minimal design, let implementation reveal the real requirements

**Trade-off:** More planning overhead provides more guidance and catches issues earlier. Less overhead enables faster iteration and discovery but requires developer judgment about when to add structure.

## Summary

| Aspect | BMAD | Superpowers | Blueprint |
|--------|------|-------------|-----------|
| **Ideation/planning** | Most structured — full PRD, architecture, agent personas | Mandatory brainstorming gate | Lean — developer-guided |
| **Execution discipline** | TDD integrated, adversarial reviews | Strictest TDD, subagent isolation | Flexible — TDD optional, parallel reviewers |
| **Learning capture** | Dev record + retrospectives | Session only (not persisted) | Most comprehensive — developer logs per story, bubbles to design |
| **Mid-flight adaptation** | Course correction workflow | Plan consumed by execution | Most flexible — add/remove/reorder stories |
| **Workflow flexibility** | Full method + Quick Dev | Single workflow, mandatory gates | Most modes — design-driven, story-first, ad-hoc |
| **Customization** | 3-layer TOML merge | User preferences override | Custom templates + derivable |
| **Dependency tracking** | Implicit (epic/story ordering) | Implicit (sequential tasks) | Explicit (`blocked-by`, cross-design) |
| **Document reviews** | Built-in checklists | Built-in self-review | Customizable per project/template |
| **Code reviews** | Adversarial (Blind Hunter, no context) | Single comprehensive reviewer | Parallel specialized + false positive validation |

## Blueprint Gaps and Trade-offs

Blueprint optimizes for leanness and flexibility. This means some features in other frameworks are intentionally absent or simplified.

### Acceptance Criteria Traceability

| Framework | AC handling |
|-----------|-------------|
| **BMAD** | Tasks link to specific ACs: `- [ ] Task 1 (AC: #)` — full traceability from requirement to implementation |
| **Superpowers** | Plans reference specs, reviewer checks plan alignment |
| **Blueprint** | Designs have Requirements/Constraints sections, but story tasks don't link back to specific requirements |

Blueprint's lean stories trade traceability for simplicity. For regulated domains or large teams needing audit trails, BMAD's explicit AC linking may be more appropriate.

### Cross-Artifact Consistency Checks

BMAD has `bmad-check-implementation-readiness` — verifies PRD, architecture, epics, and stories are aligned before development begins. Blueprint relies on individual document reviews (`/design-review`, `/story-review`) but has no cross-artifact consistency check.

**Trade-off:** Blueprint's leanness may not require comprehensive cross-checks — fewer artifacts means less drift. But for complex multi-design projects, a consistency check could catch misalignment.

### Specialized Planning Workflows

| Workflow | BMAD | Superpowers | Blueprint |
|----------|------|-------------|-----------|
| **UX Design** | Dedicated agent + workflow | Part of brainstorming | Not specialized |
| **Architecture** | Dedicated agent + workflow | Part of design doc | Part of design (Technical Approach) |
| **PRD/Requirements** | Full facilitated workflow | Part of brainstorming | Design doc (Requirements section) |

Blueprint embeds these concerns in the design document rather than providing specialized workflows. **This is intentional** — Blueprint focuses on execution and learning capture, and integrates with other frameworks for ideation and planning. Use BMAD for PRD generation and brainstorming, then feed the output into Blueprint designs (see [Combining Workflows](#combining-workflows)).

### Multi-Agent Collaboration

BMAD's Party Mode brings multiple personas (PM, Architect, Dev, UX) into one conversation for:
- Big decisions with tradeoffs
- Brainstorming sessions
- Post-mortems
- Sprint retrospectives

Blueprint uses single-agent with context switching. Users can prompt for different perspectives, but there's no structured multi-persona workflow.

### Context Isolation

| Framework | Approach | Trade-off |
|-----------|----------|-----------|
| **Superpowers** | Fresh subagent per task — clean context, no bleed-through | May lose helpful context between tasks |
| **Blueprint** | Shared context by default, selective forking in skills | Coherence across tasks, but context can grow large |

Blueprint's approach favors coherence — learnings from one task inform the next. Superpowers' isolation prevents context pollution but requires more explicit context passing.

## Verdict

**BMAD** is a full agile methodology. Choose it when you want structured guidance from ideation through retrospectives, multiple agent personas, and a proven process. Best for teams wanting a complete framework or developers new to AI-assisted development.

**Superpowers** is disciplined execution. Choose it when you want strict TDD enforcement, subagent isolation per task, and a philosophy that treats skipping process as rationalization. Best for developers who value test-first discipline and don't want escape hatches.

**Blueprint** is lean and adaptive. Choose it when you want lightweight specs, explicit dependency tracking, comprehensive learning capture, and the flexibility to choose your entry point based on what you know. Best for developers who have opinions about their workflow and want to capture decisions as they go.

All three are spec-anchored. They complement each other — BMAD's ideation can feed Blueprint's designs, and all three work for greenfield and brownfield projects.
