package ticket

import (
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/urso/claudev/project-discovery/internal/document"
	"github.com/urso/claudev/project-discovery/internal/text"
)

// rrfK is the Reciprocal Rank Fusion smoothing constant. 60 is the canonical
// default from Cormack et al.; lower values amplify top-rank dominance.
const rrfK = 60.0

// OverlapParams describes the new ticket whose existing-ticket overlap we
// want to detect.
type OverlapParams struct {
	Intent string
	Scope  string
	Tags   []string
}

// OverlapHit is one ranked overlap candidate.
type OverlapHit struct {
	Summary
	Score float64
}

// FindOverlap returns existing tickets that may overlap with params, ranked
// by RRF fusion of (a) a structural signal — scope-token overlap + watch-glob
// hits + tag overlap — and (b) a content signal — TF-IDF cosine similarity of
// intent text against ticket text. Tickets with zero combined signal are
// dropped. Sorted by score desc, then id asc.
func FindOverlap(folder TicketFolder, params OverlapParams) ([]OverlapHit, error) {
	tok := text.NewTokenizer()
	intentTokens := slices.Collect(tok.Tokenize(params.Intent))
	paramScopeStems := slices.Collect(tok.Tokenize(params.Scope))

	matches, err := filepath.Glob(filepath.Join(folder.Path(), IDPrefix+"*.md"))
	if err != nil {
		return nil, err
	}

	type candidate struct {
		summary    Summary
		structural float64
		text       []string
	}
	cands := make([]candidate, 0, len(matches))
	corpus := text.NewCorpus()

	for _, path := range matches {
		doc, err := document.ParseDocumentFile[Ticket](path)
		if err != nil {
			return nil, err
		}
		fm := doc.Frontmatter

		scopeStems := slices.Collect(tok.Tokenize(fm.Scope))
		scopeOverlap := stemOverlap(paramScopeStems, scopeStems)
		watchHits := watchGlobHits(fm.Watches, params.Scope)
		tagOverlap := tagOverlapCount(params.Tags, fm.Tags)
		structural := float64(scopeOverlap + watchHits + tagOverlap)

		titleTokens := slices.Collect(tok.Tokenize(fm.Title))
		intentTicketTokens := slices.Collect(tok.Tokenize(fm.Intention))
		bodyTokens := slices.Collect(tok.Tokenize(string(doc.Body)))
		all := make([]string, 0, len(titleTokens)+len(intentTicketTokens)+len(bodyTokens)+len(scopeStems))
		all = append(all, titleTokens...)
		all = append(all, intentTicketTokens...)
		all = append(all, bodyTokens...)
		all = append(all, scopeStems...)
		corpus.Add(slices.Values(all))

		cands = append(cands, candidate{
			summary: Summary{
				ID:     fm.ID,
				Status: fm.Status,
				Title:  fm.Title,
				Tags:   fm.Tags,
				Path:   path,
			},
			structural: structural,
			text:       all,
		})
	}

	queryVec := text.NewTermFreq(slices.Values(intentTokens)).Weight(corpus)

	type ranked struct {
		idx   int
		score float64
	}
	structRanks := make([]ranked, 0, len(cands))
	contentRanks := make([]ranked, 0, len(cands))
	for i, c := range cands {
		if c.structural > 0 {
			structRanks = append(structRanks, ranked{i, c.structural})
		}
		if !queryVec.Empty() {
			docVec := text.NewTermFreq(slices.Values(c.text)).Weight(corpus)
			sim := queryVec.Cosine(docVec)
			if sim > 0 {
				contentRanks = append(contentRanks, ranked{i, sim})
			}
		}
	}
	sort.Slice(structRanks, func(i, j int) bool { return structRanks[i].score > structRanks[j].score })
	sort.Slice(contentRanks, func(i, j int) bool { return contentRanks[i].score > contentRanks[j].score })

	rrf := make(map[int]float64, len(cands))
	for r, item := range structRanks {
		rrf[item.idx] += 1.0 / (rrfK + float64(r+1))
	}
	for r, item := range contentRanks {
		rrf[item.idx] += 1.0 / (rrfK + float64(r+1))
	}

	hits := make([]OverlapHit, 0, len(rrf))
	for idx, score := range rrf {
		hits = append(hits, OverlapHit{Summary: cands[idx].summary, Score: score})
	}
	sort.Slice(hits, func(i, j int) bool {
		if hits[i].Score != hits[j].Score {
			return hits[i].Score > hits[j].Score
		}
		return hits[i].ID < hits[j].ID
	})
	return hits, nil
}

func stemOverlap(a, b []string) int {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	set := make(map[string]struct{}, len(b))
	for _, s := range b {
		set[s] = struct{}{}
	}
	seen := make(map[string]struct{}, len(a))
	n := 0
	for _, s := range a {
		if _, dup := seen[s]; dup {
			continue
		}
		seen[s] = struct{}{}
		if _, ok := set[s]; ok {
			n++
		}
	}
	return n
}

// watchGlobHits returns the number of ticket Watches globs that match the
// param scope string. Scope is matched as both a literal path candidate and
// against each whitespace-separated token, so a multi-word scope like
// "internal/auth login flow" still tests "internal/auth" against globs.
func watchGlobHits(watches []string, scope string) int {
	if len(watches) == 0 || scope == "" {
		return 0
	}
	candidates := append([]string{scope}, strings.Fields(scope)...)
	n := 0
	for _, g := range watches {
		for _, c := range candidates {
			ok, err := filepath.Match(g, c)
			if err == nil && ok {
				n++
				break
			}
		}
	}
	return n
}

func tagOverlapCount(a, b []string) int {
	set := make(map[string]struct{}, len(b))
	for _, t := range b {
		set[t] = struct{}{}
	}
	n := 0
	for _, t := range a {
		if _, ok := set[t]; ok {
			n++
		}
	}
	return n
}
