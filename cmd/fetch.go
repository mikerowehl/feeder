/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/mikerowehl/feeder/internal/feeder"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch the content from feeds and update the local set of items",
	Long: `For the set of feeds in the local database this fetches the content from
each of the URLs and updates the items associated with the feed.`,
	Run: func(cmd *cobra.Command, args []string) {
		f, err := feeder.NewFeeder(dbFile)
		if err != nil {
			log.Fatalf("Startup error: %v", err)
		}
		defer f.Close()
		err = f.Fetch()
		if err != nil {
			log.Fatalf("Error fetching feeds: %v", err)
		}
		fmt.Println("fetch finished")
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
