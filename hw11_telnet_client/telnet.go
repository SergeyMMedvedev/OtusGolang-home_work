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

	conn       net.Conn
	inScanner  bufio.Scanner
	outScanner bufio.Scanner
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
	c.inScanner = *bufio.NewScanner(c.in)
	c.outScanner = *bufio.NewScanner(c.conn)
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
	return nil
}

func (c *telnetClient) writeFromScanner(scanner *bufio.Scanner, writer io.Writer) error {
	if scanner.Scan() {
		_, err := writer.Write([]byte(fmt.Sprintf("%s\n", scanner.Text())))
		if err != nil {
			return fmt.Errorf("write: %w", err)
		}
	}
	return scanner.Err()
}

func (c *telnetClient) Send() error {
	return c.writeFromScanner(&c.inScanner, c.conn)
}

func (c *telnetClient) Receive() error {
	return c.writeFromScanner(&c.outScanner, c.out)
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
