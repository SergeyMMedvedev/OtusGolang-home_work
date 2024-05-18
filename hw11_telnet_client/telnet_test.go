package main

import (
	"bytes"
	"io"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}

func TestArgParse(t *testing.T) {
	t.Run("check parseArgs", func(t *testing.T) {
		originalArgs := os.Args
		defer func() { os.Args = originalArgs }()
		os.Args = []string{"go-telnet", "127.0.0.1", "80"}
		address, timeout, err := parseArgs()
		require.NoError(t, err)
		require.Equal(t, "127.0.0.1:80", address)
		require.Equal(t, time.Second*10, timeout)
	})

	t.Run("check parseArgs no host&port", func(t *testing.T) {
		originalArgs := os.Args
		defer func() { os.Args = originalArgs }()
		os.Args = []string{"go-telnet"}
		_, _, err := parseArgs()
		require.Error(t, err)
		require.Equal(t, "host and port is required", err.Error())
	})

	t.Run("check parseArgs timeout", func(t *testing.T) {
		originalArgs := os.Args
		defer func() { os.Args = originalArgs }()
		os.Args = []string{"go-telnet", "--timeout=12s", "127.0.0.1", "80"}
		_, timeout, err := parseArgs()
		require.NoError(t, err)
		require.Equal(t, time.Second*12, timeout)
	})
}
