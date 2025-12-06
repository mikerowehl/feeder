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

// markCmd represents the mark command
var markCmd = &cobra.Command{
	Use:   "mark",
	Short: "Marks all of the items in the database as read",
	Long: `By default all of the items added to the database are marked as unread so they'll go
into the next page of posts. This command marks all of the items as read. This
is called automatically as part of the daily command. If you're using each 
command individually it should be called after calling read to generate a page
of posts to read:

feeder fetch # retreive the latest from the feed
feeder read  # generate a local file with all the posts to read
feeder mark  # mark everything in the database as read
feeder open  # open the genereated file in your default browser`,
	RunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Context().Value(feederKey).(*feeder.Feeder)
		err := f.MarkAll()
		if err != nil {
			return fmt.Errorf("Error marking feeds: %w", err)
		}
		fmt.Println("mark called")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(markCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// markCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// markCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
