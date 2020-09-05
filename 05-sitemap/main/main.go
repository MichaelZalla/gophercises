package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/MichaelZalla/gophercises/sitemap/crawl"
)

func main() {

	originURL := flag.String("origin", "", "the origin URL for the sitemap")

	depth := flag.Int("depth", 5, "the maximum crawl (search) depth")

	crossSite := flag.Bool("cross-site", false, "follow links to external sites")

	// out := flag.String("out", "", "output file location (default: stdout)")

	// errorOut := flag.String("error-out", "", "error file location (default: stdout)")

	// verbose := flag.Bool("verbose", false, "logs activity and errors (optional)")

	// Init our flag values

	flag.Parse()

	// Perform the site crawl

	var w io.Writer

	// if *out != "" {

	// 	file, err := os.Create(*out)

	// 	if err != nil {
	// 		log.Fatal(fmt.Sprintf("failed to open file '%s'", *out))
	// 	}

	// 	defer file.Close()

	// 	w = file

	// } else {

	w = os.Stdout

	// }

	err := crawl.GetSitemap(*originURL, *depth, *crossSite, w)

	if err != nil {
		log.Fatal(fmt.Sprintf("failed to write XML to writer '%s'", w))
	}

}
