/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mikerowehl/feeder/internal/feeder"
	"github.com/spf13/cobra"
)

const feederKey = "feeder"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "feeder",
	Short: "A basic command line syndicated feed processor",
	Long: `Use the add command to build up a list of feeds to process. The fetch
command then pulls down the feeds and merges them into a summary page.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		dbFile, err := cmd.Flags().GetString("dbFile")
		if err != nil {
			return err
		}

		f, err := feeder.NewFeeder(dbFile)
		if err != nil {
			return err
		}

		ctx := context.WithValue(cmd.Context(), feederKey, f)
		cmd.SetContext(ctx)

		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		fVal := cmd.Context().Value(feederKey)
		if fVal == nil {
			return nil
		}

		f, ok := fVal.(*feeder.Feeder)
		if !ok {
			return fmt.Errorf("%s key has wrong type", feederKey)
		}
		f.Close()

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("dbFile", "feeder.db",
		"database file name (default feeder.db)")
}
