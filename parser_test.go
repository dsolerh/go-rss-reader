package rssreader

import (
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/araddon/dateparse"
)

func TestParse(t *testing.T) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("id") {
		case "valid":
			w.Write([]byte(schema_valid))
		case "valid-date":
			w.Write([]byte(schema_valid_date))
		case "valid-source":
			w.Write([]byte(schema_valid_source))
		case "valid-no-item":
			w.Write([]byte(schema_valid_no_item))
		case "multi-items":
			w.Write([]byte(schema_multi_items))
		case "invalid":
			w.Write([]byte(schema_invalid))
		case "invalid-date":
			w.Write([]byte(schema_invalid_date))
		}
	})
	go func() {
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			t.Error("couldn't setup the test")
		}
	}()

	urls := []string{
		"http://localhost:3000?id=valid",
		"http://localhost:3000?id=valid-date",
		"http://localhost:3000?id=valid-source",
		"http://localhost:3000?id=valid-no-item",
		"http://localhost:3000?id=multi-items",
		"http://localhost:3000?id=invalid",
		"http://localhost:3000?id=invalid-date",
	}

	items := Parse(urls...)

	if len(items) != 5 {
		t.Errorf("the amount of items should be 5 but is %d", len(items))
	}
}

func Test_parseData(t *testing.T) {
	testCases := []struct {
		desc       string
		xml        string
		itemLength int
		items      []RSSItem
	}{
		{
			desc:       "a valid schema",
			xml:        schema_valid,
			itemLength: 1,
			items: []RSSItem{
				{Title: "RSS Tutorial", Link: "https://www.w3schools.com/xml/xml_rss.asp", Description: "New RSS tutorial on W3Schools"},
			},
		},
		{
			desc:       "a valid schema (2 items)",
			xml:        schema_multi_items,
			itemLength: 2,
			items: []RSSItem{
				{Title: "RSS Tutorial", Link: "https://www.w3schools.com/xml/xml_rss.asp", Description: "New RSS tutorial on W3Schools"},
				{Title: "RSS Tutorial", Link: "https://www.w3schools.com/xml/xml_rss.asp", Description: "New RSS tutorial on W3Schools"},
			},
		},
		{
			desc:       "a valid schema (source)",
			xml:        schema_valid_source,
			itemLength: 1,
			items: []RSSItem{
				{
					Title:       "RSS Tutorial",
					Link:        "https://www.w3schools.com/xml/xml_rss.asp",
					Description: "New RSS tutorial on W3Schools",
					Source:      "W3Schools.com",
					SourceURL:   "https://www.w3schools.com",
				},
			},
		},
		{
			desc:       "a valid schema (date)",
			xml:        schema_valid_date,
			itemLength: 1,
			items: []RSSItem{
				{
					Title:       "RSS Tutorial",
					Link:        "https://www.w3schools.com/xml/xml_rss.asp",
					Description: "New RSS tutorial on W3Schools",
					PublishDate: dateparse.MustParse("Thu, 27 Apr 2006"),
				},
			},
		},
		{
			desc:       "an invalid schema",
			xml:        schema_invalid,
			itemLength: 0,
			items:      []RSSItem{},
		},
		{
			desc:       "an invalid schema (date)",
			xml:        schema_invalid_date,
			itemLength: 0,
			items:      []RSSItem{},
		},
		{
			desc:       "an valid schema (no item)",
			xml:        schema_invalid,
			itemLength: 0,
			items:      []RSSItem{},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			data := strings.NewReader(tC.xml)
			items := parseData(data)
			if len(items) != tC.itemLength {
				t.Errorf("the items parsed should be %d but are %d", tC.itemLength, len(items))
			}
			if !reflect.DeepEqual(items, tC.items) {
				t.Errorf("the items should be %v", tC.items)
			}
		})
	}
}

var schema_valid = `
<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">

<channel>
  <title>W3Schools Home Page</title>
  <link>https://www.w3schools.com</link>
  <description>Free web building tutorials</description>
  <item>
    <title>RSS Tutorial</title>
    <link>https://www.w3schools.com/xml/xml_rss.asp</link>
    <description>New RSS tutorial on W3Schools</description>
  </item>
</channel>

</rss>
`

var schema_multi_items = `
<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">

<channel>
  <title>W3Schools Home Page</title>
  <link>https://www.w3schools.com</link>
  <description>Free web building tutorials</description>
  <item>
    <title>RSS Tutorial</title>
    <link>https://www.w3schools.com/xml/xml_rss.asp</link>
    <description>New RSS tutorial on W3Schools</description>
  </item>
  <item>
    <title>RSS Tutorial</title>
    <link>https://www.w3schools.com/xml/xml_rss.asp</link>
    <description>New RSS tutorial on W3Schools</description>
  </item>
</channel>

</rss>
`

var schema_valid_source = `
<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">

<channel>
  <title>W3Schools Home Page</title>
  <link>https://www.w3schools.com</link>
  <description>Free web building tutorials</description>
  <item>
    <title>RSS Tutorial</title>
    <link>https://www.w3schools.com/xml/xml_rss.asp</link>
    <description>New RSS tutorial on W3Schools</description>
    <source url="https://www.w3schools.com">W3Schools.com</source>
  </item>
</channel>

</rss>
`
var schema_valid_date = `
<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">

<channel>
  <title>W3Schools Home Page</title>
  <link>https://www.w3schools.com</link>
  <description>Free web building tutorials</description>
  <item>
    <title>RSS Tutorial</title>
    <link>https://www.w3schools.com/xml/xml_rss.asp</link>
    <description>New RSS tutorial on W3Schools</description>
    <pubDate>Thu, 27 Apr 2006</pubDate>
  </item>
</channel>

</rss>
`
var schema_invalid = ""

var schema_invalid_date = `
<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">

<channel>
  <title>W3Schools Home Page</title>
  <link>https://www.w3schools.com</link>
  <description>Free web building tutorials</description>
  <item>
    <title>RSS Tutorial</title>
    <link>https://www.w3schools.com/xml/xml_rss.asp</link>
    <description>New RSS tutorial on W3Schools</description>
    <pubDate></pubDate>
  </item>
</channel>

</rss>
`

var schema_valid_no_item = `
<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">

<channel>
</channel>

</rss>
`
