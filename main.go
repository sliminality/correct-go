package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	pathname := flag.String("corpus", "", "path to corpus file")
	flag.Parse()

	file, err := os.Open(*pathname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
