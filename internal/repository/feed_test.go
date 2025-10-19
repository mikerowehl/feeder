/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package repository_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/mikerowehl/feeder/internal/repository"
	"github.com/mikerowehl/feeder/internal/rss"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The sqlite driver uses this special name to swap to an in-memory version
// that's both really fast, and that goes away automatically after the tests
// run. So we can avoid cleanup or startup hooks. If you want to see the data
// from a test you can always just make this a normal filename and it'll write
// the file out to make it easy to debug.
var dbFilename = ":memory:"

func setupRepository(t *testing.T) *repository.FeedRepository {
	t.Helper()
	r, err := repository.NewFeedRepository(dbFilename)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	t.Cleanup(func() {
		err := r.Close()
		if err != nil {
			t.Fatalf("failed to close respository: %v", err)
		}
	})
	return r
}

func TestRepository_BasicSaveAndLoad(t *testing.T) {
	r := setupRepository(t)
	feedUrl := "https://test.com/sample.rss"
	testFeed := rss.Feed{URL: feedUrl}
	err := r.Save(&testFeed)
	require.NoError(t, err)
	fetchedFeeds, err := r.All()
	require.NoError(t, err)
	require.Len(t, fetchedFeeds, 1)
	feed1 := fetchedFeeds[0]
	assert.Equal(t, feedUrl, feed1.URL)
}

func TestRepository_UniqueURLViolation(t *testing.T) {
	r := setupRepository(t)
	feedUrl := "https://test.com/sample.rss"
	testFeed1 := rss.Feed{URL: feedUrl}
	err := r.Save(&testFeed1)
	require.NoError(t, err)
	testFeed2 := rss.Feed{URL: feedUrl}
	err = r.Save(&testFeed2)
	require.Error(t, err)
	// Unfortunately the sqlite driver doesn't return the nice duplicate key
	// GORM level error, so check the text.
	assert.Contains(t, err.Error(), "UNIQUE constraint failed")
}

func TestRepository_FeedWithItems(t *testing.T) {
	r := setupRepository(t)
	feeds := []rss.Feed{
		{Title: "Test Feed 1",
			URL:   "https://test.com/sample.rss",
			Items: []rss.Item{{GUID: "1", Content: "test item 1"}},
		}}
	err := r.Save(&(feeds[0]))
	require.NoError(t, err)
	fetched, err := r.All()
	require.NoError(t, err)
	diff := cmp.Diff(feeds, fetched)
	if diff != "" {
		t.Errorf("Mismatch single feed:\n%s", diff)
	}
}

func TestRepository_MultipleFeeds(t *testing.T) {
	r := setupRepository(t)
	feeds := []rss.Feed{
		{Title: "Feed 1", URL: "https://example.com/feed1.rss", Items: []rss.Item{
			{Title: "Feed 1 Item 1",
				Link:    "https://feed1.com/i1",
				Content: "content for 1/1",
				GUID:    "guid1"},
			{Title: "Feed 1 Item 2",
				Link:    "https://feed1.com/i2",
				Content: "content for 1/2",
				GUID:    "guid2"},
		}},
		{Title: "Feed 2", URL: "https://example.com/feed2.rss", Items: []rss.Item{
			{Title: "Feed 2 Item 1",
				Link:    "https://feed2.com/i1",
				Content: "content for 2/1",
				GUID:    "guid10"},
			{Title: "Feed 2 Item 2",
				Link:    "https://feed2.com/i2",
				Content: "content for 2/2",
				GUID:    "guid11"},
		}},
	}
	for _, feed := range feeds {
		err := r.Save(&feed)
		require.NoError(t, err)
	}
	fetched, err := r.All()
	require.NoError(t, err)
	diff := cmp.Diff(feeds, fetched,
		cmpopts.IgnoreFields(rss.Feed{}, "ID", "CreatedAt", "UpdatedAt", "DeletedAt"),
		cmpopts.IgnoreFields(rss.Item{}, "ID", "CreatedAt", "UpdatedAt", "DeletedAt", "FeedID"),
		cmpopts.SortSlices(func(a, b rss.Feed) bool {
			return a.Title < b.Title
		}),
		cmpopts.SortSlices(func(a, b rss.Item) bool {
			return a.Title < b.Title
		}),
	)
	if diff != "" {
		t.Errorf("Mismatch multiple feeds:\n%s", diff)
	}
}

func TestRepository_Unread(t *testing.T) {
	r := setupRepository(t)
	feedUrl := "https://test.com/sample.rss"
	testFeed := rss.Feed{URL: feedUrl}
	err := r.Save(&testFeed)
	require.NoError(t, err)
	feedId := testFeed.ID
	testItem1 := rss.Item{FeedID: feedId, GUID: "1", Content: "test item 1", Read: true}
	testItem2 := rss.Item{FeedID: feedId, GUID: "2", Content: "test item 2", Read: false}
	testFeed.Items = append(testFeed.Items, testItem1, testItem2)
	err = r.Save(&testFeed)
	require.NoError(t, err)
	unread, err := r.Unread()
	require.NoError(t, err)
	require.Len(t, unread, 1)
	require.Len(t, unread[0].Items, 1)
	assert.Equal(t, "2", unread[0].Items[0].GUID)
}

func TestRepository_MarkAll(t *testing.T) {
	r := setupRepository(t)
	feeds := []rss.Feed{
		{Title: "Test Feed 1",
			URL:   "https://test.com/sample.rss",
			Items: []rss.Item{{GUID: "1", Content: "test item 1"}},
		}}
	err := r.Save(&(feeds[0]))
	require.NoError(t, err)
	err = r.MarkAll()
	require.NoError(t, err)
	fetched, err := r.Unread()
	require.NoError(t, err)
	require.Len(t, fetched, 1)
	require.Len(t, fetched[0].Items, 0)
}
