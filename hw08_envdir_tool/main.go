package main

import (
	"fmt"
	"os"
)

func main() {
	// Place your code here.
	cmd := os.Args[1:]
	if len(cmd) < 2 {
		fmt.Println("No command given")
		os.Exit(1)
	}
	envDir := cmd[0]
	environment, err := ReadDir(envDir)
	if err != nil {
		fmt.Println("err:", err)
		os.Exit(1)
	}
	s := RunCmd(cmd[1:], environment)
	os.Exit(s)
}
