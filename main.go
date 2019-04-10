package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	// TODO(slim): Take the filename as an argument.
	bytes, err := ioutil.ReadFile("corpus.txt")
	if err != nil {
		// TODO(slim): Handle error lol
		return
	}
	fmt.Println(string(bytes))
}
