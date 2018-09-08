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

// PrettyPrint takes the site map and prints it as json
func PrettyPrint(site map[string][]string) {
	type SiteMapItem struct {
		URL   string   `json:"url"`
		Links []string `json:"links"`
	}

	fmt.Println("[")
	first := true
	for k, v := range site {
		item := SiteMapItem{URL: k, Links: v}
		jsonStr, err := json.Marshal(item)
		if err != nil {
			log.Printf("Error pretty printing: %v", err)
			return
		}
		if !first {
			fmt.Print(",")
		}
		fmt.Println(string(jsonStr))
		first = false
	}
	fmt.Println("]")
}

// ParseHTML takes a html body and returns a list of referred URLs
func ParseHTML(body string) []string {
	r := strings.NewReader(body)
	tokenizer := html.NewTokenizer(r)

	result := []string{}

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
						result = append(result, attr.Val)
					}
				}
			}
		}
	}

	return result
}

// FilterByDomain takes a list of URLs and remove the ones that doesn't belong
// to the given domain
func FilterByDomain(domain string, urls []string) []string {
	result := []string{}

	for _, rawurl := range urls {
		u, err := url.Parse(rawurl)
		if err != nil {
			log.Printf("Error parsing url %v: %v", rawurl, err)
			continue
		}
		if u.Hostname() == domain {
			result = append(result, rawurl)
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
func Crawl(rawurl string, level, delay int, verbose bool) (map[string][]string, error) {
	pu, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	domain := pu.Hostname()

	type Item struct {
		URL   string
		Level int
	}
	queue := []Item{Item{rawurl, 0}}

	adjList := map[string][]string{}

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		if _, ok := adjList[u.URL]; ok {
			continue
		}

		body, err := GetBody(u.URL)
		if err != nil {
			log.Printf("Error getting URL %v: %v\n", u.URL, err)
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
			log.Println(pad + u.URL)
		}

		adjList[u.URL] = urls
		filtered := FilterByDomain(domain, urls)

		if u.Level < level {
			for _, f := range filtered {
				queue = append(queue, Item{URL: f, Level: u.Level + 1})
			}
		}
	}

	return adjList, nil
}
