---
description: Analyze requirements or existing doc against codebase to surface concepts, risks, and gaps
user-invocable: true
disable-model-invocation: true
argument-hint: "[design|story|requirements file]"
allowed-tools: ["Read", "Bash", "Glob", "Grep", "Agent", "AskUserQuestion"]
---

# Analyze

Explore codebase against a target (design, story, or requirements) to surface concepts, risks, and gaps. Runs in a fresh-context sub-agent to avoid bias from conversation history.

## Variables

- **FIND_DOC**: `find-doc.sh`

## User Input
```
$ARGUMENTS
```

Parse for: design ID/name, story ID/name, or file path.

## Process

### 1. Resolve Target

Determine target type and load content:

- **Design**: Use `bash FIND_DOC design "<input>"` to locate, then read
- **Story**: Use `bash FIND_DOC story "<input>"` to locate, then read
- **File path**: Read directly (requirements doc, RFC, etc.)
- **No input**: Ask user what to analyze

Extract the key information:
- For designs: Problem Context, Goals, Requirements, Constraints
- For stories: Tasks, Technical Notes
- For requirements files: Full content

### 2. Spawn Analysis Agent

Spawn a sub-agent with fresh context to avoid bias:

```
subagent_type: general-purpose

Analyze the following against the codebase. Explore relevant code, don't assume.

## Target
<extracted content from step 1>

## Analysis Tasks

1. **Concept Identification**
   - What domain concepts does this involve?
   - Which exist in the codebase (tables, types, services)?
   - Which are new?

2. **Codebase Exploration**
   - Search for related code by concept names
   - Identify existing patterns and conventions
   - Note relevant files and their purpose

3. **Risk & Gap Analysis**
   - What's ambiguous or underspecified?
   - What edge cases aren't covered?
   - What technical constraints exist?
   - What assumptions need validation?

4. **Integration Points**
   - What existing code will this touch?
   - What APIs, services, or data flows are involved?
   - What might break?

## Output Format

### Existing Concepts
- [Concept]: [location] — [relevance]

### New Concepts Required
- [Concept]: [purpose] — [relationship to existing]

### Codebase Patterns
- [Pattern]: [where used] — [apply here?]

### Risks & Gaps
- [Risk/Gap]: [impact] — [mitigation or question to resolve]

### Integration Points
- [Component]: [how affected]

Be specific. Reference actual files and code. No assumptions.
```

### 3. Present Findings

Display the sub-agent's analysis to the user.

### 4. Suggest Next Steps

Based on target type:

- **Design**: "Update design with findings? `/design-review` when ready."
- **Story**: "Update story's Technical Notes? Create new tasks for gaps?"
- **Requirements**: "Ready to create a design? `/design-create`"

Offer to apply relevant findings to the target document.

## Guidelines

- Sub-agent runs with fresh context — no conversation history bias
- Analysis is read-only — doesn't modify code
- Findings update the target doc only with user approval
- Works at any certainty level: vague requirements or detailed specs
