/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package feeder

import (
	"bufio"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mikerowehl/feeder/internal/output"
	"github.com/mikerowehl/feeder/internal/repository"
	"github.com/mikerowehl/feeder/internal/rss"
)

type Feeder struct {
	Db     *repository.FeedRepository
	Client *http.Client
	out    io.Writer
	err    io.Writer
	in     io.Reader
}

//go:embed templates/feed.html
var feedTemplate string

const appName = "feeder"
const maxItems = 100

func NewFeeder(dbFile string, cmdOut io.Writer, cmdErr io.Writer, cmdIn io.Reader) (*Feeder, error) {
	f := &Feeder{}
	r, err := repository.NewFeedRepository(dbFile)
	if err != nil {
		return f, err
	}
	f.Db = r
	f.Client = &http.Client{Timeout: 30 * time.Second}
	f.out = cmdOut
	f.err = cmdErr
	f.in = cmdIn
	return f, nil
}

func TodayFile() string {
	return fmt.Sprintf("%s-%s.html", appName, time.Now().Format(time.DateOnly))
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

func (f *Feeder) Delete(id uint) error {
	return f.Db.Delete(id)
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
		fmt.Fprintln(f.out, "Fetched:", feed.URL)
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
	feeds, err := f.Db.AllFeeds()
	if err != nil {
		return fmt.Errorf("Error fetching feeds: %w", err)
	}
	for i := range feeds {
		feed := &feeds[i]
		fmt.Fprintf(f.out, "%d: %s (%s)\n", feed.ID, feed.Title, feed.URL)
	}
	return nil
}

func (f *Feeder) MarkAll() error {
	return f.Db.MarkAll()
}

func (f *Feeder) Open(filename string) error {
	openPath, err := exec.LookPath("open")
	if err == nil {
		cmd := exec.Command(openPath, filename)
		return cmd.Run()
	}
	xdgOpenPath, err := exec.LookPath("xdg-open")
	if err == nil {
		cmd := exec.Command(xdgOpenPath, filename)
		return cmd.Run()
	}
	return fmt.Errorf("unable find suitable open command")
}

func (f *Feeder) Export() error {
	feeds, err := f.Db.AllFeeds()
	if err != nil {
		return fmt.Errorf("Error fetching feeds: %w", err)
	}
	for i := range feeds {
		feed := &feeds[i]
		fmt.Fprintf(f.out, "%s\n", feed.URL)
	}
	return nil
}

func (f *Feeder) Import() error {
	scanner := bufio.NewScanner(f.in)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" {
			continue
		}
		if err := f.Add(url); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func (f *Feeder) Trim(maxItems int) error {
	feeds, err := f.Db.AllFeeds()
	if err != nil {
		return fmt.Errorf("error reading feeds: %w", err)
	}
	for i := range feeds {
		feed := &feeds[i]
		fmt.Fprintf(f.out, "Trimming feed %s\n", feed.URL)
		err := f.Db.TrimItems(feed.ID, maxItems)
		if err != nil {
			fmt.Fprintf(f.out, "  Error trimming feed %v", err)
		}
	}
	fmt.Fprintf(f.out, "Cleaning up database file\n")
	return f.Db.Vacuum()
}
