# go-rss-reader

![tests](https://github.com/dsolerh/go-rss-reader/actions/workflows/test.yml/badge.svg)

The `go-rss-reader` package helps manage HTTP request for RSS feeds.

The principle is simple, it provides a `Parse` function who makes the requests to the 
provided URLs (concurrently for efficiency), then process each request using this 
[XML schema](https://www.w3schools.com/xml/xml_rss.asp#rssref) for RSS feeds. 
The response is a slice of `RSSItem` which is a struct that contains the information for:

1. Title
2. Link
3. Description
4. Publish Date (Optional)
5. Source (Optional)
6. Source URL (Optional)

The dedinition for the `Parse` function and `RSSItem` can be found 
[here](https://github.com/dsolerh/go-rss-reader/blob/main/parser.go) 
and [here](https://github.com/dsolerh/go-rss-reader/blob/main/item.go) 
respectively.

## Install

`go get -u github.com/dsolerh/go-rss-reader`

## Examples

**As easy as:**

```go
package main
import (
    fmt

	reader "github.com/dsolerh/go-rss-reader"
)
func main() {
	urls := []string{
		"http://yournews.com/rss",
		"http://abc.com/feed",
	}
	items := reader.Parse(urls...)

    fmt.Println(items)
}
```

As a result you should see a list of `RSSItem` so long as the given URLs are
correct and they contain valid xml feeds.

## License

Copyright (c) 2023-present [Daniel Soler](https://github.com/dsolerh)

Licensed under [MIT License](./LICENSE)
