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

func TestRepository_BasicSaveAndLoad(t *testing.T) {
	r, err := repository.NewFeedRepository("test.db")
	require.NoError(t, err)
	testFeed := rss.Feed{URL: "https://test.com/sample.rss"}
	err = r.Save(&testFeed)
	require.NoError(t, err)
	fetchedFeeds, err := r.All()
	assert.Len(t, fetchedFeeds, 1)
}
