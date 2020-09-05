package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MichaelZalla/gophercises/link"
)

func main() {

	filepath := flag.String("filepath", "../data/ex1.html", "path to the HTML file to parse")

	flag.Parse()

	file, err := os.Open(*filepath)

	if err != nil {
		log.Fatal(fmt.Sprintf("failed to open file '%s' with error: '%s'", *filepath, err))
	}

	links, err := link.ParseLinks(file)

	if err != nil {
		log.Fatal(fmt.Sprintf("failed to parse links in file '%s' with error: '%s'", *filepath, err))
	}

	for i, l := range links {
		fmt.Printf("Link #%d: { Href: '%s', Text: '%s' }\n", i+1, l.Href, l.Text)
	}

	return

}
