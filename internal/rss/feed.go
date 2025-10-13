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

// TODO put in etag and modified check too
func (feed *Feed) Fetch(client *http.Client) error {
	req, err := http.NewRequest("GET", feed.URL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/rss+xml, application/atom+xml, application/xml;q=0.9, text/xml;q=0.8, */*;q=0.7")
	req.Header.Set("User-Agent", "Feeder/0.0 (+https://github.com/mikerowehl/feeder)")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected http status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return feed.Process(string(body))
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
	for _, i := range parsed.Items {
		found := slices.IndexFunc(feed.Items, func(search Item) bool {
			return search.GUID == i.GUID
		})
		if found == -1 {
			content := i.Content
			if content == "" {
				content = i.Description
			}
			guid := i.GUID
			if guid == "" {
				guid = i.Link
			}
			feed.Items = append(feed.Items, Item{
				FeedID:  feed.ID,
				Title:   i.Title,
				Link:    i.Link,
				Content: content,
				GUID:    guid,
				Read:    false,
			})
		}
	}
	return nil
}
