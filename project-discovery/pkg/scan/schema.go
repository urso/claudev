package scan

import (
	"github.com/invopop/jsonschema"
)

// jqRecipes is a hand-curated map of useful jq expressions over scan output.
// Generated schemas describe shape but not query patterns; recipes ride along
// to give agents a starting point.
var jqRecipes = map[string]string{
	"top_5_languages_by_files":  ".repo.extensions[:5]",
	"all_markers":               ".repo.markers",
	"top_files_in_dir":          `.churn.top_files | map(select(.path | startswith("pkg/auth/")))`,
	"current_owners_of_dir":     `.churn.dir_ownership[] | select(.dir == "pkg") | .authors | map({name, email, share, lines_recent})`,
	"active_maintainers":        ".churn.contributors | map(select(.lines_recent > 0)) | sort_by(-.lines_recent)",
	"departed_top_contributors": ".churn.contributors | map(select(.lines_recent == 0 and .lines > 0))",
}

// Schema returns a JSON Schema (draft 2020-12) describing the shape of
// `scan repo` and `scan churn` outputs, plus a `jq_recipes` companion map
// of common queries. Generated from struct tags via invopop/jsonschema —
// updating a struct field automatically updates the schema.
func Schema() map[string]any {
	r := jsonschema.Reflector{
		ExpandedStruct:            true,
		AllowAdditionalProperties: true,
	}
	return map[string]any{
		"description": "Schema for `discover scan repo` and `discover scan churn` outputs. Use jq to slice; see jq_recipes for examples.",
		"repo":        r.Reflect(&RepoSignals{}),
		"churn":       r.Reflect(&ChurnSignals{}),
		"jq_recipes":  jqRecipes,
	}
}
