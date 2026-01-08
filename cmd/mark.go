/*
Copyright (c) Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package cmd

import (
	"fmt"

	"github.com/mikerowehl/feeder/internal/feeder"
	"github.com/spf13/cobra"
)

func NewMarkCmd() *cobra.Command {
	markCmd := &cobra.Command{
		Use:   "mark",
		Short: "Marks all of the items in the database as read",
		Long: `By default all of the items added to the database are marked as unread so they'll go
into the next page of posts. This command marks all of the items as read. This
is called automatically as part of the daily command. If you're using each 
command individually it should be called after calling read to generate a page
of posts to read:

feeder fetch # retrieve the latest from the feed
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
	return markCmd
}

func init() {
	RegisterSubcommand(NewMarkCmd)
}
