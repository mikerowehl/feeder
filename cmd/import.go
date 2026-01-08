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

func NewImportCmd() *cobra.Command {
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "Reads a set of urls from standard input and adds them",
		Long: `Input on standard input should be a set of urls, one per line. This will read
the urls and add each one to the database.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			f := cmd.Context().Value(feederKey).(*feeder.Feeder)
			err := f.Import()
			if err != nil {
				return fmt.Errorf("error importing feeds: %w", err)
			}
			return nil
		},
	}
	return importCmd
}

func init() {
	RegisterSubcommand(NewImportCmd)
}
