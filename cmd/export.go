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

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Writes the list of feeds from the database out to standard output",
	Long: `Very minimal output just written to standard output. This can be captured and then
fed back into the import command when rebuilding the database. There isn't
currently any way to capture the read/unread status. So I normally get the
feeds caught up, export and import, and just mark everything read.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		f := cmd.Context().Value(feederKey).(*feeder.Feeder)
		err := f.Export()
		if err != nil {
			return fmt.Errorf("error exporting feeds: %w", err)
		}
		return nil
	},
}

func init() {
	RegisterSubcommand(exportCmd)
}
