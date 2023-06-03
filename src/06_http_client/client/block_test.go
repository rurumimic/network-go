// go test -v -timeout 10s block_test.go

package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func blockIndefinitely(w http.ResponseWriter, r *http.Request) {
	select {}
}

/*

func TestBlockIndefinitely(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(blockIndefinitely))
	_, _ = http.Get(ts.URL)
	t.Fatal("client did not indefinitely block")
}

=== RUN   TestBlockIndefinitely
panic: test timed out after 5s

*/

func TestBlockIndefinitelyWithTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(blockIndefinitely))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil) // nil = GET no payload
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Logf("Error: %v", err)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatal(err)
		}
		return
	}
	_ = resp.Body.Close() // no exception handling in tests
}

/*

=== RUN   TestBlockIndefinitelyWithTimeout
    block_test.go:43: Error: Get "http://127.0.0.1:56111": context deadline exceeded
--- PASS: TestBlockIndefinitelyWithTimeout (5.00s)

*/

func TestBlockIndefinitelyWithCancel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(blockIndefinitely))

	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(3*time.Second, cancel)
	defer func() {
		t.Log("Canceling Request")
		cancel()
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ts.URL, nil) // nil = GET no payload
	if err != nil {
		t.Fatal(err)
	}
	req.Close = true

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Log("Client timeout")
		if !errors.Is(err, context.Canceled) {
			t.Log("Not Cancelled")
			t.Fatal(err)
		}
		return
	}

	// Read Header

	h := resp.Header
	t.Logf("Header: %v", h)

	// Add 5 seconds to the timer
	timer.Reset(3 * time.Second)

	// Read Body

	_ = resp.Body.Close() // no exception handling in tests
}

/*

=== RUN   TestBlockIndefinitelyWithCancel
    block_test.go:78: Client timeout
    block_test.go:66: Canceling Request
--- PASS: TestBlockIndefinitelyWithCancel (3.00s)

*/
