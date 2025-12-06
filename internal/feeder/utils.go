/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package feeder

import (
	"os"
	"path/filepath"
	"runtime"
)

const dataDirDefault = "."

func GetDataDir() string {
	var dataDir string
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome != "" {
		dataDir = filepath.Join(xdgDataHome, appName)
	} else {
		switch runtime.GOOS {
		case "linux":
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return dataDirDefault
			}
			dataDir = filepath.Join(homeDir, ".local", "share", appName)

		case "darwin":
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return dataDirDefault
			}
			dataDir = filepath.Join(homeDir, "Library", "Application Support", appName)

		case "windows":
			if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
				dataDir = filepath.Join(localAppData, appName)
			} else {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return dataDirDefault
				}
				dataDir = filepath.Join(homeDir, "AppData", "Local", appName)
			}
		}
	}

	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return dataDirDefault
	}

	return dataDir
}

func GetConfigDir() string {
	var configDir string
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		configDir = filepath.Join(xdgConfigHome, appName)
	} else {
		switch runtime.GOOS {
		case "linux":
			homeDir, _ := os.UserHomeDir()
			configDir = filepath.Join(homeDir, ".config", appName)
		case "darwin":
			homeDir, _ := os.UserHomeDir()
			configDir = filepath.Join(homeDir, "Library", "Application Support", appName)
		case "windows":
			if appData := os.Getenv("APPDATA"); appData != "" {
				configDir = filepath.Join(appData, appName)
			} else {
				homeDir, _ := os.UserHomeDir()
				configDir = filepath.Join(homeDir, "AppData", "Roaming", appName)
			}
		default:
			configDir = "."
		}
	}

	return configDir
}
