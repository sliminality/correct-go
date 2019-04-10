package dictionary

import (
	"bufio"
	"io"
	"log"
	trie "slim/correct/trie"
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

	for scanner.Scan() {
		root.Insert(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return Dictionary{Root: &root}
}
