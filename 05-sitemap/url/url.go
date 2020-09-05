package url

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// See: https://regex101.com/r/sb0nQu/1/
var urlComponentsRegex = regexp.MustCompile(`^((?P<protocol>(http(s)?|feed|sketch):\/\/|mailto:|javascript:)?(?P<wwwSubdomain>www\.)?)(?P<canonicalDomain>(?P<domain>[a-z0-9\-@]+(\.[a-z0-9\-]+)*)(?P<tld>\.[a-z]+))?(?P<path>\/[^?#\s]*)?(?P<params>\?([^\s#])*)?(?P<search>#.*)?$`)

var fileNameRegex = regexp.MustCompile(`(?P<dirname>.+\/)(?P<filename>[^\/]+\.[^\/]+)?$`)

// Fields contains the individual field components of a given URL
type Fields struct {
	Protocol        string
	WwwSubdomain    string
	CanonicalDomain string
	Domain          string
	Tld             string
	Path            string
	Params          string
	Search          string
	CanonicalURL    string
}

// func (fields Fields) String() string {

// 	f := make([]string, 9)

// 	f[0] = fmt.Sprintf("Protocol = %s, ", fields.Protocol)
// 	f[1] = fmt.Sprintf("WwwSubdomain = %s, ", fields.WwwSubdomain)
// 	f[2] = fmt.Sprintf("CanonicalDomain = %s, ", fields.CanonicalDomain)
// 	f[3] = fmt.Sprintf("Domain = %s, ", fields.Domain)
// 	f[4] = fmt.Sprintf("Tld = %s, ", fields.Tld)
// 	f[5] = fmt.Sprintf("Path = %s, ", fields.Path)
// 	f[6] = fmt.Sprintf("Params = %s, ", fields.Params)
// 	f[7] = fmt.Sprintf("Search = %s, ", fields.Search)
// 	f[8] = fmt.Sprintf("CanonicalURL = %s, ", fields.CanonicalURL)

// 	ret := "{ "

// 	for i := range f {
// 		ret += f[i] + " "
// 	}

// 	return ret + "}"
// }

func getFields(url string) (Fields, error) {

	fields := Fields{}

	namedMatches := make(map[string]string)

	submatches := urlComponentsRegex.FindStringSubmatch(url)
	subexprNames := urlComponentsRegex.SubexpNames()

	if len(submatches) == 0 {
		return fields, fmt.Errorf("failed to parse URL '%s'", url)
	}

	for i, name := range subexprNames {
		if i != 0 && name != "" {
			namedMatches[name] = submatches[i]
		}
	}

	fields.Protocol = namedMatches["protocol"]
	fields.WwwSubdomain = namedMatches["wwwSubdomain"]
	fields.CanonicalDomain = namedMatches["canonicalDomain"]
	fields.Domain = namedMatches["domain"]
	fields.Tld = namedMatches["tld"]
	fields.Path = namedMatches["path"]
	fields.Params = namedMatches["params"]
	fields.Search = namedMatches["search"]

	if fields.Path != "" && strings.Index(fields.Path, ".") == -1 && fields.Path[len(fields.Path)-1] != '/' {
		fields.Path += "/"
	}

	fields.CanonicalURL = getCanonicalURL(fields)

	return fields, nil

}

func getDir(path string) string {

	submatches := fileNameRegex.FindStringSubmatch(path)

	if len(submatches) > 0 {
		return submatches[1]
	}

	return "/"

}

// NormalizeURL normalizes a path or href relative to its referer
func NormalizeURL(url string, referer string) (string, Fields, error) {

	var refererFields Fields
	var err error

	if referer != "" {

		refererFields, err = getFields(referer)

		if err != nil {
			return "", Fields{}, fmt.Errorf("failed to call NormalizeURL with arguments '%s', '%s' (error: %s)", url, referer, err)
		}

		if refererFields.Path != "" && refererFields.Path != "/" {

			for strings.HasPrefix(url, "./") {
				url = url[2:]
			}

			for strings.HasPrefix(url, "../") {
				url = url[3:]
				refererFields.Path = getDir(refererFields.Path)
				refererFields.Path = refererFields.Path[:strings.LastIndex(refererFields.Path, "/")]
			}

		}

	}

	f, err := getFields(url)

	if err != nil {
		return "", Fields{}, fmt.Errorf("failed to call NormalizeURL with arguments '%s', '%s' (error: %s)", url, referer, err)
	}

	if f.Protocol == "" {
		f.Protocol = "https://"
	}

	if f.Path == "" {
		if refererFields.Path != "" {
			f.Path = refererFields.Path
		} else {
			f.Path = "/"
		}
	}

	if referer != "" {

		if f.CanonicalDomain == "" {
			f.CanonicalDomain = refererFields.CanonicalDomain
			f.Domain = refererFields.Domain
			f.Tld = refererFields.Tld
		}

	}

	return getCanonicalURL(f), f, nil

}

func getCanonicalURL(f Fields) string {
	return f.Protocol + f.WwwSubdomain + f.CanonicalDomain + f.Path + f.Params + f.Search
}

// GetBodyReader returns an io.Reader reading from the response body of the given URL
func GetBodyReader(url string) (io.Reader, error) {

	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	body := string(bodyBytes)

	return strings.NewReader(body), nil

}
