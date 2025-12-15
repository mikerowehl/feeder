/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
// This has a root command generator function instead of the normal simple
// Cobra top level root static variable so that we can run tests more easily.
// We have each subcommand register with a custom hook so that as we create a
// root command for each test we can add the subcommands to it. There's also a
// flag to skip config file processing to try to keep the test state as clean
// as possible while keeping the setup clean.
package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mikerowehl/feeder/internal/feeder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const feederKey = "feeder"

var subcommands []*cobra.Command

func RegisterSubcommand(cmd *cobra.Command) {
	subcommands = append(subcommands, cmd)
}

func NewRootCommand(skipConfig bool) *cobra.Command {
	var cfgFile string

	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "feeder",
		Short: "A basic command line syndicated feed processor",
		Long: `Use the add command to build up a list of feeds to process. The fetch
command then pulls down the feeds and merges them into a summary page.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			dbDir := feeder.ExpandPath(viper.GetString("db-dir"))
			dbFile := viper.GetString("db-file")
			filename := filepath.Join(dbDir, dbFile)
			f, err := feeder.NewFeeder(filename)
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

	if !skipConfig {
		cobra.OnInitialize(func() {
			initConfig(cfgFile)
		})
	}

	defaultConfigDir := feeder.GetConfigDir()
	configHelp := fmt.Sprintf("config file (default %s/config.yaml)", defaultConfigDir)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", configHelp)

	defaultDataDir := feeder.GetDataDir()
	dataDirHelp := fmt.Sprintf("database file directory (default %s)", defaultDataDir)
	rootCmd.PersistentFlags().String("db-dir", defaultDataDir, dataDirHelp)
	rootCmd.PersistentFlags().String("db-file", "feeder.db",
		"database file name (default feeder.db)")
	rootCmd.PersistentFlags().Int("max-items", 100,
		"Maximum number of items to store per feed")

	viper.BindPFlag("db-dir", rootCmd.PersistentFlags().Lookup("db-dir"))
	viper.BindPFlag("db-file", rootCmd.PersistentFlags().Lookup("db-file"))
	viper.BindPFlag("max-items", rootCmd.PersistentFlags().Lookup("max-items"))

	for _, cmd := range subcommands {
		rootCmd.AddCommand(cmd)
	}

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd := NewRootCommand(false)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initConfig(cfgFile string) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(feeder.GetConfigDir())
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}
	viper.SetEnvPrefix("FEEDER")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// File not found, this isn't an error - no output
		} else {
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
		}
	}
}
