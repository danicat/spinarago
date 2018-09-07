package main

import (
	"io"
	"net/http"
	"testing"
)

var index0 = `
<html>
<title>Home</title>
<body>
<p><a href="http://localhost:8080/abc.html">abc</a></p>
</body>
</html>
`

var index1 = `
<html>
<title>Home</title>
<body>
<p><a href="http://localhost:8080/abc.html">abc</a></p>
<p><a href="http://localhost:8080/def.html">abc</a></p>
<p><a href="http://localhost:8080/ghi.html">abc</a></p>
</body>
</html>
`

var index2 = `
<html>
<title>Home</title>
<body>
</body>
</html>
`

var index3 = `
<html>
<title>Home</title>
<body>
<p>//localhost:8080</p>
<p>http://abc.com</p>
<p><a href="http://localhost:8080/abc.html">abc</a></p>
<p><a href="http://localhost:8080/def.html">abc</a></p>
<p><a href="http://localhost:8080/ghi.html">abc</a></p>
<p><a href="http://abc.com/ghi.html">abc</a></p>
</body>
</html>
`

var home = `
<html>
<title>Home</title>
<body>
<p><a href="http://localhost:8080/abc.html">abc</a></p>
</body>
</html>
`

var abc = `
<html>
<title>ABC</title>
<body>
<p><a href="http://localhost:8080/def.html">abc</a></p>
</body>
</html>
`

func TestParseHTML(t *testing.T) {
	table := []struct {
		input    string
		expected []string
	}{
		{
			index0,
			[]string{
				"http://localhost:8080/abc.html",
			},
		},
		{
			index1,
			[]string{
				"http://localhost:8080/abc.html",
				"http://localhost:8080/def.html",
				"http://localhost:8080/ghi.html",
			},
		},
		{
			index2,
			[]string{},
		},
		{
			index3,
			[]string{
				"http://localhost:8080/abc.html",
				"http://localhost:8080/def.html",
				"http://localhost:8080/ghi.html",
				"http://abc.com/ghi.html",
			},
		},
	}

	for testnum, test := range table {
		result := ParseHTML(test.input)
		if len(result) != len(test.expected) {
			t.Fatalf("\nTest: %v\nExpected %v\nGot %v\n", testnum, test.expected, result)
		}

		for i, r := range result {
			if r != test.expected[i] {
				t.Fatalf("\nTest: %v\nExpected %v\nGot %v\n", testnum, test.expected, result)
			}
		}
	}
}

func TestFilterByDomain(t *testing.T) {
	input := []string{
		"//abc.com/index.html",
		"//abc.com/blablabla.html",
		"//def.com/index.html",
		"",
	}
	domain := "abc.com"
	expected := []string{
		"//abc.com/index.html",
		"//abc.com/blablabla.html",
	}

	result := FilterByDomain(domain, input)
	if len(result) != len(expected) {
		t.Fatalf("\nExpected %v\nGot %v\n", expected, result)
	}

	for i, r := range result {
		if r != expected[i] {
			t.Fatalf("\nExpected %v\nGot %v\n", expected, result)
		}
	}
}

func TestCrawl(t *testing.T) {
	input := "http://localhost:8080"
	expected := map[string][]string{
		"http://localhost:8080": {
			"http://localhost:8080/abc.html",
			// "http://localhost:8080/def.html",
			// "http://localhost:8080/ghi.html",
		},
		"http://localhost:8080/abc.html": {
			"http://localhost:8080/def.html",
		},
		"http://localhost:8080/def.html": {
			// "http://localhost:8080/ghi.html",
			// "http://localhost:8080/jkl.html",
		},
		// "http://localhost:8080/ghi.html": {
		// 	"http://localhost:8080/abc.html",
		// },
	}

	homeHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, home)
	}

	abcHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, abc)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/abc.html", abcHandler)

	go http.ListenAndServe(":8080", nil)

	result, err := Crawl(input)
	if err != nil {
		t.Fatalf("Test failed with error: %v", err)
	}

	if len(result) != len(expected) {
		t.Fatalf("\nExpected: %v\nGot: %v\n", expected, result)
	}

	for k, v := range result {
		if ev, ok := expected[k]; !ok {
			t.Fatalf("\nExpected: %v\nGot: %v\n", expected, result)
		} else {
			for i, u := range v {
				if u != ev[i] {
					t.Fatalf("\nExpected: %v\nGot: %v\n", expected, result)
				}
			}
		}
	}
}
