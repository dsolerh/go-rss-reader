package rssreader

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type DefaultTimeFunc func() time.Time

var DefaultTime DefaultTimeFunc = func() time.Time { return time.Time{} }

func Parse(urls ...string) []RSSItem {
	items := make([]RSSItem, 0, len(urls))
	itemsChan := make(chan []RSSItem)
	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, url := range urls {
		go fetchURL(url, itemsChan, &wg)
	}

	var done sync.WaitGroup
	done.Add(1)
	go func() {
		for newItems := range itemsChan {
			items = append(items, newItems...)
		}
		done.Done()
	}()

	wg.Wait()
	itemsChan <- nil
	close(itemsChan)

	done.Wait()
	return items
}

func fetchURL(url string, ch chan<- []RSSItem, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}
	items := parseData(resp.Body, url)
	if len(items) == 0 {
		return
	}

	ch <- items
}

func parseData(data io.Reader, originURL string) []RSSItem {
	var r rss
	if err := xml.NewDecoder(data).Decode(&r); err != nil {
		return []RSSItem{}
	}

	rssItems := make([]RSSItem, 0, len(r.Channel.Items))
	for _, item := range r.Channel.Items {
		if item.Title == "" || item.Link == "" || item.Description == "" {
			continue
		}

		rssItem := RSSItem{
			Title:       item.Title,
			Description: item.Description,
			Link:        item.Link,
		}

		if !item.PubDate.hasValue {
			if DefaultTime == nil {
				continue
			} else {
				rssItem.PublishDate = DefaultTime()
			}
		} else {
			rssItem.PublishDate = item.PubDate.value
		}

		if item.Source != nil {
			rssItem.Source = item.Source.Value
			rssItem.SourceURL = item.Source.URL
		} else {
			host := extractSource(originURL)
			if host == "" {
				continue
			}
			rssItem.Source = host
			rssItem.SourceURL = originURL

		}

		rssItems = append(rssItems, rssItem)
	}
	return rssItems
}

func extractSource(urlRaw string) string {
	u, err := url.Parse(urlRaw)
	if err != nil {
		return ""
	}

	return u.Hostname()
}
