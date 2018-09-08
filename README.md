# spinarago

A basic web crawler written in Go.

## Why "Spinarago"?

Spinarak = A spider pokemón
Go = Go

You do the math. :)

## Install

```sh
$ go get github.com/danicat/spinarago
```

## Usage

```sh
$ ./spinarago --hostname <host> --delay <milliseconds> --level <max-depth>
```

I highly recommend for you to install [jq](https://stedolan.github.io/jq/) to pretty print the json output. Example:

```sh
$ ./spinarago --hostname http://example.com | jq
```

You can also redirect the stdout to a json file to make a site map dump:

```sh
$ ./spinarago --hostname http://example.com -level 1 -delay 10 > example_level1.json
$ cat example_level1.json | jq '.[] | { url: .url }'
```


## TODO

- Handle relative paths
- Refactor tests