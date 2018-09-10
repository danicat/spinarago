package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// PrettyPrint takes a site map and prints it as json
func PrettyPrint(site map[string][]*url.URL) {
	type SiteMapItem struct {
		URL   string   `json:"url"`
		Links []string `json:"links"`
	}

	fmt.Print("[")
	first := true
	for k, v := range site {
		links := []string{}
		for _, u := range v {
			links = append(links, u.String())
		}

		item := SiteMapItem{URL: k, Links: links}
		jsonStr, err := json.Marshal(item)
		if err != nil {
			log.Printf("Error pretty printing: %v", err)
			return
		}
		if !first {
			fmt.Print(",")
		}
		fmt.Print(string(jsonStr))
		first = false
	}
	fmt.Print("]")
}

// ParseHTML takes a html body and returns a list of referred URLs
func ParseHTML(body string) []*url.URL {
	r := strings.NewReader(body)
	tokenizer := html.NewTokenizer(r)

	result := []*url.URL{}

	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.StartTagToken {
			t := tokenizer.Token()

			if t.Data == "a" {
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						u, err := url.Parse(attr.Val)
						if err != nil {
							log.Printf("Error parsing url %v: %v", attr.Val, err)
							continue
						}
						result = append(result, u)
					}
				}
			}
		}
	}

	return result
}

// FilterByHostname takes a list of URLs and remove the ones that doesn't belong
// to the given hostname
func FilterByHostname(hostname string, urls []*url.URL) []*url.URL {
	result := []*url.URL{}

	for _, u := range urls {
		// TODO: handle relative paths
		target := u.Hostname()
		if len(target) < len(hostname) {
			continue
		}

		match := true
		for i, j := len(hostname)-1, len(target)-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
			if hostname[i] != target[j] {
				match = false
				break
			}
		}
		if match {
			result = append(result, u)
		}
	}

	return result
}

// GetBody takes an URL as parameter, performs a GET and returns its body
func GetBody(rawurl string) (string, error) {
	resp, err := http.Get(rawurl)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Crawl takes a base URL and builds a site map from it.
// delay is the time in miliseconds to wait between one request and another
func Crawl(baseURL string, level, delay int, verbose bool) (map[string][]*url.URL, error) {
	bu, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	hostname := bu.Hostname()

	type Item struct {
		URL   *url.URL
		Level int
	}
	queue := []Item{Item{bu, 0}}

	adjList := map[string][]*url.URL{}

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		currURL := u.URL.String()

		if _, ok := adjList[currURL]; ok {
			continue
		}

		body, err := GetBody(currURL)
		if err != nil {
			log.Printf("Error getting URL %v: %v\n", currURL, err)
			continue
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)

		urls := ParseHTML(body)
		if verbose {
			var pad string
			for i := 0; i <= u.Level; i++ {
				pad = pad + "-"
			}
			pad += ">"
			log.Println(pad + currURL)
		}

		adjList[currURL] = urls
		filtered := FilterByHostname(hostname, urls)

		if u.Level < level {
			for _, f := range filtered {
				queue = append(queue, Item{URL: f, Level: u.Level + 1})
			}
		}
	}

	return adjList, nil
}
