package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spf13/pflag"
)

func isFlagPassed(name string) bool {
	found := false
	pflag.Visit(func(f *pflag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func parseArgs() (address string, timeout time.Duration, err error) {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	var host string
	var port string
	defaultTimeout, err := time.ParseDuration("10s")
	if err != nil {
		return "", 0, err
	}
	fmt.Println("osArgs", os.Args)
	pflag.DurationVar(&timeout, "timeout", defaultTimeout, "connection timeout")
	pflag.Lookup("timeout").NoOptDefVal = "10s"
	pflag.Parse()

	flagPassed := isFlagPassed("timeout")

	fmt.Println("flagPassed", flagPassed)
	if !flagPassed {
		if len(os.Args) != 3 {
			return "", 0, fmt.Errorf("host and port is required")
		}
		host = os.Args[1]
		port = os.Args[2]
	} else {
		if len(os.Args) != 4 {
			return "", 0, fmt.Errorf("host and port is required")
		}
		host = os.Args[2]
		port = os.Args[3]
	}
	if err != nil {
		return "", 0, fmt.Errorf("port must be a number")
	}
	return net.JoinHostPort(host, port), timeout, nil
}
