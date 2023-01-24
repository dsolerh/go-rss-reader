package rssreader

import (
	"net/http"
)

func Parse(urls ...string) []RSSItem {
	return nil
}

func fetchURL(url string, ch chan<- *RSSItem) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}
	item, err := parseItem(resp.Body)
	if err != nil {
		return
	}

	ch <- item
}
