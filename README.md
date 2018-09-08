# spinarago

A basic web crawler written in Go.

## Why "Spinarago"?

Spinarak = A spider pokemón  
Go = Go

You do the math. :)

## Description

This is a basic web crawler that prints the site map of given URL. It does print external URLs, but doesn't follow them. This project is still a work in progress so
it's not feature complete. Feel free to make suggestions for improvements, either by creating issues or submiting pull requests.

## Install

```sh
$ go get github.com/danicat/spinarago
```

## Usage

```sh
$ spinarago --hostname <host> --delay <milliseconds> --level <max-depth>
```

I highly recommend for you to install [jq](https://stedolan.github.io/jq/) to pretty print the json output. Example:

```sh
$ spinarago --hostname http://example.com | jq
```

You can also redirect the stdout to a json file to make a site map dump:

```sh
$ ./spinarago --hostname http://example.com -level 1 -delay 10 > example_level1.json
```

`jq` is really handy to filter the output:

```sh
$ cat example_level1.json | jq '.[] | { url: .url }'
```

## TODO

- Handle relative paths
- Refactor tests

## Contributing

I'm open to contributions. Just create an issue and/or submit a pull request.

## Contact

Any comments please feel free to reach out to me at @danicat83 on Twitter.