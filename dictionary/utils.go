package dictionary

// Utilities for working with dictionaries.

import (
	"fmt"
	"strings"
)

type Any = interface{}

type LogEditsT = func(depth uint, v ...Any)

func logEdits(c chan<- []Any, maxDepth uint) LogEditsT {
	return func(depth uint, v ...Any) {
		tab := strings.Repeat("\t", int(maxDepth-depth))
		a := append([]interface{}{tab}, v...)
		c <- a
	}
}

func receiveLogs(c <-chan []Any, debug bool) {
	for {
		msg, ok := <-c
		if !ok {
			return
		}
		if debug {
			fmt.Println(msg...)
		}
	}
}
