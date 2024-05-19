package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spf13/pflag"
)

func parseArgs() (address string, timeout time.Duration, err error) {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	defaultTimeout, err := time.ParseDuration("10s")
	if err != nil {
		return "", 0, err
	}
	pflag.DurationVar(&timeout, "timeout", defaultTimeout, "connection timeout")
	pflag.Lookup("timeout").NoOptDefVal = "10s"
	pflag.Parse()
	args := pflag.Args()
	if len(args) != 2 {
		return "", 0, fmt.Errorf("host and port is required")
	}
	host := args[0]
	port := args[1]
	return net.JoinHostPort(host, port), timeout, nil
}
