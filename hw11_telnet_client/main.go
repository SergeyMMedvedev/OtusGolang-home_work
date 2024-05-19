package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	_ "sync"
	"syscall"
)

var (
	in  = os.Stdin
	out = os.Stdout
)

func main() {
	flags, err := parseArgs()
	if err != nil {
		fmt.Println(err)
		return
	}
	telnet := NewTelnetClient(flags.address, flags.timeout, in, out)
	err = telnet.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer telnet.Close()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("exiting Send")
			return
		default:
			err := telnet.Send()
			if err != nil {
				fmt.Println(err)
				stop()
			}
		}
	}()
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("exiting Receive")
			return
		default:
			err := telnet.Receive()
			if err != nil {
				fmt.Println(err)
				stop()
			}
		}
	}()
	<-ctx.Done()
}
