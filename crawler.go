package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func main() {
	//var hostname string
	//flag.String(&hostname, "hostname", "The hostname to crawl.")
	hostname := flag.String("hostname", "", "the name of the host to crawl")
	level := flag.Int("level", 10, "the max level of depth to crawl")
	delay := flag.Int("delay", 1000, "delay between each request in milliseconds")
	flag.Parse()

	fmt.Printf("Parsing %v with %v milliseconds interval.\n", *hostname, *delay)
	siteMap, err := Crawl(*hostname, *level, *delay)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(siteMap)
	}
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

// Crawl takes a base URL and builds a site map from it.
// delay is the time in miliseconds to wait between one request and another
func Crawl(rawurl string, level, delay int) (map[string][]string, error) {
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
			// return nil, err
			log.Printf("Error getting URL %v: %v\n", u.URL, err)
		}
		time.Sleep(time.Duration(delay) * time.Millisecond)

		urls := ParseHTML(body)
		filtered := FilterByDomain(domain, urls)

		adjList[u.URL] = filtered

		if u.Level < level {
			for _, f := range filtered {
				queue = append(queue, Item{URL: f, Level: u.Level + 1})
			}
		}
	}

	return adjList, nil
}
