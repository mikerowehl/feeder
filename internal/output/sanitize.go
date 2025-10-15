/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package output

import (
	"net/url"

	"github.com/mikerowehl/feeder/internal/rss"
)

// Ensure a url is actually just a parsable http or https url, don't allow
// anything else. If this isn't a valid url return just a hash character, so a
// link will just go nowhere.
func SafeURL(s string) string {
	u, err := url.Parse(s)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return "#"
	}
	return u.String()
}

func SanitizeItems(raw []rss.Item) []rss.Item {
	var sanitizedItems []rss.Item
	for _, rawItem := range raw {
		sanitizedItem := rss.Item{
			Title: rawItem.Title,
			Link:  SafeURL(rawItem.Link),
		}
		sanitizedItems = append(sanitizedItems, sanitizedItem)
	}
	return sanitizedItems
}

func SanitizeFeeds(raw []rss.Feed) []rss.Feed {
	var sanitizedFeeds []rss.Feed
	for _, rawFeed := range raw {
		sanitizedFeed := rss.Feed{
			Title: rawFeed.Title,
			URL:   SafeURL(rawFeed.URL),
			Items: SanitizeItems(rawFeed.Items),
		}
		sanitizedFeeds = append(sanitizedFeeds, sanitizedFeed)
	}
	return sanitizedFeeds
}
