package text

import (
	"math"
	"slices"
	"testing"
)

func seq(s ...string) func(yield func(string) bool) {
	return func(yield func(string) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

func TestTermFreq(t *testing.T) {
	tf := NewTermFreq(slices.Values([]string{"go", "go", "rust", "go"}))
	if got := tf.Get("go"); got != 0.75 {
		t.Errorf("Get(go) = %f, want 0.75", got)
	}
	if got := tf.Get("rust"); got != 0.25 {
		t.Errorf("Get(rust) = %f, want 0.25", got)
	}
}

func TestTermFreq_Empty(t *testing.T) {
	tf := NewTermFreq(seq())
	if !tf.Empty() {
		t.Errorf("expected empty, got %+v", tf)
	}
}

func TestCorpus_IDF(t *testing.T) {
	c := NewCorpus()
	c.Add(slices.Values([]string{"go", "rust"}))
	c.Add(slices.Values([]string{"go", "python"}))
	c.Add(slices.Values([]string{"rust", "python"}))

	idfGo := c.IDF("go")
	expected := math.Log(2.0)
	if math.Abs(idfGo-expected) > 0.001 {
		t.Errorf("IDF(go) = %f, want %f", idfGo, expected)
	}

	idfUnknown := c.IDF("java")
	expected = math.Log(4.0)
	if math.Abs(idfUnknown-expected) > 0.001 {
		t.Errorf("IDF(java) = %f, want %f", idfUnknown, expected)
	}
}

func TestCorpus_IDF_AlwaysPositive(t *testing.T) {
	c := NewCorpus()
	c.Add(slices.Values([]string{"go"}))

	idf := c.IDF("go")
	if idf <= 0 {
		t.Errorf("IDF should be positive, got %f", idf)
	}
}

func TestCorpus_DuplicateTerms(t *testing.T) {
	c := NewCorpus()
	c.Add(slices.Values([]string{"go", "go", "go"}))
	c.Add(slices.Values([]string{"rust"}))

	idf := c.IDF("go")
	expected := math.Log(2.0)
	if math.Abs(idf-expected) > 0.001 {
		t.Errorf("IDF(go) = %f, want %f", idf, expected)
	}
}

func TestWeight(t *testing.T) {
	c := NewCorpus()
	c.Add(slices.Values([]string{"go", "deploy"}))
	c.Add(slices.Values([]string{"go", "test"}))
	c.Add(slices.Values([]string{"deploy", "test"}))

	tf := NewTermFreq(slices.Values([]string{"go", "deploy"}))
	v := tf.Weight(c)
	if v.Empty() {
		t.Fatal("expected non-empty vector")
	}
	// Self-cosine should be 1.
	if got := v.Cosine(v); math.Abs(got-1.0) > 0.001 {
		t.Errorf("self-cosine = %f, want 1.0", got)
	}
}
