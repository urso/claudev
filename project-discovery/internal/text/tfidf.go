package text

import (
	"iter"
	"math"
)

// Corpus accumulates document frequency stats for IDF computation.
type Corpus struct {
	docCount int
	docFreq  map[string]int
}

// NewCorpus creates an empty Corpus.
func NewCorpus() *Corpus {
	return &Corpus{docFreq: make(map[string]int)}
}

// Add registers one document's terms. Only unique terms per document are counted.
func (c *Corpus) Add(terms iter.Seq[string]) {
	c.docCount++
	seen := make(map[string]struct{})
	for t := range terms {
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		c.docFreq[t]++
	}
}

// IDF returns the inverse document frequency for a term.
func (c *Corpus) IDF(term string) float64 {
	if c.docCount == 0 {
		return 0
	}
	return math.Log(1 + float64(c.docCount)/(1+float64(c.docFreq[term])))
}

// TermFreq holds normalized term-frequency counts for a single document.
// Comparable across documents only after weighting via Weight; raw TermFreq
// values from different documents are not meaningfully comparable.
type TermFreq struct {
	tf map[string]float64
}

// NewTermFreq consumes terms once and returns a TermFreq with normalized
// frequencies (count / total).
func NewTermFreq(terms iter.Seq[string]) TermFreq {
	counts := make(map[string]int)
	total := 0
	for t := range terms {
		counts[t]++
		total++
	}
	if total == 0 {
		return TermFreq{}
	}
	tf := make(map[string]float64, len(counts))
	denom := float64(total)
	for t, n := range counts {
		tf[t] = float64(n) / denom
	}
	return TermFreq{tf: tf}
}

// Get returns the normalized term frequency for term, or 0 if absent.
func (t TermFreq) Get(term string) float64 {
	return t.tf[term]
}

// Empty reports whether this TermFreq has no terms.
func (t TermFreq) Empty() bool { return len(t.tf) == 0 }

// Weight applies IDF from corpus to produce a TFIDFVector suitable for
// cosine comparison.
func (t TermFreq) Weight(corpus *Corpus) TFIDFVector {
	if len(t.tf) == 0 {
		return TFIDFVector{}
	}
	v := make(map[string]float64, len(t.tf))
	for term, freq := range t.tf {
		v[term] = freq * corpus.IDF(term)
	}
	return TFIDFVector{v: v}
}
