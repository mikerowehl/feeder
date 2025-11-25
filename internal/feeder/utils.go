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

func GetDataDir() (string, error) {
	var dataDir string
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if xdgDataHome != "" {
		dataDir = filepath.Join(xdgDataHome, appName)
	} else {
		switch runtime.GOOS {
		case "linux":
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			dataDir = filepath.Join(homeDir, ".local", "share", appName)

		case "darwin":
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			dataDir = filepath.Join(homeDir, "Library", "Application Support", appName)

		case "windows":
			if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
				dataDir = filepath.Join(localAppData, appName)
			} else {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return "", err
				}
				dataDir = filepath.Join(homeDir, "AppData", "Local", appName)
			}
		}
	}

	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return "", err
	}

	return dataDir, nil
}
