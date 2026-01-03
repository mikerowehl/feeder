/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package rss_test

import (
	"context"
	"sort"
	"testing"

	"github.com/mikerowehl/feeder/internal/rss"
	"github.com/mikerowehl/feeder/test/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var basicFeed = `<?xml version="1.0" encoding="UTF-8" ?>
<rss version="2.0">
  <channel>
    <title>Simple RSS Feed</title>
    <link>https://example.com/rss.xml</link>
    <description>A minimal example of an RSS feed</description>
    <item>
      <title>First Post</title>
      <link>https://example.com/post1</link>
      <description>This is the first post in the feed.</description>
	  <pubDate>Mon, 03 Nov 2025 12:00:00 GMT</pubDate>
    </item>
    <item>
      <title>Second Post</title>
      <link>https://example.com/post2</link>
      <description>This is the second post in the feed.</description>
	  <pubDate>Sun, 02 Nov 2025 12:00:00 GMT</pubDate>
    </item>
  </channel>
</rss>`

func DateSortItems(items []rss.Item) {
	sort.Slice(items, func(a, b int) bool {
		return items[a].Published.Before(items[b].Published)
	})
}

func TestFeed_EmptyFeed(t *testing.T) {
	feed := rss.Feed{}
	err := feed.Process(basicFeed, 25)
	require.NoError(t, err)
	assert.Len(t, feed.Items, 2)
}

func TestFeed_FetchSimple(t *testing.T) {
	client := mock.NewMockClient(basicFeed, 200)
	feed := rss.Feed{URL: "https://testing.com/dummyfeed.rss"}
	err := feed.Fetch(client, 25)
	require.NoError(t, err)
	assert.Len(t, feed.Items, 2)
	DateSortItems(feed.Items)
	firstItem := feed.Items[0]
	require.Equal(t, "Second Post", firstItem.Title)
}

func TestFeed_FetchInvalidRSS(t *testing.T) {
	client := mock.NewMockClient("This isn't a feed", 200)
	feed := rss.Feed{URL: "https://testing.com/dummyfeed.rss"}
	err := feed.Fetch(client, 25)
	require.Error(t, err)
}

func TestFeed_FetchNetworkError(t *testing.T) {
	client := mock.NewMockClientWithError(context.DeadlineExceeded)
	feed := rss.Feed{URL: "https://testing.com/dummyfeed.rss"}
	err := feed.Fetch(client, 25)
	require.Error(t, err)
}
