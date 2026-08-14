---
description: Continue development on the focused design or story
user-invocable: true
disable-model-invocation: true
argument-hint: ""
allowed-tools: ["Read", "Write", "Edit", "Bash", "Glob", "Grep", "AskUserQuestion"]
---

# Work

Continue development on the currently focused design or story. Reads `.claude-focus` and delegates to `develop-story` with the appropriate context.

## Variables

- **QUERY_WORK**: `query-work.sh`
- **FIND_DOC**: `find-doc.sh`
- **DEVELOP_STORY**: `blueprint:develop-story`

## Process

### 1. Read Focus

```bash
GIT_ROOT=$(git rev-parse --show-toplevel)
```

Read `$GIT_ROOT/.claude-focus` if it exists.

If no focus file exists:
- Tell the user: "No focus set. Run `/focus` to set a design or story."
- Stop here.

### 2. Determine Work Target

Parse the focus file (YAML format):

```yaml
design: <path>   # optional
story: <path>    # optional
```

**Cases:**

- **Story is set** → use the story path as target
- **Design only, no story** → will need to select a story (step 3)
- **Neither set** → invalid state, treat as no focus

### 3. Story Selection (if needed)

If focus has design but no story:

1. Run `query-work.sh --actionable --search "<design-id>"` to find actionable stories
2. If one story → auto-select it, update `.claude-focus` with story path
3. If multiple → present options, let user pick, update `.claude-focus`
4. If none → inform user all stories are blocked or done, stop

### 4. Delegate to develop-story

Invoke the develop-story skill with the story path:

```
/develop-story <story-path>
```

This handles:
- Loading design context
- Loading rules
- Task selection
- Planning and implementation
- Story document updates

### 5. Post-Development: Check Story Status

After develop-story completes, read the story document and check its status.

If status changed to `done`:
1. Remove the `story:` line from `.claude-focus` (keep design if present)
2. Inform user: "Story completed. Run `/work` to pick the next story."

This ensures the next `/work` invocation will prompt for story selection from remaining actionable work.

## Examples

**Continue focused story:**
```
> /work
Continuing story 0008 - CLI Foundation...
[delegates to develop-story]
```

**Pick story from focused design:**
```
> /work
Focus: design 0002 - Dreams (no active story)

Actionable stories:
1. 0008 - CLI Foundation (ready)
2. 0009 - Import Command (ready)

Which story? 
> 1
Setting active story to 0008.
[delegates to develop-story]
```

**Story completed:**
```
> /work
[... development completes ...]

Story 0008 - CLI Foundation marked done.
Cleared active story. Run /work to pick the next one.
```

**No focus:**
```
> /work
No focus set. Run /focus to set a design or story.
```
