package main

import (
	"bufio"
	"context"
	_ "errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	Address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer

	inChan  chan string
	outChan chan string
	conn    net.Conn

	outScanner bufio.Scanner
}

func (c *telnetClient) Send() error {
	go func() {
		for str := range c.inChan {
			c.conn.Write([]byte(fmt.Sprintf("%s\n", str)))
		}
	}()
	return nil
}

func (c *telnetClient) inScan() {
	go func() {
		scanner := bufio.NewScanner(c.in)
		for scanner.Scan() {
			c.inChan <- scanner.Text()
		}
		if scanner.Err() != nil {
			close(c.inChan)
		}
	}()
}

func (c *telnetClient) Connect() error {
	dialer := &net.Dialer{}
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	var err error
	c.conn, err = dialer.DialContext(ctx, "tcp", c.Address)
	if err != nil {
		cancel()
		fmt.Println("err", err)
		return err
	}
	log.Printf("connect from %s to %s\n", c.conn.LocalAddr(), c.conn.RemoteAddr())
	c.outScanner = *bufio.NewScanner(c.conn)
	c.inScan()

	return nil
}

func (c *telnetClient) Close() error {
	fmt.Println("close connection")
	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("close connection: %w", err)
	}
	err = c.in.Close()
	if err != nil {
		return fmt.Errorf("close input: %w", err)
	}
	close(c.inChan)
	return nil
}

func (c *telnetClient) Receive() error {
	if c.outScanner.Scan() {
		c.out.Write([]byte(c.outScanner.Text() + "\n"))
	}
	if c.outScanner.Err() != nil {
		return c.outScanner.Err()
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	inChan := make(chan string)
	outChan := make(chan string)
	return &telnetClient{
		Address: address,
		timeout: timeout,
		in:      in,
		out:     out,

		inChan:  inChan,
		outChan: outChan,
	}
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
