package text

import (
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
	sw := make(map[string]struct{}, len(defaultStopwords))
	for _, w := range defaultStopwords {
		sw[w] = struct{}{}
	}
	return &Tokenizer{stopwords: sw}
}

// Tokenize splits s into stemmed tokens, removing stopwords.
func (t *Tokenizer) Tokenize(s string) []string {
	words := splitWords(s)
	var tokens []string
	for _, w := range words {
		w = strings.ToLower(w)
		if _, stop := t.stopwords[w]; stop {
			continue
		}
		stemmed, err := snowball.Stem(w, "english", false)
		if err != nil {
			stemmed = w
		}
		if stemmed != "" {
			tokens = append(tokens, stemmed)
		}
	}
	return tokens
}

func splitWords(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
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
