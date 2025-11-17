/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package cmd

import (
	"fmt"

	"github.com/mikerowehl/feeder/internal/feeder"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch the content from feeds and update the local set of items",
	Long: `For the set of feeds in the local database this fetches the content from
each of the URLs and updates the items associated with the feed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := feeder.NewFeeder(dbFile)
		if err != nil {
			return fmt.Errorf("startup error: %w", err)
		}
		defer f.Close()
		err = f.Fetch()
		if err != nil {
			return fmt.Errorf("error fetching feeds: %w", err)
		}
		fmt.Println("fetch finished")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
