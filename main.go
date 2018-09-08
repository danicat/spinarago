package main

import (
	"flag"
	"log"
)

func main() {
	hostname := flag.String("hostname", "", "the name of the host to crawl")
	level := flag.Int("level", 10, "the max level of depth to crawl")
	delay := flag.Int("delay", 1000, "delay between each request in milliseconds")
	verbose := flag.Bool("verbose", false, "prints results as they come through")
	flag.Parse()

	if *hostname == "" {
		flag.PrintDefaults()
		return
	}

	log.Printf("Parsing %v with %v milliseconds interval and max depth %v.\n", *hostname, *delay, *level)
	siteMap, err := Crawl(*hostname, *level, *delay, *verbose)
	if err != nil {
		log.Println(err)
	} else {
		PrettyPrint(siteMap)
	}
}
