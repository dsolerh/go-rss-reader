package rssreader

import "time"

// RSSItem is the representation of an item
// retrieved from an RSS feed
type RSSItem struct {
	Title       string    // Defines the title of the item
	Source      string    // Specifies a third-party source for the item
	SourceURL   string    // Specifies the link to the source
	Link        string    // Defines the hyperlink to the item
	PublishDate time.Time // Defines the last-publication date for the item
	Description string    // Describes the item
}
