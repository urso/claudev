package text

import "math"

// Cosine returns the cosine similarity between two TF-IDF vectors.
// Returns 0 if either vector is empty.
func Cosine(a, b map[string]float64) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}

	var dot, normA, normB float64
	for term, va := range a {
		normA += va * va
		if vb, ok := b[term]; ok {
			dot += va * vb
		}
	}
	for _, vb := range b {
		normB += vb * vb
	}

	denom := math.Sqrt(normA) * math.Sqrt(normB)
	if denom == 0 {
		return 0
	}
	return dot / denom
}
