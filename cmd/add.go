/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mikerowehl/feeder/internal/repository"
	"github.com/mikerowehl/feeder/internal/rss"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add URL",
	Short: "Add a URL to the list of feeds",
	Long: `Adds the URL given to the list in the database. The URL should be
passed directly as an argument. Make sure the URL is a full URL with all 
components, and note that you might need to quote the URL depending on the 
contents of the URL and what environment you're running in.

ex: feeder add "https://rowehl.com/feed.xml"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		feedUrl := args[0]
		fmt.Println("Adding feed:", feedUrl)
		r, err := repository.NewFeedRepository(dbFile)
		if err != nil {
			log.Fatalf("Error setting up database: %v", err)
		}
		defer r.Close()
		client := &http.Client{}
		feed, err := rss.FeedFromURL(feedUrl, client)
		if err != nil {
			log.Fatalf("Error creating feed from url: %v", err)
		}
		err = r.Save(&feed)
		if err != nil {
			log.Fatalf("Error adding feed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
