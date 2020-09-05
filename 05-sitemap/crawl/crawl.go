package crawl

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/MichaelZalla/gophercises/link"
	"github.com/MichaelZalla/gophercises/sitemap/url"
)

type urlCacheEntry struct {
	URL          string
	FirstReferer string
	Depth        int
	Timestamp    int64
}

type urlCache map[string]urlCacheEntry

type edge struct {
	From string
	To   string
}

func crawl(originURL string, maxDepth int, crossOrigin bool) []urlCacheEntry {

	log.Printf("Starting crawl...")

	start := time.Now()

	// Initialize a cache to keep track of which URL fetches have been scheduled

	scheduledCache := make(urlCache)

	// Perform BFS to iterate through all nested (descendent) links efficiently

	normalizedOriginURL, originURLFields, err := url.NormalizeURL(originURL, "")

	if err != nil {
		log.Fatal(fmt.Printf("failed to initiate crawl from '%s' (error: %s) ", originURL, err))
	}

	currentFrontier := []edge{}

	futureFrontier := []edge{edge{"", normalizedOriginURL}}

	// While our queue is not empty...

	for currentDepth := 0; currentDepth <= maxDepth; currentDepth++ {

		currentFrontier, futureFrontier = futureFrontier, []edge{}

		finished := 0

		futureCh := make(chan edge)

		loggerCh := make(chan string)

		doneCh := make(chan bool, len(currentFrontier))

		for _, e := range currentFrontier {

			firstReferer := e.From

			href := e.To

			// Check whether or not this URL has already been visited

			if _, ok := scheduledCache[href]; ok {
				doneCh <- true
				continue
			}

			entry := scheduledCache[href]
			entry.URL = href
			entry.FirstReferer = firstReferer
			entry.Depth = currentDepth
			entry.Timestamp = time.Now().Unix()

			scheduledCache[href] = entry

			// If we're processing the last level of the crawl, mark this href in our cache, but don't crawl it

			if currentDepth == maxDepth {
				doneCh <- true
				continue
			}

			// Otherwise, we need to crawl this href, in parallel with the others at this level

			go func(future chan<- edge, done chan<- bool, logger chan<- string) {

				defer func() {
					done <- true
				}()

				// Fetch the response (body) for this href

				logger <- fmt.Sprintf(`GET %s`, href)

				body, err := url.GetBodyReader(href)

				if err != nil {
					logger <- fmt.Sprintf("failed to GET %s (error: %s)", href, err)
					return
				}

				// Parse all suitable links from the response (body)

				links, err := link.ParseLinks(body)

				if err != nil {
					logger <- fmt.Sprintf("failed to parse links from response body for %s (error: %s)", href, err)
					return
				}

				// Process links

				for _, link := range links {

					// In this context, 'href' is the first referer to link.Href

					parentURL := href

					childURL := link.Href

					// Ignore empty links, unsupported protocols, and unsupported file types

					if childURL == "" {
						continue
					}

					if strings.HasPrefix(childURL, "#") {
						continue
					}

					if strings.Contains(childURL, "feed:") ||
						strings.Contains(childURL, "mailto:") ||
						strings.Contains(childURL, "sketch:") ||
						strings.Contains(childURL, "javascript:") {
						continue
					}

					if strings.Contains(childURL, ".dmg") ||
						strings.Contains(childURL, ".zip") ||
						strings.Contains(childURL, ".pdf") {
						continue
					}

					// Normalize the child URL, using the parent URL if needed

					var normalizedChildURL string
					var childURLFields url.Fields
					var err error

					if childURL[0] == '/' || childURL[0] == '.' || childURL[0] == '#' {
						normalizedChildURL, childURLFields, err = url.NormalizeURL(childURL, parentURL)
					} else {
						normalizedChildURL, childURLFields, err = url.NormalizeURL(childURL, "")
					}

					if err != nil {
						logger <- fmt.Sprintf(err.Error())
					}

					// Ignore cross-origin child URLs, unless these were requested for the crawl

					if crossOrigin == false && childURLFields.CanonicalDomain != originURLFields.CanonicalDomain {
						// logger <- fmt.Sprintf("Skipping external link: %s", normalizedChildURL)
						continue
					}

					// Send a new edge to our future frontier (channel)

					future <- edge{parentURL, normalizedChildURL}

				}

			}(futureCh, doneCh, loggerCh)

		}

		// Wait for a signal that all of our current frontier's paths have been processed

		// All sends to futureCh and loggerCn should happen before we call close(...)

	readChannels:
		for {
			select {
			case e := <-futureCh:
				if _, ok := scheduledCache[e.To]; !ok {
					futureFrontier = append(futureFrontier, e)
				}
			case msg := <-loggerCh:
				fmt.Println(msg)
			case <-doneCh:
				finished++
				if finished == len(currentFrontier) {
					close(futureCh)
					close(loggerCh)
					close(doneCh)
					break readChannels
				}
			}

		}

	}

	// Collect all URL cache entries into a list, and return it

	ret := make([]urlCacheEntry, 0, len(scheduledCache))

	for _, entry := range scheduledCache {
		ret = append(ret, entry)
	}

	// Time profiling

	log.Printf("\n\nCrawled %d unique webpages (duration %s seconds)", len(ret), time.Since(start))

	return ret

}

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	URLs  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

// GetSitemap generates a standard XML representation of a sitemap for the given origin, up to a maximum crawl depth
func GetSitemap(originURL string, maxDepth int, crossOrigin bool, w io.Writer) error {

	entries := crawl(originURL, maxDepth, crossOrigin)

	toXML := urlset{
		URLs:  make([]loc, len(entries)),
		Xmlns: xmlns,
	}

	for i, entry := range entries {
		toXML.URLs[i] = loc{entry.URL}
	}

	enc := xml.NewEncoder(w)

	enc.Indent("", "  ")

	fmt.Printf(xml.Header)

	err := enc.Encode(toXML)

	fmt.Printf("\n")

	if err != nil {
		return err
	}

	// ret := xml.Header
	// ret += "<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n"

	// xml := make([]string, len(entries))

	// for i, entry := range entries {
	// 	xml[i] = fmt.Sprintf("\t<url><loc>%s</loc>"+ /*"<referer>%s</referer>" + */ "<lastmod>%s</lastmod></url>\n", entry.URL /*entry.FirstReferer, */, time.Unix(entry.Timestamp, 0))
	// }

	// sort.Sort(sort.StringSlice(xml))

	// ret += strings.Join(xml, "")

	// ret += "</urlset>"

	return nil

}
