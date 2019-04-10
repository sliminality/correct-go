package trie

import (
	"testing"
)

func TestInit(t *testing.T) {
	root := new(Node)
	if root.Count != 0 {
		t.Error("When initialized, node has a count of 0")
	}
	if len(root.children) != 0 {
		t.Error("When initialized, node has 0 children")
	}
}

func TestInsert(t *testing.T) {
	root := CreateNode("")
	word := "hello"
	root.Insert(word)

	checkWord(t, &root, word, 1)
}

func TestInsertTwice(t *testing.T) {
	root := CreateNode("")
	word := "hello"

	root.Insert(word)
	root.Insert(word)

	checkWord(t, &root, word, 2)
}

func TestInsertSplit(t *testing.T) {
	root := CreateNode("")
	a := "hello"
	b := "goodbye"

	root.Insert(a)
	root.Insert(b)

	checkWord(t, &root, a, 1)
	checkWord(t, &root, b, 1)
}

func TestInsertPrefix(t *testing.T) {
	root := CreateNode("")
	word := "hello"
	cutIndex := 4
	prefix := word[:cutIndex] // "hell"

	root.Insert(word)
	root.Insert(prefix)

	node := &root
	for i, c := range word {
		if i == cutIndex {
			if node.Count != 1 {
				t.Error("prefix terminal node has count of 1")
			}
		} else if node.Count != 0 {
			t.Error("non-terminal node has count of 0")
		}

		if child, ok := node.children[c]; ok {
			node = child
		} else {
			t.Error("Child lookup failed for character", c)
		}
	}

	// Handle the terminal node separately.
	if node.Count != 1 {
		t.Error("terminal node has a count of 1")
	}
	if len(node.children) != 0 {
		t.Error("terminal node has no children")
	}
}

// Check that the given word has been inserted into the given trie node, the
// given number of times.
// TODO(slim): Rewrite this to take an array of (rune, count) pairs.
func checkWord(t *testing.T, root *Node, word string, timesInserted uint) {
	node := root

	for _, c := range word {
		if node.Count != 0 {
			t.Error("non-terminal node has count of 1")
		}
		if child, ok := node.children[c]; ok {
			node = child
		} else {
			t.Error("child lookup failed for character", c, "in word", word)
		}
	}

	// Handle the terminal node separately.
	if node.Count != timesInserted {
		t.Error("terminal node has a count of", timesInserted)
	}
	if len(node.children) != 0 {
		t.Error("terminal node has no children")
	}
}

func TestLookupWord(t *testing.T) {
	root := CreateNode("")
	root.Insert("hello")
	root.Insert("word")
	root.Insert("world")
	root.Insert("word")

	assertWordInTrie(t, &root, "hello", 1)
	assertWordInTrie(t, &root, "world", 1)
	assertWordInTrie(t, &root, "word", 2)
}

// Words in the trie should return (*node, true).
func assertWordInTrie(t *testing.T, root *Node, s string, count uint) {
	node, exists := root.Lookup(s)
	if !exists {
		t.Error("expected node for", s, "to exist")
	}
	if node.Count != count {
		t.Error("expected node to have count", count, "got", node.Count)
	}
}

func TestLookupPrefix(t *testing.T) {
	root := CreateNode("")
	root.Insert("hello")
	root.Insert("word")
	root.Insert("world")
	root.Insert("word")

	assertInTrieButNotWord(t, &root, "hell")
	assertInTrieButNotWord(t, &root, "wor")
}

// Prefixes in the trie should return (*node, false).
func assertInTrieButNotWord(t *testing.T, root *Node, s string) {
	node, exists := root.Lookup(s)
	if exists {
		t.Error("expected node for", s, "to not exist")
	}
	if node.Count != 0 {
		t.Error("expected node to have count 0, got", node.Count)
	}
}

func TestLookupFail(t *testing.T) {
	root := CreateNode("")
	root.Insert("hello")
	root.Insert("word")
	root.Insert("world")
	root.Insert("word")

	assertNotInTrie(t, &root, "hi")
	assertNotInTrie(t, &root, "worlddddddddd")
}

// Words that aren't in the trie should return (nil, false).
func assertNotInTrie(t *testing.T, root *Node, s string) {
	node, exists := root.Lookup(s)
	if exists {
		t.Error("expected node for", s, "to not exist")
	}
	if node != nil {
		t.Error("expected nil pointer for node")
	}
}
