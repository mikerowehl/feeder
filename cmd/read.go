/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package cmd

import (
	"fmt"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Context().Value(feederKey).(*feeder.Feeder)
		outfile := fmt.Sprintf("feeder-%s.html", time.Now().Format(time.DateOnly))
		err := f.WriteUnread(outfile)
		if err != nil {
			return fmt.Errorf("error writing out unread: %w", err)
		}
		fmt.Println("read called")
		return nil
	},
}

func init() {
	RegisterSubcommand(readCmd)
}
