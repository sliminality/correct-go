package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	dict "slim/correct/dictionary"
)

func main() {
	pathname := flag.String("corpus", "", "Path to corpus file")
	debug := flag.Bool("debug", false, "Print debug path")
	flag.Parse()

	file, err := os.Open(*pathname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	d := dict.CreateDictionary(file)

	// Read words from stdin.
	corrections := make(chan []string)

	go func() {
		for {
			c, ok := <-corrections
			if !ok {
				return
			}
			output := make([]byte, 0)
			for _, s := range c {
				output = append(output, []byte(s)...)
				output = append(output, byte(' '))
			}
			output = append(output, byte('\n'))
			os.Stdout.Write(output)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		s := scanner.Text()
		go func() {
			correct, suggestions, err := d.Check(s, 2, 3, *debug)
			if err != nil {
				panic(err)
			}
			if correct {
				result := []string{"✅", s}
				corrections <- result
			} else {
				result := []string{"❌", s}
				if len(suggestions) > 0 {
					result = append(result, "|")
				}
				corrections <- append(result, suggestions...)
			}
		}()
	}

	close(corrections)
}
