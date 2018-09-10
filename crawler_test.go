package main

import (
	"io"
	"net/http"
	"os"
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
<p><a href="http://localhost:8080/def.html">def</a></p>
</body>
</html>
`

var def = `
<html>
<title>DEF</title>
<body>
<p><a href="http://localhost:8080/abc.html">abc</a></p>
<p><a href="http://localhost:8080/ghi.html">ghi</a></p>
</body>
</html>
`

var ghi = `
<html>
<title>GHI</title>
<body>
<p><a href="http://localhost:8080/jkl.html">jkl</a></p>
</body>
</html>
`

var jkl = `
<html>
<title>JKL</title>
</html>
`

func Example() {
	os.Args = []string{"spinarago", "-hostname", "http://example.com"}
	main()
	// output:
	// [{"url":"http://example.com","links":["http://www.iana.org/domains/example"]}]
}

func TestMain(m *testing.M) {
	// Setup Test web server on localhost:8080
	hndlFunc := func(path, body string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != path {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, body)
		}
	}

	http.HandleFunc("/", hndlFunc("/", home))
	http.HandleFunc("/abc.html", hndlFunc("/abc.html", abc))
	http.HandleFunc("/def.html", hndlFunc("/def.html", def))
	http.HandleFunc("/ghi.html", hndlFunc("/ghi.html", ghi))
	http.HandleFunc("/jkl.html", hndlFunc("/jkl.html", jkl))

	go http.ListenAndServe(":8080", nil)
	os.Exit(m.Run())
}

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

func TestFilterByHostname(t *testing.T) {
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

	result := FilterByHostname(domain, input)
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
		},
		"http://localhost:8080/abc.html": {
			"http://localhost:8080/def.html",
		},
		// Test for cycles (should not visit abc again)
		"http://localhost:8080/def.html": {
			"http://localhost:8080/abc.html",
			"http://localhost:8080/ghi.html",
		},
		// Test for level (at level = 3, jkl should not appear as a key)
		"http://localhost:8080/ghi.html": {
			"http://localhost:8080/jkl.html",
		},
	}

	result, err := Crawl(input, 3, 10, false)
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
