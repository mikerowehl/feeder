/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package cmd

import (
	"log"

	"github.com/mikerowehl/feeder/internal/feeder"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the active feeds",
	Long:  `Outputs the title and URL of each feed from the database onto standard output.`,
	Run: func(cmd *cobra.Command, args []string) {
		f, err := feeder.NewFeeder(dbFile)
		if err != nil {
			log.Fatalf("Startup error: %v", err)
		}
		defer f.Close()
		err = f.List()
		if err != nil {
			log.Fatalf("Error fetching feeds: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
