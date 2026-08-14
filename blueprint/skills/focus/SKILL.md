---
description: Set worktree focus to a design or story for streamlined development sessions
user-invocable: true
disable-model-invocation: true
argument-hint: "[design|story|clear] [search terms]"
allowed-tools: ["Read", "Write", "Bash", "AskUserQuestion"]
---

# Focus

Set worktree focus to a design or story. Focus persists in `.claude-focus` at the git root, so subsequent `/work` calls know what you're working on without re-specifying.

## Variables

- **QUERY_WORK**: `query-work.sh`
- **FIND_DOC**: `find-doc.sh`

## User Input
```
$ARGUMENTS
```

## Process

### 1. Find Git Root

```bash
GIT_ROOT=$(git rev-parse --show-toplevel)
```

The focus file lives at `$GIT_ROOT/.claude-focus`.

### 2. Parse Input

Interpret the user's input conversationally:

- **No input** → Show current focus, or if none set, ask what they want to work on
- **`clear`** → Remove `.claude-focus` file, confirm focus cleared
- **`design`** → List/search designs, help user pick one
- **`story`** → List/search stories (within focused design if set, otherwise all), help user pick one
- **`design <terms>`** → Search designs matching terms
- **`story <terms>`** → Search stories matching terms
- **`<terms>`** → Search both designs and stories, discuss matches

### 3. Search and Select

Use the query-work script to find matching documents:

```bash
# List all actionable work
bash QUERY_WORK --actionable

# Search with terms
bash QUERY_WORK --search "<terms>"

# Find specific document
bash FIND_DOC "<name or id>"
```

Output format for query-work: `TYPE | ID | FILENAME | TITLE | STATUS | BLOCKED-BY | DESCRIPTION`

**Dialogue patterns:**

- If search returns one match → confirm and set focus
- If search returns multiple → present options, let user pick
- If no matches → suggest alternatives or ask for clarification
- If user asks to explore → list available designs/stories with descriptions
- If story is blocked → inform user, ask if they want to focus anyway

### 4. Determine Document Type

When a document is selected, read its frontmatter to determine type:

- Has `design:` field → it's a story (field references parent design)
- No `design:` field → it's a design

### 5. Write Focus File

Write `.claude-focus` as YAML at the git root.

**IMPORTANT**: Use the **exact file path** returned by FIND_DOC or shown in QUERY_WORK output (the FILENAME column). Do NOT invent or simplify paths. The path must match the actual file location (e.g., `docs/ai/designs/0001-greeter.md`, not `docs/designs/...`).

**Focusing on a design:**
```yaml
design: <exact-path-from-find-doc>
```

Clears any previous story focus. User will pick a story on next `/work`.

**Focusing on a story with parent design:**
```yaml
design: <exact-path-to-parent-design>
story: <exact-path-to-story>
```

Resolve the parent design from the story's `design:` frontmatter field using FIND_DOC.

**Focusing on a standalone story (no design field):**
```yaml
story: <exact-path-to-story>
```

### 6. Confirm

Tell the user what focus was set:

- For design: mention how many actionable stories are available
- For story: mention the story title and status
- Suggest running `/work` to start development

## Examples

**Show current focus:**
```
> /focus
Currently focused on:
  Design: 0002 - Dreams
  Story: 0008 - CLI Foundation (in-progress)

Run /work to continue development.
```

**Set design focus:**
```
> /focus design dreams
Setting focus to design 0002 - Dreams.
3 actionable stories. Run /work to pick one.
```

**Set story focus:**
```
> /focus story 0008
Setting focus to story 0008 - CLI Foundation.
Parent design: 0002 - Dreams
Run /work to continue.
```

**Explore options:**
```
> /focus
No focus set. What would you like to work on?

> something about docs
Found:
1. Design 0001 - docnav (in-progress) - Local docs navigation
2. Story 0001-0004 - Tags and Stale (ready)

Which one?
```

**Clear focus:**
```
> /focus clear
Focus cleared.
```
