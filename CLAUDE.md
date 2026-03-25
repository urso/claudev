# claudev - Plugin Marketplace

This repository is a curated Claude Code plugin marketplace. It serves as a registry of third-party plugins that can be installed via:

```
/plugin marketplace add urso/claudev
/plugin install <plugin-name>@claudev
```

## Structure

- `.claude-plugin/marketplace.json` - The marketplace catalog.

## Adding Plugins

Add new entries to `.claude-plugin/marketplace.json`.

### Local plugins (in this repo)

Use a relative path string as the source:
```json
{ "name": "my-plugin", "source": "./my-plugin", ... }
```

### External plugins (from another repo)

Check the upstream repo's `.claude-plugin/marketplace.json` to understand how the plugin is structured.

1. If the plugin source is `"./"` (repo root), use `"source": "github"` with the repo.
2. If the plugin source is a subdirectory (e.g. `"./plugins/foo"`), use `"source": "git-subdir"` with the path.
3. Do NOT pin versions in our marketplace — we want users to always get the latest from upstream.
