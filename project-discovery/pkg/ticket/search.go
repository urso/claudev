package ticket

import (
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/urso/claudev/project-discovery/internal/document"
	"github.com/urso/claudev/project-discovery/internal/text"
)

// Per-field weights. Title hits dominate; body is background signal.
const (
	weightTitle     = 80.0
	weightTags      = 60.0
	weightScope     = 50.0
	weightIntention = 40.0
	weightBody      = 10.0
)

// SearchFilter narrows candidates before scoring. Empty fields match everything.
// Tags matches exact frontmatter tag values (no stemming — tags are identifiers).
// Scope matches by stemmed-token overlap against the ticket's scope.
type SearchFilter struct {
	Tags  []string
	Scope string
}

// SearchHit is one ranked search result.
type SearchHit struct {
	Summary
	Score float64
}

type indexedTicket struct {
	summary Summary
	score   *text.WeightedFields
}

// Search ranks tickets in folder by relevance to query. Filter is applied first;
// remaining tickets form the corpus. Tickets with zero score are dropped.
// Results are sorted by score desc, then id asc.
func Search(folder TicketFolder, query string, filter SearchFilter) ([]SearchHit, error) {
	tok := text.NewTokenizer()
	queryTokens := slices.Collect(tok.Tokenize(query))
	if len(queryTokens) == 0 {
		return nil, nil
	}
	filterScopeStems := slices.Collect(tok.Tokenize(filter.Scope))

	matches, err := filepath.Glob(filepath.Join(folder.Path(), IDPrefix+"*.md"))
	if err != nil {
		return nil, err
	}

	docs := make([]indexedTicket, 0, len(matches))
	corpus := text.NewCorpus()

	for _, path := range matches {
		doc, err := document.ParseDocumentFile[Ticket](path)
		if err != nil {
			return nil, err
		}
		fm := doc.Frontmatter

		if !matchTagsFilter(fm.Tags, filter.Tags) {
			continue
		}
		if !matchScopeFilter(tok, fm.Scope, filterScopeStems) {
			continue
		}

		titleTokens := slices.Collect(tok.Tokenize(fm.Title))
		tagsTokens := slices.Collect(tok.Tokenize(strings.Join(fm.Tags, " ")))
		scopeTokens := slices.Collect(tok.Tokenize(fm.Scope))
		intentTokens := slices.Collect(tok.Tokenize(fm.Intention))
		bodyTokens := slices.Collect(tok.Tokenize(string(doc.Body)))

		score := text.NewWeightedFields().
			Add(weightTitle, text.NewTermFreq(slices.Values(titleTokens))).
			Add(weightTags, text.NewTermFreq(slices.Values(tagsTokens))).
			Add(weightScope, text.NewTermFreq(slices.Values(scopeTokens))).
			Add(weightIntention, text.NewTermFreq(slices.Values(intentTokens))).
			Add(weightBody, text.NewTermFreq(slices.Values(bodyTokens)))

		corpus.Add(concatSeq(titleTokens, tagsTokens, scopeTokens, intentTokens, bodyTokens))

		docs = append(docs, indexedTicket{
			summary: Summary{
				ID:     fm.ID,
				Status: fm.Status,
				Title:  fm.Title,
				Tags:   fm.Tags,
				Path:   path,
			},
			score: score,
		})
	}

	hits := make([]SearchHit, 0, len(docs))
	for _, d := range docs {
		score := d.score.Score(slices.Values(queryTokens), corpus)
		if score == 0 {
			continue
		}
		hits = append(hits, SearchHit{Summary: d.summary, Score: score})
	}
	sort.Slice(hits, func(i, j int) bool {
		if hits[i].Score != hits[j].Score {
			return hits[i].Score > hits[j].Score
		}
		return hits[i].ID < hits[j].ID
	})
	return hits, nil
}

func matchTagsFilter(have, want []string) bool {
	for _, w := range want {
		if !slices.Contains(have, w) {
			return false
		}
	}
	return true
}

func matchScopeFilter(tok *text.Tokenizer, ticketScope string, wantStems []string) bool {
	if len(wantStems) == 0 {
		return true
	}
	for s := range tok.Tokenize(ticketScope) {
		if slices.Contains(wantStems, s) {
			return true
		}
	}
	return false
}
