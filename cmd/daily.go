/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/mikerowehl/feeder/internal/feeder"
	"github.com/spf13/cobra"
)

// dailyCmd represents the daily command
var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Fetches all the feeds, makes a page of posts, and marks all read",
	Long: `Just a convenience wrapper around fetch, read, and mark. Just checks at each
operation and only goes to the next if everything is okay.`,
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
		outfile := fmt.Sprintf("feeder-%s.html", time.Now().Format(time.DateOnly))
		err = f.WriteUnread(outfile)
		if err != nil {
			log.Fatalf("Error writing out unread: %v", err)
		}
		err = f.MarkAll()
		if err != nil {
			log.Fatalf("Error marking feeds: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(dailyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dailyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dailyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
