package main

import (
	_ "bufio"
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	address = "localhost:4242"
	timeout = time.Duration(500) * time.Second
	in      = os.Stdin
	out     = os.Stdout
)

func main() {
	telnet := NewTelnetClient(address, timeout, in, out)
	err := telnet.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer telnet.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := telnet.Send()
		if err != nil {
			fmt.Println(err)
		}
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
