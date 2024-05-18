package main

import (
	"context"
	"fmt"
	"os"
	"sync"
)

var (
	in  = os.Stdin
	out = os.Stdout
)

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
