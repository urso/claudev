package text

import (
	"math"
	"testing"
)

func TestTF(t *testing.T) {
	terms := []string{"go", "go", "rust", "go"}
	tf := TF(terms)
	if got := tf["go"]; got != 0.75 {
		t.Errorf("TF(go) = %f, want 0.75", got)
	}
	if got := tf["rust"]; got != 0.25 {
		t.Errorf("TF(rust) = %f, want 0.25", got)
	}
}

func TestTF_Empty(t *testing.T) {
	tf := TF(nil)
	if tf != nil {
		t.Errorf("expected nil, got %v", tf)
	}
}

func TestCorpus_IDF(t *testing.T) {
	c := NewCorpus()
	c.Add([]string{"go", "rust"})
	c.Add([]string{"go", "python"})
	c.Add([]string{"rust", "python"})

	// IDF formula: log(1 + N/(1+df))
	// "go" appears in 2/3 docs: log(1 + 3/3) = log(2)
	idfGo := c.IDF("go")
	expected := math.Log(2.0)
	if math.Abs(idfGo-expected) > 0.001 {
		t.Errorf("IDF(go) = %f, want %f", idfGo, expected)
	}

	// unknown term: log(1 + 3/1) = log(4)
	idfUnknown := c.IDF("java")
	expected = math.Log(4.0)
	if math.Abs(idfUnknown-expected) > 0.001 {
		t.Errorf("IDF(java) = %f, want %f", idfUnknown, expected)
	}
}

func TestCorpus_IDF_AlwaysPositive(t *testing.T) {
	c := NewCorpus()
	c.Add([]string{"go"})

	// Single doc, term in it: log(1 + 1/2) > 0
	idf := c.IDF("go")
	if idf <= 0 {
		t.Errorf("IDF should be positive, got %f", idf)
	}
}

func TestCorpus_DuplicateTerms(t *testing.T) {
	c := NewCorpus()
	c.Add([]string{"go", "go", "go"}) // should count as 1 doc
	c.Add([]string{"rust"})

	// "go" in 1/2 docs: log(1 + 2/2) = log(2)
	idf := c.IDF("go")
	expected := math.Log(2.0)
	if math.Abs(idf-expected) > 0.001 {
		t.Errorf("IDF(go) = %f, want %f", idf, expected)
	}
}

func TestTFIDF(t *testing.T) {
	c := NewCorpus()
	c.Add([]string{"go", "deploy"})
	c.Add([]string{"go", "test"})
	c.Add([]string{"deploy", "test"})

	tf := TF([]string{"go", "deploy"})
	scores := TFIDF(tf, c)

	// Both terms in 2/3 docs: TF=0.5, IDF=log(1+3/3)=log(2)
	expected := 0.5 * math.Log(2.0)
	for term, score := range scores {
		if math.Abs(score-expected) > 0.001 {
			t.Errorf("TFIDF(%s) = %f, want %f", term, score, expected)
		}
	}
}
