/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package feeder

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mikerowehl/feeder/internal/repository"
)

type Feeder struct {
	Db     *repository.FeedRepository
	Client *http.Client
}

func NewFeeder(dbFile string) (Feeder, error) {
	f := Feeder{}
	r, err := repository.NewFeedRepository(dbFile)
	if err != nil {
		return f, err
	}
	f.Db = r
	f.Client = &http.Client{}
	return f, nil
}

func (f Feeder) Close() {
	if f.Db != nil {
		f.Db.Close()
	}
}

func (f Feeder) Fetch() error {
	feeds, err := f.Db.All()
	if err != nil {
		return fmt.Errorf("Error fetching feeds: %W", err)
	}
	for _, feed := range feeds {
		err := feed.Fetch(f.Client)
		if err != nil {
			return fmt.Errorf("Error fetching feed %s: %w", feed.URL, err)
		}
		err = f.Db.Save(&feed)
		if err != nil {
			return fmt.Errorf("Error saving feed %s: %w", feed.URL, err)
		}
		log.Println("Fetched:", feed.URL)
	}
	return nil
}
