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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the active feeds",
	Long:  `Outputs the title and URL of each feed from the database onto standard output.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Context().Value(feederKey).(*feeder.Feeder)
		err := f.List()
		if err != nil {
			return fmt.Errorf("error fetching feeds: %w", err)
		}
		return nil
	},
}

func init() {
	RegisterSubcommand(listCmd)
}
