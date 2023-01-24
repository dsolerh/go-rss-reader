package rssreader

import (
	"encoding/xml"
	"time"

	"github.com/araddon/dateparse"
)

// rss represents the shcema of an RSS feed
type rss struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel channel  `xml:"channel"`
}

// in the oficial shcema channel contains more than just `item`
// but there is no need to use those fields
type channel struct {
	XMLName xml.Name `xml:"channel"`
	Items   []item   `xml:"item"`
}

// item represent the actual feed for each news
type item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     xmlTime  `xml:"pubDate"`
	Source      *source  `xml:"source"`
}

// this is for custom unmarshaling of date
type xmlTime struct {
	value time.Time
}

func (t *xmlTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var dateString string
	d.DecodeElement(&dateString, &start)
	dateTime, err := dateparse.ParseAny(dateString)
	if err != nil {
		return err
	}
	*t = xmlTime{value: dateTime}
	return nil
}

// this represents a cource tag
type source struct {
	XMLName xml.Name `xml:"source"`
	URL     string   `xml:"url,attr"`
	Value   string   `xml:",chardata"`
}
