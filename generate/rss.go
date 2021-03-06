package generate

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// FeedLink ...
type FeedLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

// FeedCategory ...
type FeedCategory struct {
	Term string `xml:"term,attr"`
}

// FeedEntry ...
type FeedEntry struct {
	ID       string       `xml:"id"`
	Title    string       `xml:"title"`
	Updated  time.Time    `xml:"updated"`
	Link     []FeedLink   `xml:"link"`
	Category FeedCategory `xml:"category"`
}

// FeedAuthor ...
type FeedAuthor struct {
	Name string `xml:"name"`
	URI  string `xml:"uri"`
}

type feed struct {
	NS      string      `xml:"xmlns,attr"`
	ID      string      `xml:"id"`
	Title   string      `xml:"title"`
	Updated time.Time   `xml:"updated"`
	Author  FeedAuthor  `xml:"author"`
	Link    []FeedLink  `xml:"link"`
	Entry   []FeedEntry `xml:"entry"`
}

func (ic IndexContent) GenerateFeed() error {
  if len(ic.Current) < 1 {
    return nil
  }

	r, err := os.Create("newfeed.xml")
	if err != nil {
		return fmt.Errorf("generateFeed create newfeed: %w", err)
	}

	_, err = r.Write([]byte(xml.Header))
	if err != nil {
		return fmt.Errorf("generateFeed write header: %w", err)
	}

	entries := []FeedEntry{}
	entries = flatternAndOrder(entries, getItems(ic.Blog, "blog"))
	entries = flatternAndOrder(entries, getItems(ic.Current, "projects/current"))
	entries = flatternAndOrder(entries, getItems(ic.Past, "projects/past"))

	f := feed{
		NS:      "http://www.w3.org/2005/Atom",
		ID:      "https://www.keloran.dev/",
		Title:   "Keloran blog and stuffs",
		Updated: time.Now(),
		Author: FeedAuthor{
			Name: "Keloran",
			URI:  "https://github.com/keloran",
		},
		Link: []FeedLink{
			{
				Href: "https://www.keloran.dev/",
				Rel:  "self",
			},
			{
				Href: "https://keloran.github.io",
				Rel:  "alternate",
			},
		},
		Entry: entries,
	}

	x, err := xml.Marshal(f)
	if err != nil {
		return fmt.Errorf("generateFeed marshal: %w", err)
	}

	_, err = r.Write(x)
	if err != nil {
		return fmt.Errorf("generateFeed write xml: %w", err)
	}

	err = r.Close()
	if err != nil {
		return fmt.Errorf("generateFeed closeFile: %w", err)
	}

	if fileExists("feed.xml") {
		err = os.Remove("feed.xml")
		if err != nil {
			return fmt.Errorf("remove feedxml: %w", err)
		}
	}
	err = os.Rename("newfeed.xml", "feed.xml")
	if err != nil {
		return fmt.Errorf("move feed: %w", err)
	}

	return nil
}

func getItems(items []File, category string) []FeedEntry {
	feedItems := []FeedEntry{}
	for _, item := range items {
		link := strings.Split(item.Path, ".md")

		entry := FeedEntry{
			ID:      fmt.Sprintf("https://www.keloran.dev/%s.html", link[0][2:]),
			Title:   item.Title,
			Updated: item.Info.ModTime(),
			Link: []FeedLink{
				{
					Href: fmt.Sprintf("https://github.com/keloran/keloran.github.io/blob/master/%s%s", category, item.CleanPath),
					Rel:  "alternate",
				},
			},
			Category: FeedCategory{
				Term: category,
			},
		}
		feedItems = append(feedItems, entry)
	}

	return feedItems
}

func orderItems(items []FeedEntry) []FeedEntry {
	if len(items) < 2 {
		return items
	}

	left, right := 0, len(items)-1
	pivot := rand.Int() % len(items)
	items[pivot], items[right] = items[right], items[pivot]

	for i := range items {
		if items[i].Updated.Unix() > items[right].Updated.Unix() {
			items[left], items[i] = items[i], items[left]
			left++
		}
	}
	items[left], items[right] = items[right], items[left]
	orderItems(items[:left])
	orderItems(items[left+1:])

	return items
}

func flatternAndOrder(flat []FeedEntry, items []FeedEntry) []FeedEntry {
	flat = append(flat, items...)
	flat = orderItems(flat)

	return flat
}
