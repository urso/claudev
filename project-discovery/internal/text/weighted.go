package text

import "iter"

// WeightedFields combines per-field TermFreqs into a single weighted score.
// Use it when a document has multiple text regions with different importance
// (title, body, tags, ...) that should contribute additively to scoring
// against a query.
//
// The combined TF is not a probability distribution (it sums to more than 1
// across fields) and is intentionally not exposed as a TFIDFVector — it has
// no meaningful cosine. It is a scoring primitive only.
type WeightedFields struct {
	tf map[string]float64
}

// NewWeightedFields creates an empty WeightedFields.
func NewWeightedFields() *WeightedFields {
	return &WeightedFields{tf: make(map[string]float64)}
}

// Add accumulates weight*tf into the combined vector and returns w for chaining.
func (w *WeightedFields) Add(weight float64, tf TermFreq) *WeightedFields {
	for term, freq := range tf.tf {
		w.tf[term] += weight * freq
	}
	return w
}

// Score sums weight*tf*idf over query terms against corpus.
// Terms with zero IDF or absent from the combined vector contribute 0.
func (w *WeightedFields) Score(query iter.Seq[string], corpus *Corpus) float64 {
	var total float64
	for term := range query {
		tf := w.tf[term]
		if tf == 0 {
			continue
		}
		idf := corpus.IDF(term)
		if idf == 0 {
			continue
		}
		total += tf * idf
	}
	return total
}
