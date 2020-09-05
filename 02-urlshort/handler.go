package urlshort

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

type pathURL struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func getPathURLs(yml []byte) ([]pathURL, error) {

	var pathUrls []pathURL

	err := yaml.Unmarshal(yml, &pathUrls)

	if err != nil {
		return nil, err
	}

	return pathUrls, nil

}

func makeMap(pathUrls []pathURL) map[string]string {

	r := make(map[string]string)

	for _, pathURL := range pathUrls {
		r[pathURL.Path] = pathURL.URL
	}

	return r

}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		oldPath := r.URL.Path

		if newPath, ok := pathsToUrls[oldPath]; ok {
			http.Redirect(w, r, newPath, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)

	}

}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	pathUrls, err := getPathURLs(yml)

	if err != nil {
		return nil, err
	}

	pathsToUrls := makeMap(pathUrls)

	return MapHandler(pathsToUrls, fallback), nil

}
