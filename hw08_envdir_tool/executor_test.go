package main

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	shell := "/bin/bash"
	testDir := "./testdata/env2"
	env, _ := ReadDir(testDir)
	if runtime.GOOS == "windows" {
		shell = "cmd.exe"
	}
	cmd := []string{shell, "./testdata/echo.sh", "arg1=1", "arg2=2"}
	returnCode := RunCmd(cmd, env)
	require.Equal(t, 0, returnCode)
}
