/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package cmd

import (
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/mikerowehl/feeder/internal/output"
	"github.com/mikerowehl/feeder/internal/repository"
	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Write a page with all unread items",
	Long: `Searches through the local database for any items not yet marked as read (so
the feeds must have already been pulled with fetch) and writes out a single
page in the current directory with a table of all the unread items.`,
	Run: func(cmd *cobra.Command, args []string) {
		r, err := repository.NewFeedRepository(dbFile)
		if err != nil {
			log.Fatalf("Error setting up database: %v", err)
		}
		defer r.Close()
		unread, err := r.Unread()
		if err != nil {
			log.Fatalf("Error fetching feeds: %v", err)
		}
		tmpl, err := template.ParseFiles("templates/feed.html")
		if err != nil {
			log.Fatalf("Error opening template: %v", err)
		}
		outFile, err := os.Create("feeder.html")
		if err != nil {
			log.Fatalf("Error writing to feeder.html: %v", err)
		}
		defer outFile.Close()
		err = tmpl.Execute(outFile, output.SanitizeFeeds(unread))
		if err != nil {
			log.Fatalf("Error executing template: %v", err)
		}
		fmt.Println("read called")
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
