package main

import (
	"bufio"
	"context"
	_ "errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func (c telnetClient) readRoutine(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	log.Printf("start readRoutine")
	defer wg.Done()
	fmt.Println("conn", conn)
	scanner := bufio.NewScanner(conn)
	fmt.Println("scanner", scanner)
OUTER:
	for {
		fmt.Println("0")
		select {
		case <-ctx.Done():
			fmt.Println("readRoutine ctx.Done")
			break OUTER
		default:
			fmt.Println("1")
			scan := scanner.Scan()
			fmt.Println("scan", scan)
			if !scan {
				log.Printf("CANNOT SCAN")
				break OUTER
			}
			text := scanner.Text() + "\n"
			log.Printf("From server: %s", text)
			c.out.Write([]byte(text))
		}
	}
	log.Printf("Finished readRoutine")
}

func (c *telnetClient) Send() error {
	go func() {
		str := <-c.inChan
		log.Printf("To server %v\n", str)
		fmt.Println("str", str)
		fmt.Println("c.conn", c.conn)
		c.conn.Write([]byte(fmt.Sprintf("%s\n", str)))
	}()
	return nil
}

type telnetClient struct {
	Address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer

	inChan chan string
	conn   net.Conn
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
	fmt.Println("c.timeout", c.timeout)
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	fmt.Println("c.Address", c.Address)
	var err error
	c.conn, err = dialer.DialContext(ctx, "tcp", c.Address)
	if err != nil {
		cancel()
		fmt.Println("err", err)
		return err
	}
	log.Printf("connect from %s to %s\n", c.conn.LocalAddr(), c.conn.RemoteAddr())
	c.inScan()

	return nil
}

func (c *telnetClient) Close() error {
	return c.in.Close()
}

func writeRoutine(ctx context.Context, conn net.Conn, wg *sync.WaitGroup, stdin chan string) {
	log.Printf("start writeRoutine")
	conn.Write([]byte(fmt.Sprintf("hello\n")))
	defer wg.Done()
OUTER:

	for {
		fmt.Println("len(stdin)", len(stdin))
		select {
		case <-ctx.Done():
			fmt.Println("writeRoutine ctx.Done")
			break OUTER
		case str := <-stdin:
			fmt.Println("str", str)
			fmt.Println("len(stdin)", len(stdin))
			log.Printf("To server %v\n", str)
			conn.Write([]byte(fmt.Sprintf("%s\n", str)))
		}
	}
	log.Printf("Finished writeRoutine")
}

func (c *telnetClient) Receive() error {
	scanner := bufio.NewScanner(c.conn)
	if scanner.Scan() {
		c.out.Write([]byte(scanner.Text() + "\n"))
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	inChan := make(chan string)
	return &telnetClient{
		Address: address,
		timeout: timeout,
		in:      in,
		out:     out,

		inChan: inChan,
	}
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
