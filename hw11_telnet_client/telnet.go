package main

import (
	"context"
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

	conn net.Conn
}

func (c *telnetClient) Connect() error {
	dialer := &net.Dialer{}
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	var err error
	c.conn, err = dialer.DialContext(ctx, "tcp", c.Address)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	log.Printf("connect from %s to %s\n", c.conn.LocalAddr(), c.conn.RemoteAddr())
	return nil
}

func (c *telnetClient) Close() error {
	fmt.Println("close connection")
	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("close connection: %w", err)
	}
	err = c.in.Close()
	if err != nil {
		return fmt.Errorf("close input: %w", err)
	}
	return nil
}

func (c *telnetClient) Send() error {
	_, err := io.Copy(c.conn, c.in)
	return err
}

func (c *telnetClient) Receive() error {
	_, err := io.Copy(c.out, c.conn)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		Address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
