package text

import "math"

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
func (c *Corpus) Add(terms []string) {
	c.docCount++
	seen := make(map[string]struct{}, len(terms))
	for _, t := range terms {
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

// TF computes term frequency for a slice of tokens.
func TF(terms []string) map[string]float64 {
	if len(terms) == 0 {
		return nil
	}
	counts := make(map[string]int, len(terms))
	for _, t := range terms {
		counts[t]++
	}
	total := float64(len(terms))
	tf := make(map[string]float64, len(counts))
	for t, n := range counts {
		tf[t] = float64(n) / total
	}
	return tf
}

// TFIDF scores terms against a corpus, returning a score per term.
func TFIDF(tf map[string]float64, corpus *Corpus) map[string]float64 {
	scores := make(map[string]float64, len(tf))
	for term, freq := range tf {
		scores[term] = freq * corpus.IDF(term)
	}
	return scores
}
