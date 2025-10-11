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

func TestRepository_BasicSaveAndLoad(t *testing.T) {
	r, err := repository.NewFeedRepository(dbFilename)
	require.NoError(t, err)
	feedUrl := "https://test.com/sample.rss"
	testFeed := rss.Feed{URL: feedUrl}
	err = r.Save(&testFeed)
	require.NoError(t, err)
	fetchedFeeds, err := r.All()
	require.Len(t, fetchedFeeds, 1)
	feed1 := fetchedFeeds[0]
	assert.Equal(t, feedUrl, feed1.URL)
}
