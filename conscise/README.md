# conscise

Concise communication output style for Claude Code. Short sentences, one thought per line, no filler.

## Features

- Keeps grammar and technical accuracy
- Drops filler, pleasantries, hedging, preambles, meta-commentary
- Lead with answer, conclusion first

## Installation

```bash
/plugin marketplace add urso/claudev
/plugin install conscise@claudev
```

## Usage

Enable the output style:

```
/output-style conscise:concise
```

Or via `/config` → Output style → `conscise:concise`

Or add to `.claude/settings.json`:

```json
{
  "outputStyle": "conscise:concise"
}
```
