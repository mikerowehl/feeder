/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package cmd

import (
	"fmt"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		feedUrl := args[0]
		fmt.Println("Adding feed:", feedUrl)
		r, err := repository.NewFeedRepository(dbFile)
		if err != nil {
			return fmt.Errorf("error setting up database: %w", err)
		}
		defer r.Close()
		client := &http.Client{}
		feed, err := rss.FeedFromURL(feedUrl, client)
		if err != nil {
			return fmt.Errorf("error creating feed from url: %w", err)
		}
		err = r.Save(&feed)
		if err != nil {
			return fmt.Errorf("error adding feed: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
