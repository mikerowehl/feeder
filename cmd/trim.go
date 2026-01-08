/*
Copyright (c) Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package cmd

import (
	"github.com/mikerowehl/feeder/internal/feeder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewTrimCmd() *cobra.Command {
	trimCmd := &cobra.Command{
		Use:   "trim",
		Short: "Trims the number of items per feed to a max number",
		Long: `Keeps only the most recent items for each feed. This keeps the local database
from getting too large and keeps things running quickly. Also runs some
housekeeping on the database file to optimize performance.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			f := cmd.Context().Value(feederKey).(*feeder.Feeder)
			maxItems := viper.GetInt("max-items")
			return f.Trim(maxItems)
		},
	}
	return trimCmd
}

func init() {
	RegisterSubcommand(NewTrimCmd)
}
