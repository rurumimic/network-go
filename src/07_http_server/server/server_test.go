// go get github.com/awoodbeck/gnp/ch09/handlers
// go test -v server_test.go

package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/awoodbeck/gnp/ch09/handlers"
)

func TestSimpleHTTPServer(t *testing.T) {
	srv := &http.Server{
		Addr:        ":8081",
		Handler:     http.TimeoutHandler(handlers.DefaultHandler(), 2*time.Minute, ""),
		IdleTimeout: 5 * time.Minute,
		ReadTimeout: 1 * time.Minute,
	}

	l, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		err := srv.Serve(l)
		if err != http.ErrServerClosed {
			t.Error(err)
		}
	}()

	testCases := []struct {
		method   string
		body     io.Reader
		code     int
		response string
	}{
		{http.MethodGet, nil, http.StatusOK, "Hello, friend!"},
		{http.MethodPost, bytes.NewBufferString("<world>"), http.StatusOK, "Hello, &lt;world&gt;!"},
		{http.MethodHead, nil, http.StatusMethodNotAllowed, ""},
	}

	client := new(http.Client)
	path := fmt.Sprintf("http://%s/", srv.Addr)

	for i, c := range testCases {
		r, err := http.NewRequest(c.method, path, c.body)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		resp, err := client.Do(r)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		if resp.StatusCode != c.code {
			t.Errorf("%d: unexpected status code %q", i, resp.Status)
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}
		_ = resp.Body.Close()

		if c.response != string(b) {
			t.Errorf("%d: expected %q, got %q", i, c.response, b)
		} else {
			t.Logf("%d: %s - %s", i, resp.Status, b)
		}

	}

	if err := srv.Close(); err != nil {
		t.Fatal(err)
	}
}

/*

=== RUN   TestSimpleHTTPServer
    server_test.go:79: 0: 200 OK - Hello, friend!
    server_test.go:79: 1: 200 OK - Hello, &lt;world&gt;!
    server_test.go:79: 2: 405 Method Not Allowed -
--- PASS: TestSimpleHTTPServer (0.00s)
PASS

*/
