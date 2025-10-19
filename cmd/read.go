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

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Write a page with all unread items",
	Long: `Searches through the local database for any items not yet marked as read (so
the feeds must have already been pulled with fetch) and writes out a single
page in the current directory with a table of all the unread items.`,
	Run: func(cmd *cobra.Command, args []string) {
		f, err := feeder.NewFeeder(dbFile)
		if err != nil {
			log.Fatalf("Startup error: %v", err)
		}
		defer f.Close()
		outfile := fmt.Sprintf("feeder-%s.html", time.Now().Format(time.DateOnly))
		err = f.WriteUnread(outfile)
		if err != nil {
			log.Fatalf("Error writing out unread: %v", err)
		}
		fmt.Println("read called")
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}
