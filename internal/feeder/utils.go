/*
Copyright (c) Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package feeder

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const dataDirDefault = "."

// GetDataDir returns the default directory to use for data files manged by
// the program. If XDG_DATA_HOME is set we use a subdirectory under that
// location. If not we try to use the platform to pick an appropriate default.
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

// GetConfigDir returns the default directory to use for configuration files.
// If the XDG_CONFIG_HOME environment variable is set we use a subdirectory
// under that. If not we default to a different platform specific location
// based on common practice.
func GetConfigDir() string {
	var configDir string
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		configDir = filepath.Join(xdgConfigHome, appName)
	} else {
		switch runtime.GOOS {
		case "linux":
			homeDir, err := os.UserHomeDir()
			if err != nil {
				homeDir = "."
			}
			configDir = filepath.Join(homeDir, ".config", appName)
		case "darwin":
			homeDir, err := os.UserHomeDir()
			if err != nil {
				homeDir = "."
			}
			configDir = filepath.Join(homeDir, "Library", "Application Support", appName)
		case "windows":
			if appData := os.Getenv("APPDATA"); appData != "" {
				configDir = filepath.Join(appData, appName)
			} else {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					homeDir = "."
				}
				configDir = filepath.Join(homeDir, "AppData", "Roaming", appName)
			}
		default:
			configDir = "."
		}
	}

	return configDir
}

// ExpandPath provides the home directory shortcut (~) expansion and expands
// Unix and Windows style environment variable references
func ExpandPath(path string) string {
	if path == "" {
		return path
	}

	if strings.HasPrefix(path, "~/") {
		if homeDir, err := os.UserHomeDir(); err == nil {
			path = filepath.Join(homeDir, path[2:])
		}
	} else if path == "~" {
		if homeDir, err := os.UserHomeDir(); err == nil {
			path = homeDir
		}
	}
	path = os.ExpandEnv(path)
	return path
}

// Wrapper for Fprintf that writes a log to standard output if there's a
// problem outputting to the desired handle.
func LoggedPrint(w io.Writer, format string, args ...any) {
	_, err := fmt.Fprintf(w, format, args...)
	if err != nil {
		log.Printf("Error writing output: %v", err)
	}
}
