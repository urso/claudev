package text

import "math"

// TFIDFVector is an IDF-weighted document vector. Two TFIDFVectors derived
// from the same Corpus may be compared by Cosine.
type TFIDFVector struct {
	v map[string]float64
}

// Empty reports whether this vector has no terms.
func (a TFIDFVector) Empty() bool { return len(a.v) == 0 }

// Cosine returns the cosine similarity between a and b.
// Returns 0 if either vector is empty.
func (a TFIDFVector) Cosine(b TFIDFVector) float64 {
	if len(a.v) == 0 || len(b.v) == 0 {
		return 0
	}

	var dot, normA, normB float64
	for term, va := range a.v {
		normA += va * va
		if vb, ok := b.v[term]; ok {
			dot += va * vb
		}
	}
	for _, vb := range b.v {
		normB += vb * vb
	}

	denom := math.Sqrt(normA) * math.Sqrt(normB)
	if denom == 0 {
		return 0
	}
	return dot / denom
}
