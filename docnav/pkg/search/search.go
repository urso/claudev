package search

import (
	"iter"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/urso/claudev/docnav/pkg/document"
	"github.com/urso/claudev/docnav/pkg/parser"
	"github.com/urso/claudev/docnav/pkg/text"
)

type Result struct {
	Path         string  `json:"path"`
	Title        string  `json:"title"`
	Score        float64 `json:"score"`
	MatchedField string  `json:"matched_field"`
	IsHub        bool    `json:"is_hub"`
}

type scoredField struct {
	name   string
	weight float64
	terms  []string
}

type candidate struct {
	path     string
	title    string
	isHub    bool
	fields   []scoredField
	allTerms []string
}

// Searcher performs document search with tokenization, stemming, and TF-IDF scoring.
type Searcher struct {
	tokenizer  *text.Tokenizer
	typeFilter string
	tagsFilter []string
}

type Option func(*Searcher)

// WithTypeFilter filters results to documents matching the given type.
func WithTypeFilter(t string) Option {
	return func(s *Searcher) { s.typeFilter = t }
}

// WithTagsFilter filters results to documents containing all given tags.
func WithTagsFilter(tags ...string) Option {
	return func(s *Searcher) { s.tagsFilter = tags }
}

// NewSearcher creates a Searcher with a default tokenizer and the given options.
func NewSearcher(opts ...Option) *Searcher {
	s := &Searcher{tokenizer: text.NewTokenizer()}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Search finds documents matching query in the given file iterator.
func (s *Searcher) Search(files iter.Seq2[string, error], query string) ([]Result, error) {
	queryTerms := s.tokenizer.Tokenize(query)
	if len(queryTerms) == 0 {
		return nil, nil
	}
	querySet := make(map[string]struct{}, len(queryTerms))
	for _, t := range queryTerms {
		querySet[t] = struct{}{}
	}

	// Pass 1: parallel read/parse/tokenize, pre-filter candidates.
	// Workers send candidates via channel to a collector that builds the corpus.
	numWorkers := runtime.GOMAXPROCS(0)
	if numWorkers < 1 {
		numWorkers = 1
	}

	paths := make(chan string)
	found := make(chan candidate)
	var wg sync.WaitGroup

	for range numWorkers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range paths {
				if c, ok := s.buildCandidate(path, querySet); ok {
					found <- c
				}
			}
		}()
	}

	// Collector: accumulates candidates and builds corpus as they arrive.
	var candidates []candidate
	corpus := text.NewCorpus()
	collectorDone := make(chan struct{})
	go func() {
		for c := range found {
			corpus.Add(c.allTerms)
			candidates = append(candidates, c)
		}
		close(collectorDone)
	}()

	for path, err := range files {
		if err != nil {
			slog.Warn("walk error", "error", err)
			continue
		}
		paths <- path
	}
	close(paths)
	wg.Wait()
	close(found)
	<-collectorDone

	if len(candidates) == 0 {
		return nil, nil
	}

	results := make([]Result, 0, len(candidates))
	for _, c := range candidates {
		score, field := scoreCandidate(c, queryTerms, corpus)
		if score == 0 {
			continue
		}
		if c.isHub {
			score *= 1.5
		}
		results = append(results, Result{
			Path:         c.path,
			Title:        c.title,
			Score:        score,
			MatchedField: field,
			IsHub:        c.isHub,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}

func (s *Searcher) buildCandidate(path string, querySet map[string]struct{}) (candidate, bool) {
	content, err := os.ReadFile(path)
	if err != nil {
		return candidate{}, false
	}

	doc := parser.ParseDocument(path, content)
	if !s.matchFilters(doc) {
		return candidate{}, false
	}

	// Tokenize each field.
	filename := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	fields := []scoredField{
		{name: "filename", weight: 100, terms: s.tokenizer.Tokenize(filename)},
		{name: "title", weight: 80, terms: s.tokenizer.Tokenize(doc.Frontmatter.Title)},
		{name: "tags", weight: 60, terms: s.tokenizer.Tokenize(strings.Join(doc.Frontmatter.Tags, " "))},
		{name: "scope", weight: 50, terms: s.tokenizer.Tokenize(doc.Frontmatter.Scope)},
		{name: "headings", weight: 40, terms: s.tokenizer.Tokenize(strings.Join(doc.Headings, " "))},
		{name: "content", weight: 10, terms: s.tokenizer.Tokenize(string(content))},
	}

	// Collect all terms and check if any query term matches.
	var allTerms []string
	hasMatch := false
	for _, f := range fields {
		allTerms = append(allTerms, f.terms...)
		if !hasMatch {
			for _, t := range f.terms {
				if _, ok := querySet[t]; ok {
					hasMatch = true
					break
				}
			}
		}
	}
	if !hasMatch {
		return candidate{}, false
	}

	title := doc.Frontmatter.Title
	if title == "" {
		title = filepath.Base(path)
	}

	return candidate{
		path:     path,
		title:    title,
		isHub:    doc.IsHub,
		fields:   fields,
		allTerms: allTerms,
	}, true
}

func scoreCandidate(c candidate, queryTerms []string, corpus *text.Corpus) (float64, string) {
	var bestScore float64
	var bestField string

	for _, f := range c.fields {
		if len(f.terms) == 0 {
			continue
		}
		tf := text.TF(f.terms)
		tfidf := text.TFIDF(tf, corpus)

		// Sum TF-IDF scores for query terms present in this field.
		var fieldScore float64
		for _, qt := range queryTerms {
			fieldScore += tfidf[qt]
		}
		fieldScore *= f.weight

		if fieldScore > bestScore {
			bestScore = fieldScore
			bestField = f.name
		}
	}

	return bestScore, bestField
}

func (s *Searcher) matchFilters(doc document.Document) bool {
	if s.typeFilter != "" && !strings.EqualFold(doc.Frontmatter.Type, s.typeFilter) {
		return false
	}
	if len(s.tagsFilter) > 0 {
		tagSet := make(map[string]bool, len(doc.Frontmatter.Tags))
		for _, t := range doc.Frontmatter.Tags {
			tagSet[strings.ToLower(t)] = true
		}
		for _, required := range s.tagsFilter {
			if !tagSet[strings.ToLower(required)] {
				return false
			}
		}
	}
	return true
}
