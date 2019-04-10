package trie

import (
	"fmt"
	"strings"
)

// Represents either the root of a trie, or a character.
type Node struct {
	// Prefix represented by the current node.
	Val string
	// Frequency of words in the training corpus ending at the current trie node.
	Count uint
	// Trie nodes corresponding to suffix extensions of the current node.
	Children map[rune]*Node
}

func CreateNode(val string) Node {
	return Node{Val: val, Count: 0, Children: make(map[rune]*Node)}
}

// Insert a word into the trie.
func (n *Node) Insert(s string) {
	node := n
	for i, c := range s {
		if child, ok := node.Advance(c); ok {
			node = child
		} else {
			fresh := CreateNode(s[:i+1])
			node.Children[c] = &fresh
			node = &fresh
		}
	}
	node.Count += 1
}

// Does the given trie contain the given word?
func (n *Node) Lookup(s string) (*Node, bool) {
	node, ok := n.findNode(s)
	if !ok {
		return nil, false
	}
	// Node exists, but may not be a word.
	return node, node.IsWord()
}

// Does the current node correspond to a word from the corpus?
func (n *Node) IsWord() bool {
	return n.Count > 0
}

// Find the node corresponding to the given prefix.
func (n *Node) findNode(s string) (*Node, bool) {
	node := n
	for _, c := range s {
		if child, ok := node.Advance(c); ok {
			node = child
		} else {
			return nil, false
		}
	}
	return node, true
}

// Check whether the node has a child corresponding to the given rune.
func (n *Node) Advance(c rune) (*Node, bool) {
	child, ok := n.Children[c]
	return child, ok
}

// Print a trie for debugging.
func (n *Node) Debug() {
	n.debugRecur(0)
}

func (n *Node) debugRecur(depth int) {
	tab := strings.Repeat("  ", depth)
	fmt.Println(tab, n.Val)

	for _, c := range n.Children {
		c.debugRecur(depth + 1)
	}
}
