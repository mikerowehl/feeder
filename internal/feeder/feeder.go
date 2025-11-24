/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package feeder

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mikerowehl/feeder/internal/output"
	"github.com/mikerowehl/feeder/internal/repository"
	"github.com/mikerowehl/feeder/internal/rss"
)

type Feeder struct {
	Db     *repository.FeedRepository
	Client *http.Client
}

//go:embed templates/feed.html
var feedTemplate string

const maxItems = 100

func NewFeeder(dbFile string) (*Feeder, error) {
	f := &Feeder{}
	r, err := repository.NewFeedRepository(dbFile)
	if err != nil {
		return f, err
	}
	f.Db = r
	f.Client = &http.Client{Timeout: 30 * time.Second}
	return f, nil
}

func (f *Feeder) Close() {
	if f.Db != nil {
		f.Db.Close()
	}
}

func (f *Feeder) Add(url string) error {
	feed, err := rss.FeedFromURL(url, f.Client)
	if err != nil {
		return fmt.Errorf("error creating feed from url %s: %w", url, err)
	}
	err = f.Db.Save(&feed)
	if err != nil {
		return fmt.Errorf("error adding feed: %w", err)
	}
	return nil
}

func (f *Feeder) Fetch() error {
	feeds, err := f.Db.All()
	if err != nil {
		return fmt.Errorf("Error fetching feeds: %w", err)
	}
	for i := range feeds {
		feed := &feeds[i]
		err := feed.Fetch(f.Client, maxItems)
		if err != nil {
			return fmt.Errorf("Error fetching feed %s: %w", feed.URL, err)
		}
		err = f.Db.Save(feed)
		if err != nil {
			return fmt.Errorf("Error saving feed %s: %w", feed.URL, err)
		}
		log.Println("Fetched:", feed.URL)
	}
	return nil
}

func (f *Feeder) WriteUnread(outFilename string) error {
	unread, err := f.Db.Unread()
	if err != nil {
		return fmt.Errorf("Error fetching feeds: %w", err)
	}
	tmpl, err := template.New("feed").Parse(feedTemplate)
	if err != nil {
		return fmt.Errorf("Error opening template: %v", err)
	}
	outFile, err := os.OpenFile(outFilename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("Error opening output file %s: %w", outFilename, err)
	}
	defer outFile.Close()
	err = tmpl.Execute(outFile, output.SanitizeFeeds(unread))
	if err != nil {
		return fmt.Errorf("Error executing template: %w", err)
	}
	return nil
}

func (f *Feeder) List() error {
	feeds, err := f.Db.All()
	if err != nil {
		return fmt.Errorf("Error fetching feeds: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("%d: %s (%s)\n", feed.ID, feed.Title, feed.URL)
	}
	return nil
}

func (f *Feeder) MarkAll() error {
	return f.Db.MarkAll()
}
