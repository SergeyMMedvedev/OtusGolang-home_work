package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/pflag"
)

var (
	in  = os.Stdin
	out = os.Stdout
)

func parseArgs() (address string, timeout time.Duration, err error) {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	var host string
	var port int
	timeout, err = time.ParseDuration("10s")
	if err != nil {
		return "", 0, err
	}
	timeoutVal := pflag.String("timeout", "", "connection timeout")
	pflag.Lookup("timeout").NoOptDefVal = "10s"
	pflag.Parse()
	if *timeoutVal == "" {
		if len(os.Args) != 3 {
			return "", 0, fmt.Errorf("host and port is required")
		}
		host = os.Args[1]
		port, err = strconv.Atoi(os.Args[2])
	} else {
		if len(os.Args) != 4 {
			return "", 0, fmt.Errorf("host and port is required")
		}
		timeout, err = time.ParseDuration(*timeoutVal)
		if err != nil {
			return "", 0, fmt.Errorf("timeout parse error: %w", err)
		}
		host = os.Args[2]
		port, err = strconv.Atoi(os.Args[3])
	}
	if err != nil {
		return "", 0, fmt.Errorf("port must be a number")
	}
	address = fmt.Sprintf("%s:%d", host, port)
	return address, timeout, nil
}

func main() {
	address, timeout, err := parseArgs()
	if err != nil {
		fmt.Println(err)
		return
	}
	telnet := NewTelnetClient(address, timeout, in, out)
	err = telnet.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer telnet.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
	OUTER:
		for {
			select {
			case <-ctx.Done():
				break OUTER
			default:
				err := telnet.Send()
				if err != nil {
					fmt.Println(err)
					cancel()
					break OUTER
				}
			}
		}
		defer wg.Done()
	}()
	wg.Add(1)
	go func() {
	OUTER:
		for {
			select {
			case <-ctx.Done():
				break OUTER
			default:
				err := telnet.Receive()
				if err != nil {
					fmt.Println(err)
					cancel()
					break OUTER
				}
			}
		}
		defer wg.Done()
	}()
	wg.Wait()
}
