package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func main() {

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
			fmt.Fprintf(os.Stderr, "Error parsing url %v: %v", rawurl, err)
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

func Crawl(rawurl string) (map[string][]string, error) {
	pu, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	domain := pu.Hostname()

	queue := []string{rawurl}
	adjList := map[string][]string{}

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]

		body, err := GetBody(u)
		if err != nil {
			return nil, err
		}
		time.Sleep(100)

		urls := ParseHTML(body)
		filtered := FilterByDomain(domain, urls)

		adjList[u] = filtered
		for _, f := range filtered {
			queue = append(queue, f)
		}
	}

	return adjList, nil
}
