package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Feed struct {
	gorm.Model
	URL   string
	Items []Item
}

type Item struct {
	gorm.Model
	FeedID    uint
	Title     string
	Link      string
	Content   string
	GUID      string `gorm:"unique"`
	State     string
	Published time.Time
}

type ItemListPage struct {
	Items []Item
}

func handleItemsUnread(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	unreadItem := Item{State: "unread"}
	itemEntries := []Item{}
	res := db.Limit(10).Order("published").Where(&unreadItem).Find(&itemEntries)
	if res.Error != nil {
		fmt.Fprintf(w, "Error reading from DB: %s", res.Error)
		return
	}
	tmpl := template.Must(template.ParseFiles("itemlist.html"))
	page := ItemListPage{
		Items: itemEntries,
	}
	tmpl.Execute(w, page)
}

func processItem(db *gorm.DB, feedEntry *Feed, item *gofeed.Item) (added bool, err error) {
	fmt.Printf("Processing item with GUID %s\n", item.GUID)
	var itemEntry Item
	res := db.Where(&Item{GUID: item.GUID}).First(&itemEntry)
	if res.Error == nil {
		fmt.Printf("Found a matching record for %s\n", item.GUID)
		return
	}
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return
	}
	fmt.Printf("No match, adding to feed\n")
	content := item.Content
	if content == "" {
		content = item.Description
	}
	newItem := Item{
		Title:     item.Title,
		Link:      item.Link,
		Content:   content,
		GUID:      item.GUID,
		State:     "unread",
		Published: *item.PublishedParsed,
	}
	feedEntry.Items = append(feedEntry.Items, newItem)
	added = true
	return
}

func processFeed(db *gorm.DB, feedEntry *Feed) (err error) {
	fmt.Printf("Processing %v\n", feedEntry.URL)
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedEntry.URL)
	if err != nil {
		fmt.Printf("Error fetching feed %s: %v\n", feedEntry.URL, err)
		return
	}
	needSave := false
	for _, i := range feed.Items {
		added, err := processItem(db, feedEntry, i)
		if err != nil {
			return err
		}
		needSave = needSave || added
	}
	if needSave {
		fmt.Printf("Saving feed to database\n")
		db.Save(feedEntry)
	}
	return
}

func pollFeeds(db *gorm.DB) {
	time.Sleep(20 * time.Minute) // Just sleep this so it doesn't happen every run
	for {
		feeds := []Feed{}
		result := db.Find(&feeds)
		if result.Error != nil {
			fmt.Printf("Error querying database: %v", result.Error)
		}
		for _, f := range feeds {
			err := processFeed(db, &f)
			if err != nil {
				fmt.Printf("Error processing %s: %v\n", f.URL, err)
				continue
			}
		}
		time.Sleep(10 * time.Minute)
	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("feeder.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Feed{}, &Item{})
	go pollFeeds(db)
	log.Println("Starting the http server")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleItemsUnread(db, w, r)
	})
	log.Fatal(http.ListenAndServe(":9090", nil))
}