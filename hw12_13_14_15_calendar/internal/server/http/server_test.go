package internalhttp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMyHandlerByNewServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer ts.Close()

	cases := []struct {
		name         string
		method       string
		target       string
		body         io.Reader
		responseCode int
	}{
		{"ok", http.MethodGet, "/hello", nil, http.StatusOK},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, c.method, ts.URL+c.target, c.body)
			require.NoError(t, err)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer res.Body.Close()
			require.Equal(t, c.responseCode, res.StatusCode)
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			fmt.Println(string(body))
			require.True(t, strings.HasPrefix(string(body), "Hello, "))
		})
	}
}
