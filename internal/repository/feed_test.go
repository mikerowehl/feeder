/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package repository_test

import (
	"testing"

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
	feedUrl := "https://test.com/sample.rss"
	testFeed := rss.Feed{URL: feedUrl}
	err := r.Save(&testFeed)
	require.NoError(t, err)
	feedId := testFeed.ID
	testItem1 := rss.Item{FeedID: feedId, GUID: "1", Content: "test item 1"}
	testItem2 := rss.Item{FeedID: feedId, GUID: "2", Content: "test item 2"}
	testFeed.Items = append(testFeed.Items, testItem1, testItem2)
	err = r.Save(&testFeed)
	require.NoError(t, err)
	fetchedFeeds, err := r.All()
	require.NoError(t, err)
	require.Len(t, fetchedFeeds, 1)
	require.Len(t, fetchedFeeds[0].Items, 2)
}
