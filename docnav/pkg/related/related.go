package related

import (
	"sort"

	"github.com/urso/claudev/docnav/pkg/document"
	"github.com/urso/claudev/docnav/pkg/linkgraph"
	"github.com/urso/claudev/docnav/pkg/text"
)

// Result holds the relatedness score for a candidate document.
type Result struct {
	Path          string  `json:"path"`
	Score         float64 `json:"score"`
	CoCitation    int     `json:"co_citation"`
	SharedTags    int     `json:"shared_tags"`
	SharedWatches int     `json:"shared_watches"`
	ContentSim    float64 `json:"content_similarity"`
}

// Related finds documents related to target using reciprocal rank fusion of
// two independent scorers: structural (co-citation, shared tags, overlapping
// watches) and content (TF-IDF cosine similarity). Direct links are excluded.
func Related(target document.Document, graph *linkgraph.Graph, docs []document.Document, contents map[string][]byte, topN int) []Result {
	exclude := buildExcludeSet(target, graph)

	structural := structuralRank(target, graph, docs, exclude)
	content := contentRank(target, docs, contents, exclude)

	return fuse(structural, content, topN)
}

// rankedEntry holds per-candidate detail from one scorer.
type rankedEntry struct {
	path          string
	coCitation    int
	sharedTags    int
	sharedWatches int
	contentSim    float64
}

func buildExcludeSet(target document.Document, graph *linkgraph.Graph) map[string]bool {
	exclude := make(map[string]bool)
	exclude[target.Path] = true
	for _, link := range target.Links {
		exclude[link.Path] = true
	}
	for _, src := range graph.Reverse[target.Path] {
		exclude[src] = true
	}
	return exclude
}

// structuralRank scores by co-citation, shared tags, and overlapping watches,
// returning candidates ordered by total structural score descending.
func structuralRank(target document.Document, graph *linkgraph.Graph, docs []document.Document, exclude map[string]bool) []rankedEntry {
	targetTags := toSet(target.Frontmatter.Tags)
	targetWatches := toSet(target.Frontmatter.Watches)

	// Co-citation: for each link target of the input doc, find other docs
	// that also link there.
	coCitation := make(map[string]int)
	for _, link := range target.Links {
		for _, src := range graph.Reverse[link.Path] {
			if !exclude[src] {
				coCitation[src]++
			}
		}
	}

	type entry struct {
		path                          string
		coCitation, tags, watches, total int
	}
	scores := make(map[string]*entry)

	for path, cc := range coCitation {
		scores[path] = &entry{path: path, coCitation: cc, total: cc}
	}

	for _, doc := range docs {
		if exclude[doc.Path] {
			continue
		}
		tags := countOverlap(targetTags, doc.Frontmatter.Tags)
		watches := countOverlap(targetWatches, doc.Frontmatter.Watches)
		if tags == 0 && watches == 0 && scores[doc.Path] == nil {
			continue
		}
		e, ok := scores[doc.Path]
		if !ok {
			e = &entry{path: doc.Path}
			scores[doc.Path] = e
		}
		e.tags = tags
		e.watches = watches
		e.total += tags + watches
	}

	ranked := make([]rankedEntry, 0, len(scores))
	for _, e := range scores {
		ranked = append(ranked, rankedEntry{
			path:          e.path,
			coCitation:    e.coCitation,
			sharedTags:    e.tags,
			sharedWatches: e.watches,
		})
	}
	sort.Slice(ranked, func(i, j int) bool {
		ti := ranked[i].coCitation + ranked[i].sharedTags + ranked[i].sharedWatches
		tj := ranked[j].coCitation + ranked[j].sharedTags + ranked[j].sharedWatches
		if ti != tj {
			return ti > tj
		}
		return ranked[i].path < ranked[j].path
	})
	return ranked
}

// contentRank scores by TF-IDF cosine similarity of document content,
// returning candidates ordered by similarity descending.
func contentRank(target document.Document, docs []document.Document, contents map[string][]byte, exclude map[string]bool) []rankedEntry {
	tok := text.NewTokenizer()
	corpus := text.NewCorpus()

	// Tokenize all docs and build corpus.
	docTerms := make(map[string][]string, len(docs))
	for _, doc := range docs {
		content := contents[doc.Path]
		terms := tok.Tokenize(string(content))
		docTerms[doc.Path] = terms
		corpus.Add(terms)
	}

	// Compute target TF-IDF vector.
	targetTerms := docTerms[target.Path]
	targetVec := text.TFIDF(text.TF(targetTerms), corpus)

	type scored struct {
		path string
		sim  float64
	}
	var candidates []scored
	for _, doc := range docs {
		if exclude[doc.Path] {
			continue
		}
		vec := text.TFIDF(text.TF(docTerms[doc.Path]), corpus)
		sim := text.Cosine(targetVec, vec)
		if sim > 0 {
			candidates = append(candidates, scored{path: doc.Path, sim: sim})
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].sim != candidates[j].sim {
			return candidates[i].sim > candidates[j].sim
		}
		return candidates[i].path < candidates[j].path
	})

	ranked := make([]rankedEntry, len(candidates))
	for i, c := range candidates {
		ranked[i] = rankedEntry{path: c.path, contentSim: c.sim}
	}
	return ranked
}

// fuse merges two ranked lists using reciprocal rank fusion (k=60).
func fuse(structural, content []rankedEntry, topN int) []Result {
	const k = 60

	// Build detail maps from each scorer.
	structDetail := make(map[string]rankedEntry, len(structural))
	for _, e := range structural {
		structDetail[e.path] = e
	}
	contentDetail := make(map[string]rankedEntry, len(content))
	for _, e := range content {
		contentDetail[e.path] = e
	}

	// RRF: score = sum of 1/(k+rank) across lists. Rank is 1-based.
	rrfScores := make(map[string]float64)
	for rank, e := range structural {
		rrfScores[e.path] += 1.0 / float64(k+rank+1)
	}
	for rank, e := range content {
		rrfScores[e.path] += 1.0 / float64(k+rank+1)
	}

	results := make([]Result, 0, len(rrfScores))
	for path, score := range rrfScores {
		r := Result{Path: path, Score: score}
		if s, ok := structDetail[path]; ok {
			r.CoCitation = s.coCitation
			r.SharedTags = s.sharedTags
			r.SharedWatches = s.sharedWatches
		}
		if c, ok := contentDetail[path]; ok {
			r.ContentSim = c.contentSim
		}
		results = append(results, r)
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return results[i].Path < results[j].Path
	})

	if topN > 0 && len(results) > topN {
		results = results[:topN]
	}
	return results
}

func toSet(items []string) map[string]bool {
	s := make(map[string]bool, len(items))
	for _, item := range items {
		s[item] = true
	}
	return s
}

func countOverlap(set map[string]bool, items []string) int {
	n := 0
	for _, item := range items {
		if set[item] {
			n++
		}
	}
	return n
}
