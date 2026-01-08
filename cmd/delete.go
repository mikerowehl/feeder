/*
Copyright (c) Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package cmd

import (
	"strconv"

	"github.com/mikerowehl/feeder/internal/feeder"
	"github.com/spf13/cobra"
)

func NewDeleteCmd() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete ID",
		Short: "Remove a feed from the database using its ID",
		Long: `Removes a feed from the database using the internal ID of the feed given by
the list command.

ex: feeder delete 5`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			u64, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return err
			}
			feedId := uint(u64)
			f := cmd.Context().Value(feederKey).(*feeder.Feeder)
			return f.Delete(feedId)
		},
	}
	return deleteCmd
}

func init() {
	RegisterSubcommand(NewDeleteCmd)
}
