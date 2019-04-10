package dictionary

import (
	assert "slim/correct/assert"
	trie "slim/correct/trie"
	"testing"
)

func TestCheckCorrect(t *testing.T) {
	d := createTestDictionary([]string{
		"hello",
		"word",
		"world",
		"word",
	})
	assertCorrect(t, &d, "hello")
	assertCorrect(t, &d, "word")
	assertCorrect(t, &d, "world")
}

func TestCheckSuggestion(t *testing.T) {
	d := createTestDictionary([]string{
		"hello",
		"word",
		"world",
		"word",
	})
	assertSuggestion(t, &d, "hell", "hello")
	assertSuggestion(t, &d, "wor", "word")
	assertSuggestion(t, &d, "worldd", "word") // "word" is more frequent than "world"
	assertSuggestion(t, &d, "worlldd", "world")
}

func TestCheckNoSuggestion(t *testing.T) {
	d := createTestDictionary([]string{
		"hello",
		"word",
		"world",
		"word",
	})
	assertNoSuggestion(t, &d, "he")
	assertNoSuggestion(t, &d, "worldddd")
	assertNoSuggestion(t, &d, "wrrrly")
}

func TestSuggestAll_zero_edits(t *testing.T) {
	d := createTestDictionary([]string{
		"word",
		"work",
		"world",
		"word",
		"work",
		"work",
	})
	results := d.suggestAll("owrd", 0)

	assert.Equal(t)(len(results), 0,
		"Expected no suggestions, received", len(results))
}

func TestSuggestAll_one_edit(t *testing.T) {
	d := createTestDictionary([]string{
		"word",
		"work",
		"world",
		"word",
		"work",
		"work",
	})
	expected := map[string]uint{
		"word": 2,
	}
	results := d.suggestAll("owrd", 1)

	assert.Equal(t)(len(results), 1,
		"Expected 1 suggestion, received", len(results))

	for s, wt := range results {
		assert.Equal(t)(wt, expected[s],
			"Expected suggestion", s, "to have weight", expected[s])
	}
}

func TestSuggestAll_two_edits(t *testing.T) {
	d := createTestDictionary([]string{
		"word",
		"work",
		"world",
		"word",
		"work",
		"work",
	})
	expected := map[string]uint{
		"work":  3,
		"world": 1,
		"word":  2,
	}
	results := d.suggestAll("owrd", 2)

	assert.Equal(t)(len(results), len(expected),
		"Expected", len(expected), "suggestions, received", len(results))

	for s, wt := range results {
		assert.Equal(t)(wt, expected[s],
			"Expected suggestion", s, "to have weight", expected[s])
	}
}

func TestSuggestAll_three_edits(t *testing.T) {
	d := createTestDictionary([]string{
		"word",
		"work",
		"world",
		"word",
		"work",
		"work",
	})
	expected := map[string]uint{
		"work":  3,
		"world": 1,
		"word":  2,
	}
	results := d.suggestAll("owrd", 3)

	assert.Equal(t)(len(results), len(expected),
		"Expected", len(expected), "suggestions, received", len(results))

	for s, wt := range results {
		assert.Equal(t)(wt, expected[s],
			"Expected suggestion", s, "to have weight", expected[s])
	}
}

func createTestDictionary(words []string) Dictionary {
	root := trie.CreateNode("")
	for _, s := range words {
		root.Insert(s)
	}
	d := Dictionary{Root: &root}
	return d
}

func assertCorrect(t *testing.T, d *Dictionary, q string) {
	correct, suggestions, _ := d.Check(q, 2, 1)
	assert.True(t)(correct,
		"Expected", q, "to be spelled correctly")
	assert.Equal(t)(len(suggestions), 0,
		"Expected no suggestions for", q, "but got", suggestions)
}

func assertSuggestion(t *testing.T, d *Dictionary, q string, exp string) {
	correct, suggestions, _ := d.Check(q, 2, 1)
	assert.False(t)(correct,
		"Expected", q, "to be misspelled")
	assert.Equal(t)(len(suggestions), 1,
		"Expected 1 suggestion, but got", len(suggestions))
	s := suggestions[0]
	assert.Equal(t)(s, exp,
		"Expected suggestion", exp, "but got nil")
}

func assertNoSuggestion(t *testing.T, d *Dictionary, q string) {
	correct, suggestions, _ := d.Check(q, 2, 1)
	assert.False(t)(correct,
		"Expected", q, "to be misspelled")
	assert.Equal(t)(len(suggestions), 0,
		"Expected no suggestions for", q, "but got", suggestions)
}
