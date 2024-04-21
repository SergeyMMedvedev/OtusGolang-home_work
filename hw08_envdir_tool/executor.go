package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	command.Env = os.Environ()
	for k, envValue := range env {
		if envValue.NeedRemove {
			os.Unsetenv(k)
		}
		command.Env = append(command.Env, fmt.Sprintf("%s=%s", k, envValue.Value))
	}
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Start()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	err = command.Wait()
	if err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}
