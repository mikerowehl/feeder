/*
Copyright (c) Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package rss

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
	"gorm.io/gorm"
)

type Feed struct {
	gorm.Model
	URL   string `gorm:"unique"`
	Title string
	Items []Item `gorm:"constraint:OnDelete:CASCADE;"`
}

type Item struct {
	gorm.Model
	FeedID    uint
	Title     string
	Link      string
	Content   string
	GUID      string `gorm:"unique"`
	Published time.Time
	Read      bool
}

// Makes the web request to fetch the content of the feed, setting headers and
// checking the return. If no error, the returned string is the full content
// of the feed.
// TODO put in etag and modified check
func FetchFeedContent(url string, client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/rss+xml, application/atom+xml, application/xml;q=0.9, text/xml;q=0.8, */*;q=0.7")
	req.Header.Set("User-Agent", "Feeder/0.0 (+https://github.com/mikerowehl/feeder)")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("failed to close feed body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected http status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Given a URL try to figure out if this is a feed URL, and lok up the feed
// URL if this isn't a feed already. We do a HEAD request and check the
// content type returned to try to figure out if this is a feed. If it's not a
// feed, but it's HTML, parse the HTML and look for a feed alternative link in
// the document header and return that.
func GetFeedURL(givenURL string, client *http.Client) (string, error) {
	req, err := http.NewRequest("HEAD", givenURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status for http request %v", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/rss+xml") ||
		strings.Contains(contentType, "application/atom+xml") ||
		strings.Contains(contentType, "application/xml") ||
		strings.Contains(contentType, "text/xml") {
		return givenURL, nil
	}

	if strings.Contains(contentType, "text/html") {
		return DiscoverFeed(givenURL, client)
	}

	return "", fmt.Errorf("unexpected content type: %s", contentType)
}

// Look for an alternative link header in the HTML content of a page. This is
// called after we do a HEAD on the URL given and we know it's HTML, so we
// just need to fetch it and try to parse.
func DiscoverFeed(givenURL string, client *http.Client) (string, error) {
	req, err := http.NewRequest("GET", givenURL, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	feedURL, err := FindFeedLink(resp.Body)
	if err != nil {
		return "", err
	}

	base, err := url.Parse(givenURL)
	if err != nil {
		return "", err
	}
	feed, err := url.Parse(feedURL)
	if err != nil {
		return "", err
	}

	return base.ResolveReference(feed).String(), nil
}

// Parse the content of a page looking for the alternate link.
func FindFeedLink(r io.Reader) (string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", err
	}

	var feedURL string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if feedURL != "" {
			return // Already found
		}

		if n.Type == html.ElementNode && n.Data == "link" {
			var rel, typ, href string
			for _, attr := range n.Attr {
				switch attr.Key {
				case "rel":
					rel = attr.Val
				case "type":
					typ = attr.Val
				case "href":
					href = attr.Val
				}
			}

			// Check if it's an alternate feed link
			if rel == "alternate" &&
				(typ == "application/rss+xml" || typ == "application/atom+xml") &&
				href != "" {
				feedURL = href
				return
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if feedURL == "" {
		return "", fmt.Errorf("no feed found")
	}

	return feedURL, nil
}

// Initially fetch a feed given a URL. Updates just the metadata necessary to
// make the feed itself. Doesn't process items in the feed. But does account
// for the user giving the URL of content that has feed alternative link. So
// first we do a HEAD request and look at the content type. If needed we try
// to determine the feed URL from the content URL. That means the URL that
// ends up in the Feed entry in the DB might not match what the user put in.
func FeedFromURL(url string, client *http.Client) (Feed, error) {
	feedUrl, err := GetFeedURL(url, client)
	if err != nil {
		feedUrl = url
	}
	feed := Feed{URL: feedUrl}
	content, err := FetchFeedContent(feedUrl, client)
	if err != nil {
		return feed, err
	}
	fp := gofeed.NewParser()
	parsed, err := fp.ParseString(content)
	if err != nil {
		return feed, err
	}
	if parsed.Title != "" {
		feed.Title = parsed.Title
	} else {
		feed.Title = url
	}
	return feed, nil
}

// Turn a gofeed version of an item into our item.
func ParsedToItem(parsed *gofeed.Item) Item {
	guid := parsed.GUID
	if guid == "" {
		guid = parsed.Link
	}
	content := parsed.Content
	if content == "" {
		content = parsed.Description
	}
	var published time.Time
	if parsed.PublishedParsed != nil {
		published = *parsed.PublishedParsed
	} else {
		published = time.Now()
	}
	return Item{
		Title:     parsed.Title,
		Link:      parsed.Link,
		Content:   content,
		GUID:      guid,
		Published: published,
		Read:      false,
	}
}

func (feed *Feed) Fetch(client *http.Client, maxItems int) error {
	content, err := FetchFeedContent(feed.URL, client)
	if err != nil {
		return err
	}

	return feed.Process(content, maxItems)
}

// Process the current content of the feed and parse into items. If there are
// already items in the list attached to the feed we only create new items for
// the entries we don't have. New items are populated with Read set to false.
func (feed *Feed) Process(content string, maxItems int) error {
	fp := gofeed.NewParser()
	parsed, err := fp.ParseString(content)
	if err != nil {
		return err
	}
	useItems := parsed.Items
	if len(useItems) > maxItems {
		sort.Sort(parsed)
		useItems = parsed.Items[len(parsed.Items)-maxItems:]
	}
	for _, parsedItem := range useItems {
		item := ParsedToItem(parsedItem)
		found := slices.IndexFunc(feed.Items, func(search Item) bool {
			return search.GUID == item.GUID
		})
		if found == -1 {
			feed.Items = append(feed.Items, item)
		}
	}
	return nil
}
