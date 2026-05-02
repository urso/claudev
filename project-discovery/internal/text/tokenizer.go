package text

import (
	"iter"
	"strings"
	"unicode"

	"github.com/kljensen/snowball"
)

// Tokenizer splits text into stemmed, lowercase tokens with stopword removal.
type Tokenizer struct {
	stopwords map[string]struct{}
}

// NewTokenizer creates a Tokenizer with default English stopwords.
func NewTokenizer() *Tokenizer {
	return &Tokenizer{stopwords: defaultStopwordsSet}
}

var defaultStopwordsSet = func() map[string]struct{} {
	sw := make(map[string]struct{}, len(defaultStopwords))
	for _, w := range defaultStopwords {
		sw[w] = struct{}{}
	}
	return sw
}()

// Tokenize yields stemmed, stopword-filtered, lowercase tokens from s.
// Returns an iter.Seq so callers can stream straight into TermFreq/Corpus
// without an intermediate slice; use slices.Collect when materialization is
// needed for fan-out.
func (t *Tokenizer) Tokenize(s string) iter.Seq[string] {
	return func(yield func(string) bool) {
		for w := range splitWords(s) {
			w = strings.ToLower(w)
			if _, stop := t.stopwords[w]; stop {
				continue
			}
			stemmed, err := snowball.Stem(w, "english", false)
			if err != nil {
				stemmed = w
			}
			if stemmed == "" {
				continue
			}
			if !yield(stemmed) {
				return
			}
		}
	}
}

// splitWords yields runs of letter/digit runes from s as substrings, without
// allocating a slice. Equivalent to strings.FieldsFunc with the same predicate.
func splitWords(s string) iter.Seq[string] {
	return func(yield func(string) bool) {
		start := -1
		for i, r := range s {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				if start < 0 {
					start = i
				}
				continue
			}
			if start >= 0 {
				if !yield(s[start:i]) {
					return
				}
				start = -1
			}
		}
		if start >= 0 {
			yield(s[start:])
		}
	}
}

var defaultStopwords = []string{
	"a", "an", "and", "are", "as", "at", "be", "but", "by",
	"for", "from", "has", "have", "he", "her", "his",
	"if", "in", "into", "is", "it", "its",
	"no", "not", "of", "on", "or",
	"she", "so", "than", "that", "the", "their", "them", "then",
	"there", "these", "they", "this", "to",
	"was", "we", "were", "what", "when", "which", "who", "will", "with",
	"you", "your",
}
