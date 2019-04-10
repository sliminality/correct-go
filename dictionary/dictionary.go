package dictionary

import (
	"bufio"
	pq "github.com/golang-collections/go-datastructures/queue"
	"io"
	"log"
	trie "slim/correct/trie"
	"strings"
	"unicode"
)

// Represents a dictionary of words with corresponding frequencies.
type Dictionary struct {
	Root *trie.Node
}

// Given a file handle to a corpus, constructs a trie and inserts each word in
// the corpus.
func CreateDictionary(corpus io.Reader) Dictionary {
	root := trie.CreateNode("")
	scanner := bufio.NewScanner(corpus)

	isDelimiter := func(c rune) bool {
		return !unicode.IsLetter(c)
	}

	for scanner.Scan() {
		s := scanner.Text()
		for _, w := range strings.FieldsFunc(s, isDelimiter) {
			root.Insert(strings.ToLower(w))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return Dictionary{Root: &root}
}

// Check if the given word is in the dictionary.
//
// Returns one of the following results:
// (true,  nil)      if the word is in the dictionary
// (false, []string) if the word is not in the dictionary, with list of suggestions up to maxSuggestions in length
func (d *Dictionary) Check(q string, maxEdits uint, maxSuggestions uint, debug bool) (bool, []string, error) {
	_, found := d.Root.Lookup(q)
	if found {
		return found, nil, nil
	}
	// If word is not in dictionary, search for suggestions.
	suggestions, err := d.TopSuggestions(q, maxEdits, maxSuggestions, debug)
	return found, suggestions, err
}

// Returns the top correct suggestion by frequency, or nil if no suggestion is found.
func (d *Dictionary) TopSuggestions(q string, maxEdits uint, maxSuggestions uint, debug bool) ([]string, error) {
	all := d.suggestAll(q, maxEdits, debug)
	if len(all) == 0 {
		return nil, nil
	}

	heap := pq.NewPriorityQueue(len(all))
	for suggestion, wt := range all {
		s := Suggestion{Val: suggestion, Weight: int(wt)}
		heap.Put(s)
	}

	suggestions, err := heap.Get(int(maxSuggestions))
	if err != nil || len(suggestions) == 0 {
		return nil, err
	}
	results := make([]string, len(suggestions))
	for i, s := range suggestions {
		results[i] = s.(Suggestion).Val
	}
	return results, nil
}

type suggestions = map[string]uint

type Suggestion struct {
	Val    string
	Weight int
}

func (self Suggestion) Compare(other pq.Item) int {
	otherSuggestion := other.(Suggestion)
	if self.Weight == otherSuggestion.Weight {
		return 0
	}
	if self.Weight > otherSuggestion.Weight {
		return -1
	}
	return 1
}

// Generate a list of suggested corrections for the given word.
// TODO(slim): Use goroutines (hhhhhh)
// TODO(slim): Cache things maybe
func (d *Dictionary) suggestAll(q string, maxEdits uint, debug bool) suggestions {
	results := make(suggestions)
	logs := make(chan []Any)
	go receiveLogs(logs, debug)
	suggestAllRecur(d.Root, &results, q, maxEdits, logEdits(logs, maxEdits))
	close(logs)
	return results
}

// Recursively adds suggestions to the list.
func suggestAllRecur(n *trie.Node, results *suggestions, q string, maxEdits uint, logEdits LogEditsT) {
	// Base case: if no q remaining or no edits remaining, check if the current word is valid.
	if maxEdits == 0 {
		if node, ok := n.Lookup(q); ok {
			(*results)[node.Val] = node.Count
		}
		return
	}
	if len(q) == 0 {
		// Add the current word to the results list.
		if n.IsWord() {
			(*results)[n.Val] = n.Count
		}
		// Check for insertions at the end.
		suggestInsertions(n, results, q, 0, logEdits)
		return
	}

	logEdits(maxEdits, n.Val+"_"+q)

	// Recursive case #1: generate edits with a distance of 1.
	suggestInsertions(n, results, q, maxEdits-1, logEdits)
	suggestDeletions(n, results, q, maxEdits-1, logEdits)
	suggestTranspositions(n, results, q, maxEdits-1, logEdits)
	suggestReplacements(n, results, q, maxEdits-1, logEdits)

	// Recursive case #2: don't edit here, just continue.
	if child, ok := n.Advance(rune(q[0])); ok {
		suggestAllRecur(child, results, q[1:], maxEdits, logEdits)
	}
}

func suggestInsertions(n *trie.Node, results *suggestions, q string, maxEdits uint, logEdits LogEditsT) {
	for ch, c := range n.Children {
		logEdits(maxEdits, "Inserting", string(ch), "->", n.Val+string(ch)+q)
		suggestAllRecur(c, results, q, maxEdits, logEdits)
	}
}

func suggestDeletions(n *trie.Node, results *suggestions, q string, maxEdits uint, logEdits LogEditsT) {
	del := rune(q[0])
	rest := q[1:]
	logEdits(maxEdits, "Deleting", string(del), "->", n.Val+rest)
	suggestAllRecur(n, results, rest, maxEdits, logEdits)
}

// Swap the first two characters of the query.
func suggestTranspositions(n *trie.Node, results *suggestions, q string, maxEdits uint, logEdits LogEditsT) {
	if len(q) > 1 {
		snd := q[1]
		if node, ok1 := n.Advance(rune(snd)); ok1 {
			fst := string(q[0])
			rest := q[2:]
			logEdits(maxEdits, "Transposing", q[0:2],
				"->", n.Val+string(snd)+fst+rest)
			suggestAllRecur(node, results, fst+rest, maxEdits, logEdits)
		}
	}
}

// Replace the first character of the query.
func suggestReplacements(n *trie.Node, results *suggestions, q string, maxEdits uint, logEdits LogEditsT) {
	if len(q) > 0 {
		rest := q[1:]
		fst := string(q[0])
		for ch, child := range n.Children {
			logEdits(maxEdits, "Replacing", fst, "with", string(ch),
				"->", n.Val+string(ch)+rest)
			suggestAllRecur(child, results, rest, maxEdits, logEdits)
		}
	}
}
