/*
Copyright (c) 2025 - Mike Rowehl <mikerowehl@gmail.com>
This software may be modified and distributed under the terms of the MIT license.
See LICENSE in the project root for full license information.
*/
package integration

import (
	"bytes"
	"testing"

	"github.com/mikerowehl/feeder/cmd"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func executeCommand(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	viper.Reset()
	rootCmd := cmd.NewRootCommand(true)

	stdoutBuf := new(bytes.Buffer)
	rootCmd.SetOut(stdoutBuf)
	stderrBuf := new(bytes.Buffer)
	rootCmd.SetErr(stderrBuf)
	rootCmd.SetArgs(args)

	err := rootCmd.Execute()
	return stdoutBuf.String(), stderrBuf.String(), err
}

func TestIntegration_Help(t *testing.T) {
	commands := []string{"add"}

	for _, cmd := range commands {
		t.Run(cmd, func(t *testing.T) {
			stdout, _, err := executeCommand(t, cmd, "--help")
			require.NoError(t, err)
			assert.NotEmpty(t, stdout)
			assert.Contains(t, stdout, "Usage:", "help should show usage")
		})
	}
}

func TestIntegration_Add(t *testing.T) {
	tmpDir := t.TempDir()
	server := startTestFeedServer(t)
	feedURL := getTestFeedURL(server, "basic.xml")
	testArgs := []string{"--db-dir", tmpDir, "--db-file", "test.db"}

	stdout, _, err := executeCommand(t, append(testArgs, "add", feedURL)...)
	require.NoError(t, err)

	stdout, _, err = executeCommand(t, append(testArgs, "list")...)
	require.NoError(t, err)
	assert.Contains(t, stdout, "Feeder Basic Integration Test")
}
