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
	"log"
	"os"
	"path/filepath"

	"github.com/mikerowehl/feeder/internal/feeder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type feederKeyType struct{}

var feederKey = feederKeyType{}

// In a simple Cobra app there's a root command and the init() for each
// subcommand adds itself to the root. But we want to be able to create new
// rootCmd instances on the fly for testing. So instead we have the
// subcommands register generator functions that we can call to make a new
// instance for each new rootCmd.
type SubcommandFactory func() *cobra.Command

var subcommands []SubcommandFactory

func RegisterSubcommand(factory SubcommandFactory) {
	subcommands = append(subcommands, factory)
}

func checkedBinding(name string, cmd *cobra.Command) {
	err := viper.BindPFlag(name, cmd.PersistentFlags().Lookup(name))
	if err != nil {
		log.Printf("Error binding %s flag\n", name)
	}
}

func NewRootCommand(skipConfig bool) *cobra.Command {
	var cfgFile string

	var rootCmd = &cobra.Command{
		Use:   "feeder",
		Short: "A basic command line syndicated feed processor",
		Long: `Use the add command to build up a list of feeds to process. The fetch
command then pulls down the feeds and merges them into a summary page.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			dbDir := feeder.ExpandPath(viper.GetString("db-dir"))
			dbFile := viper.GetString("db-file")
			filename := filepath.Join(dbDir, dbFile)
			f, err := feeder.NewFeeder(filename, cmd.OutOrStdout(), cmd.ErrOrStderr(), cmd.InOrStdin())
			if err != nil {
				return err
			}
			f.Verbose = viper.GetBool("verbose")

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
	rootCmd.PersistentFlags().String("output", "", "filename to output HTML")
	rootCmd.PersistentFlags().Bool("verbose", false, "Output additional info during run")

	checkedBinding("db-dir", rootCmd)
	checkedBinding("db-file", rootCmd)
	checkedBinding("max-items", rootCmd)
	checkedBinding("output", rootCmd)
	checkedBinding("verbose", rootCmd)

	for _, factory := range subcommands {
		rootCmd.AddCommand(factory())
	}

	return rootCmd
}

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
