package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/MichaelZalla/gophercises/cyoa"
)

var myStoryTemplate = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Choose Your Own Adventure</title>
	</head>
	<body>
		<section class="page">
			<h1>{{ .Title }}</h1>
			{{ range .Paragraphs }}
				<p>{{ . }}</p>
			{{ end }}
			<ul>
				{{ range .Options }}
				<li><a href="/{{ .Chapter }}">{{ .Text }}</a></li>
				{{ end }}
			</ul>
		</section>
	</body>
	</html>`

func myPathFn(r *http.Request) string {

	path := strings.TrimSpace(r.URL.Path)

	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}

	path = path[len("/story/"):]

	return path

}

func main() {

	port := flag.Int("port", 3000, "the port to start the CYOA web application on")

	filename := flag.String("file", "gopher.json", "the JSON file holding the CYOA story")

	flag.Parse()

	fmt.Printf("Using the story in %s.\n", *filename)

	f, err := os.Open(*filename)

	if err != nil {
		panic(err)
	}

	story, err := cyoa.GetStoryFromJson(f)

	if err != nil {
		panic(err)
	}

	tpl := template.Must(template.New("").Parse(myStoryTemplate))

	h := cyoa.NewHandler(story,
		cyoa.WithTemplate(tpl),
		cyoa.WithPathFn(myPathFn),
	)

	mux := http.NewServeMux()

	mux.Handle("/story/", h)

	fmt.Printf("Starting the server on port: %d\n", *port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))

}
