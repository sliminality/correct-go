# spellcheck

An implementation of [Peter Norvig](http://norvig.com/spell-correct.html)'s statistical spelling corrector, written in Go.

The concept and description below are from [Jesse Tov](http://users.eecs.northwestern.edu/~jesse/course/eecs396rust-sp18/)'s Rust course at Northwestern University.

## Concept

Given a word, an edit action is one of the following:

  - the deletion of one letter;
  - the transposition of two neighboring letters;
  - the replacement of one letter with another letter; and
  - the insertion of a letter at any position.

In this context, Norvig suggests that "small edits" means the application of one edit action possibly followed by the application of a second one to the result of the first.

Once the second part has generated all possible candidate for a potentially misspelled word, it picks the most frequently used one from the training corpus. If none of the candidates is a correct word, the program reports a failure.

## Implementation details

The program consumes a training file on the command line and then reads words—one per line—from standard input. For each word from standard in, the program prints one line. The line consists of just the word if it is spelled correctly. If the word is not correctly spelled, the program prints the word and the best improvement or "-" if there aren't any improvements found.

For example:

````
$ cat corpus.txt
hello world hello word hello world
$ cat input.txt
hello
hell
word
wordl
wor
wo
w
$ go run main.go --corpus corpus.txt < input.txt
hello
hell, hello
word
wordl, world
wor, world
wo, word
w, -
````
