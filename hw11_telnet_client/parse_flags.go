package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spf13/pflag"
)

type Flags struct {
	address string
	timeout time.Duration
}

func parseArgs() (flags Flags, err error) {
	var timeout time.Duration
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	defaultTimeout, err := time.ParseDuration("10s")
	if err != nil {
		return flags, err
	}
	pflag.DurationVar(&timeout, "timeout", defaultTimeout, "connection timeout")
	pflag.Lookup("timeout").NoOptDefVal = "10s"
	pflag.Parse()
	args := pflag.Args()
	if len(args) != 2 {
		return flags, fmt.Errorf("host and port is required")
	}
	return Flags{
		address: net.JoinHostPort(args[0], args[1]),
		timeout: timeout,
	}, nil
}
