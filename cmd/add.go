/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package cmd

import (
	"github.com/mikerowehl/feeder/internal/feeder"
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
		f := cmd.Context().Value(feederKey).(*feeder.Feeder)
		return f.Add(feedUrl)
	},
}

func init() {
	RegisterSubcommand(addCmd)
}
