package text

import (
	"slices"
	"strings"
	"testing"
)

func TestTokenizer_Basic(t *testing.T) {
	tok := NewTokenizer()
	tokens := slices.Collect(tok.Tokenize("Running deployments quickly"))
	got := strings.Join(tokens, " ")
	if !strings.Contains(got, "run") {
		t.Errorf("expected stemmed 'run', got %q", got)
	}
	if !strings.Contains(got, "deploy") {
		t.Errorf("expected stemmed 'deploy', got %q", got)
	}
}

func TestTokenizer_StopwordsRemoved(t *testing.T) {
	tok := NewTokenizer()
	for tk := range tok.Tokenize("the quick and the dead") {
		if tk == "the" || tk == "and" {
			t.Errorf("stopword %q should have been removed", tk)
		}
	}
}

func TestTokenizer_Empty(t *testing.T) {
	tok := NewTokenizer()
	tokens := slices.Collect(tok.Tokenize(""))
	if len(tokens) != 0 {
		t.Errorf("expected 0 tokens, got %d", len(tokens))
	}
}

func TestTokenizer_Stemming(t *testing.T) {
	tok := NewTokenizer()
	t1 := slices.Collect(tok.Tokenize("architecture"))
	t2 := slices.Collect(tok.Tokenize("architectural"))
	if len(t1) == 0 || len(t2) == 0 {
		t.Fatal("expected tokens")
	}
	if t1[0] != t2[0] {
		t.Errorf("expected same stem, got %q and %q", t1[0], t2[0])
	}
}
