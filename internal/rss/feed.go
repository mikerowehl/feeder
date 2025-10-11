package rss

import (
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

// Process the current content of the feed and parse into items. If there are
// already items in the list attached to the feed we only create new items for
// the entries we don't have. New items are populated with Read set to false.
func (feed *Feed) Process(content string) (err error) {
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
