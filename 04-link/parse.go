package link

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link a hyperlink in an HTML document
type Link struct {
	Href string
	Text string
}

func parse(r io.Reader) (*html.Node, error) {

	root, err := html.Parse(r)

	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML with error: '%s'", err)
	}

	htmlRoot := html.Node(*root)

	return &htmlRoot, nil

}

// ParseLinks will take an HTML document and return a list of links, or an error
func ParseLinks(r io.Reader) ([]Link, error) {

	root, err := parse(r)

	if err != nil {
		return nil, err
	}

	anchors := getAnchorNodes(root)

	links := getLinks(anchors)

	return links, nil

}

func getAnchorNodes(r *html.Node) []*html.Node {

	var anchors []*html.Node

	dfs(r, func(n *html.Node) {

		if n.Type == html.ElementNode && n.Data == "a" {
			anchors = append(anchors, n)
		}

	}, 0)

	return anchors

}

func getAttribute(n *html.Node, key string) (string, error) {

	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val, nil
		}
	}

	return "", fmt.Errorf("node %v has no value for attribute key '%s'", n, key)

}

func getText(n *html.Node) string {

	text := ""

	dfs(n, func(n *html.Node) {

		if n.Type == html.TextNode {
			text += n.Data
		}

	}, 0)

	return strings.Join(strings.Fields(text), " ")

}

func getLink(anchor *html.Node) Link {

	ret := Link{}

	href, err := getAttribute(anchor, "href")

	if err == nil {
		ret.Href = strings.TrimSpace(href)
	}

	ret.Text = strings.TrimSpace(getText(anchor))

	return ret

}

func getLinks(anchors []*html.Node) []Link {

	links := make([]Link, len(anchors))

	for i, a := range anchors {
		links[i] = getLink(a)
	}

	return links

}

type visitFn func(n *html.Node)

func dfs(r *html.Node, visit visitFn, depth int) {

	out := r.Data

	if r.Type == html.ElementNode {
		out = "<" + out + ">"
	}

	// fmt.Printf("%s", strings.Repeat("\t", depth))
	// fmt.Printf("%v\n", out)

	visit(r)

	if r.FirstChild != nil {

		depth++

		for child := r.FirstChild; child != nil; child = child.NextSibling {
			dfs(child, visit, depth)
		}

		depth--

	}

	if r.Type == html.ElementNode {
		out = "</" + out[1:]
		// fmt.Printf("%s", strings.Repeat("\t", depth))
		// fmt.Printf("%v\n", out)
	}

	return

}
