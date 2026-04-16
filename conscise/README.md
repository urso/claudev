# conscise

Concise communication mode for Claude Code. Short sentences, one thought per line, no filler.

## Features

- Auto-activates on session start via SessionStart hook
- Toggle with `/conscise off` or `/conscise on`
- Keeps grammar and technical accuracy
- Drops filler, pleasantries, hedging, preambles, meta-commentary

## Installation

```bash
/plugin marketplace add urso/claudev
/plugin install conscise@claudev
```

## Usage

Once installed, conscise mode activates automatically at session start.

To disable mid-session:
```
/conscise off
```

To re-enable:
```
/conscise on
```
