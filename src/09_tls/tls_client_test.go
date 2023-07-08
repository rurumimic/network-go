// go test -race -v tls_client_test.go

package mytls

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/http2"
)

func TestClientTLS(t *testing.T) {
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.TLS == nil {
				u := "https://" + r.Host + r.RequestURI
				http.Redirect(w, r, u, http.StatusMovedPermanently)
				return
			}
			w.WriteHeader(http.StatusOK)
		}))
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d; actual status %d", http.StatusOK, resp.StatusCode)
	}

	tp := &http.Transport{
		TLSClientConfig: &tls.Config{
			CurvePreferences: []tls.CurveID{tls.CurveP256},
			MinVersion:       tls.VersionTLS12,
		},
	}

	err = http2.ConfigureTransport(tp)
	if err != nil {
		t.Fatal(err)
	}

	client2 := &http.Client{Transport: tp}

	_, err = client2.Get(ts.URL)
	if err == nil || !strings.Contains(err.Error(), "certificate signed by unknown authority") {
		t.Fatalf("expected unknown authority error; actual: %q", err)
	}

	tp.TLSClientConfig.InsecureSkipVerify = true
	resp, err = client2.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d; actual status %d", http.StatusOK, resp.StatusCode)
	}

}

func TestClientTLSGoogle(t *testing.T) {
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: 30 * time.Second},
		"tcp",
		"www.google.com:443",
		&tls.Config{
			CurvePreferences: []tls.CurveID{tls.CurveP256},
			MinVersion:       tls.VersionTLS12,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	state := conn.ConnectionState()
	t.Logf("TLS 1.%d", state.Version-tls.VersionTLS10)
	t.Log(tls.CipherSuiteName(state.CipherSuite))
	t.Log(state.VerifiedChains[0][0].Issuer.Organization[0])

	_ = conn.Close()
}

/*

go test -race -v tls_client_test.go

=== RUN   TestClientTLS
2023/07/08 15:32:21 http: TLS handshake error from 127.0.0.1:54291: remote error: tls: bad certificate
--- PASS: TestClientTLS (0.05s)
=== RUN   TestClientTLSGoogle
    tls_client_test.go:84: TLS 1.3
    tls_client_test.go:85: TLS_AES_128_GCM_SHA256
    tls_client_test.go:86: Google Trust Services LLC
--- PASS: TestClientTLSGoogle (0.11s)
PASS
ok  	command-line-arguments	0.387s

*/

/*

go test -v -race -run TestClientTLSGoogle tls_client_test.go

=== RUN   TestClientTLSGoogle
    tls_client_test.go:83: TLS 1.3
    tls_client_test.go:84: TLS_AES_128_GCM_SHA256
    tls_client_test.go:85: Google Trust Services LLC
--- PASS: TestClientTLSGoogle (0.25s)
PASS
ok  	command-line-arguments	0.825s

*/
