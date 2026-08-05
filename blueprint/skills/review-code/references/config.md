# Configuration Review Instructions

Review configuration files for correctness, security, and embedded code defects.

Covers Kubernetes manifests, Helm charts, CI workflows, Terraform, and general
config formats (YAML, JSON, TOML, INI).

## Embedded Code

Config files routinely embed executable content — that's where the serious bugs
live.

Review embedded code **in place**. The surrounding manifest determines whether a
bug is reachable or harmless; you need both in view.

Look for:

- `command:` / `args:` arrays — especially `["/bin/sh", "-c", "<script>"]`
- ConfigMap/Secret `data:` values holding scripts, nginx configs, Lua, SQL
- Lifecycle hooks, probe handlers, init container scripts
- Job/CronJob commands
- `_helpers.tpl` and `*.tpl` files
- CI `run:` blocks
- Any value with a shebang

For each piece of embedded code, apply
`${CLAUDE_PLUGIN_ROOT}/skills/review-code/references/bugs.md`. Embedded shell is
code — review it as code.

### Interpolation = Injection

Template values substituted into shell strings are injection surfaces:

```yaml
command:
  - /bin/sh
  - -c
  - psql -c "SELECT * FROM {{ .Values.tableName }}"   # injection via values
```

Values come from `values.yaml`, which chart consumers control. Treat as
untrusted input.

## Severity

- **[error]** — security exposure, credential leak, config that fails to
  render/deploy, definite bug in embedded code
- **[warning]** — missing limits/probes, unpinned versions, type coercion risk,
  dead config, context-dependent issues

**Do NOT report:**

- Pre-existing issues in unchanged files passed as context
- Style preferences not backed by a documented rule
- Speculation without a concrete failure path

## Report Format

For embedded code findings, use the formats in `bugs.md`.

For structural findings:

```
[warning] path/to/file.yaml:LINE

**Issue:** What is wrong.
**Context:** What this config controls; when it takes effect.
**Impact:** What happens at deploy or runtime.
**Fix:** How to address it.
```

For trivial mechanical issues:

```
[error] path/to/file.yaml:LINE
Issue: Description.
Impact: Consequence.
```
