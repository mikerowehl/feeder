/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package rss

import (
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/mmcdole/gofeed"
	"gorm.io/gorm"
)

type Feed struct {
	gorm.Model
	URL   string `gorm:"unique"`
	Title string
	Items []Item
}

type Item struct {
	gorm.Model
	FeedID  uint
	Title   string
	Link    string
	Content string
	GUID    string `gorm:"unique"`
	Read    bool
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected http status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// Initialy fetch a feed given a URL. Updates just the metadata necessary to
// make the feed itself. Doesn't process items in the feed. The feed at this
// point doesn't have a database ID, so we wouldn't be able to make those
// child items and link them to the parent yet.
func FeedFromURL(url string, client *http.Client) (Feed, error) {
	feed := Feed{URL: url}
	content, err := FetchFeedContent(url, client)
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
	return Item{
		Title:   parsed.Title,
		Link:    parsed.Link,
		Content: content,
		GUID:    guid,
		Read:    false,
	}
}

func (feed *Feed) Fetch(client *http.Client) error {
	content, err := FetchFeedContent(feed.URL, client)
	if err != nil {
		return err
	}

	return feed.Process(content)
}

// Process the current content of the feed and parse into items. If there are
// already items in the list attached to the feed we only create new items for
// the entries we don't have. New items are populated with Read set to false.
func (feed *Feed) Process(content string) error {
	fp := gofeed.NewParser()
	parsed, err := fp.ParseString(content)
	if err != nil {
		return err
	}
	for _, parsedItem := range parsed.Items {
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
