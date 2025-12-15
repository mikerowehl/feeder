/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewConfigCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Display all the config values currently active",
		Long: `Outputs the location of the config file being used and the values of all the
	settings configured using any config mechanism (flags, config file, or environment).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfgFile := viper.ConfigFileUsed(); cfgFile != "" {
				fmt.Printf("Config file: %s\n", cfgFile)
			} else {
				fmt.Printf("No config file being used\n")
			}
			settings := viper.AllSettings()
			for key, value := range settings {
				fmt.Printf("%s: %v\n", key, value)
			}
			return nil
		},
	}
	return configCmd
}

func init() {
	RegisterSubcommand(NewConfigCmd)
}
