package text

import (
	"strings"
	"testing"
)

func TestTokenizer_Basic(t *testing.T) {
	tok := NewTokenizer()
	tokens := tok.Tokenize("Running deployments quickly")
	// "running" → "run", "deployments" → "deploy", "quickly" → "quick"
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
	tokens := tok.Tokenize("the quick and the dead")
	for _, tok := range tokens {
		if tok == "the" || tok == "and" {
			t.Errorf("stopword %q should have been removed", tok)
		}
	}
}

func TestTokenizer_Empty(t *testing.T) {
	tok := NewTokenizer()
	tokens := tok.Tokenize("")
	if len(tokens) != 0 {
		t.Errorf("expected 0 tokens, got %d", len(tokens))
	}
}

func TestTokenizer_Punctuation(t *testing.T) {
	tok := NewTokenizer()
	tokens := tok.Tokenize("hello-world foo_bar baz.qux")
	if len(tokens) != 6 {
		t.Errorf("expected 6 tokens, got %d: %v", len(tokens), tokens)
	}
}

func TestTokenizer_Stemming(t *testing.T) {
	tok := NewTokenizer()
	// "architecture" and "architectural" should stem to the same thing
	t1 := tok.Tokenize("architecture")
	t2 := tok.Tokenize("architectural")
	if len(t1) == 0 || len(t2) == 0 {
		t.Fatal("expected tokens")
	}
	if t1[0] != t2[0] {
		t.Errorf("expected same stem, got %q and %q", t1[0], t2[0])
	}
}
