package main

import (
	"flag"
	"log"
	"os"
	dict "slim/correct/dictionary"
)

func main() {
	pathname := flag.String("corpus", "", "path to corpus file")
	flag.Parse()

	file, err := os.Open(*pathname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	d := dict.CreateDictionary(file)
	d.Root.Debug()
}
