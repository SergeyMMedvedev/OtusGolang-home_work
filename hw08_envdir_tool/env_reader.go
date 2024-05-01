package main

import (
	"bytes"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := Environment{}
	for _, file := range files {
		var needRemove bool
		if file.IsDir() {
			continue
		}
		content, err := os.ReadFile(dir + "/" + file.Name())
		if err != nil {
			return nil, err
		}
		if len(content) == 0 {
			os.Unsetenv(file.Name())
			continue
		}
		contains := strings.Contains(file.Name(), "=")
		if contains {
			continue
		}
		_, ok := os.LookupEnv(file.Name())
		if ok {
			needRemove = true
		}
		lines := strings.Split(string(content), "\n")
		newContent := string(bytes.ReplaceAll([]byte(lines[0]), []byte{0x00}, []byte{'\n'}))
		newLine := strings.TrimRight(newContent, "\t ")
		env[file.Name()] = EnvValue{
			Value:      newLine,
			NeedRemove: needRemove,
		}
	}
	return env, nil
}
