package text

import (
	"math"
	"slices"
	"testing"
)

func TestWeightedFields_Score(t *testing.T) {
	c := NewCorpus()
	c.Add(slices.Values([]string{"auth", "login"}))
	c.Add(slices.Values([]string{"deploy", "ci"}))

	title := NewTermFreq(slices.Values([]string{"auth"}))      // tf[auth]=1
	body := NewTermFreq(slices.Values([]string{"auth", "ci"})) // tf[auth]=0.5, tf[ci]=0.5

	w := NewWeightedFields().
		Add(80, title).
		Add(10, body)

	got := w.Score(slices.Values([]string{"auth"}), c)
	// combined tf[auth] = 80*1 + 10*0.5 = 85
	// idf(auth) = log(1 + 2/(1+1)) = log(2)
	want := 85.0 * math.Log(2.0)
	if math.Abs(got-want) > 0.001 {
		t.Errorf("Score = %f, want %f", got, want)
	}
}

func TestWeightedFields_NoQueryHit(t *testing.T) {
	c := NewCorpus()
	c.Add(slices.Values([]string{"x"}))
	w := NewWeightedFields().Add(1, NewTermFreq(slices.Values([]string{"x"})))
	if got := w.Score(slices.Values([]string{"y"}), c); got != 0 {
		t.Errorf("Score = %f, want 0", got)
	}
}

func TestWeightedFields_Empty(t *testing.T) {
	c := NewCorpus()
	w := NewWeightedFields()
	if got := w.Score(slices.Values([]string{"a"}), c); got != 0 {
		t.Errorf("Score = %f, want 0", got)
	}
}
