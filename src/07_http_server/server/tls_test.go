// go get github.com/awoodbeck/gnp/ch09/handlers
// go test -v tls_test.go

package server

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/awoodbeck/gnp/ch09/handlers"
)

func TestSimpleHTTPSServer(t *testing.T) {
	srv := &http.Server{
		Addr:        "127.0.0.1:8443",
		Handler:     http.TimeoutHandler(handlers.DefaultHandler(), 2*time.Minute, ""),
		IdleTimeout: 5 * time.Minute,
		ReadTimeout: 1 * time.Minute,
	}

	l, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		err := srv.ServeTLS(l, "./certs/cert.pem", "./certs/key.pem")
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

	certPEM, err := os.ReadFile("./certs/cert.pem")
	if err != nil {
		t.Fatal(err)
	}

	rootCAs, _ := x509.SystemCertPool()
	if ok := rootCAs.AppendCertsFromPEM(certPEM); !ok {
		t.Fatal(err)
	}
	config := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            rootCAs,
	}
	tr := &http.Transport{TLSClientConfig: config}
	client := new(http.Client)
	client.Transport = tr
	path := fmt.Sprintf("https://%s/", srv.Addr)

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

=== RUN   TestSimpleHTTPSServer
    tls_test.go:97: 0: 200 OK - Hello, friend!
    tls_test.go:97: 1: 200 OK - Hello, &lt;world&gt;!
    tls_test.go:97: 2: 405 Method Not Allowed -
--- PASS: TestSimpleHTTPSServer (0.06s)
PASS

*/
