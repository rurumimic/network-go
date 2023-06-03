// go test -v -timeout 10s time_test.go

package client

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestHeadTime(t *testing.T) {
	url := "https://www.time.gov/"

	t.Logf("curl -I %s", url)
	resp, err := http.Head(url)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close() // no exception handling in tests

	now := time.Now()
	t.Logf("Now: %s", now)
	date := resp.Header.Get("Date")
	if date == "" {
		t.Fatal("no Date header received fro time.gov")
	}
	t.Logf("Date: %s", date)

	dt, err := time.Parse(time.RFC1123, date)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("time.gov: %s (skew %s)", dt, now.Sub(dt))
}

/*

=== RUN   TestHeadTime
    time_test.go:14: curl -I https://www.time.gov/
    time_test.go:22: Now: 2023-06-03 13:11:30.376903 +0900 KST m=+1.601415411
    time_test.go:27: Date: Sat, 03 Jun 2023 04:11:30 GMT
    time_test.go:34: time.gov: 2023-06-03 04:11:30 +0000 GMT (skew 376.903ms)
--- PASS: TestHeadTime (1.60s)
PASS
ok  	command-line-arguments	1.884s

*/

func TestBodyCopy(t *testing.T) {
	url := "https://www.time.gov/"

	t.Logf("curl -i %s", url)
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	message := make([]byte, 70)
	resp.Body.Read(message)
	t.Logf("Reads 70 bytes: '%s'", message)

	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close() // no exception handling in tests
}

/*

=== RUN   TestBodyCopy
    body_test.go:14: curl -i https://www.time.gov/
    body_test.go:22: Reads 70 bytes: '<!DOCTYPE html>
        <html lang="en">
        	<head>
        		<meta charset="utf-8" />'
--- PASS: TestBodyCopy (0.82s)
PASS
ok  	command-line-arguments	1.106s

*/

func TestHeadTimeWithTimeout(t *testing.T) {
	url := "https://www.time.gov/"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil) // nil = HEAD no payload
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close() // no exception handling in tests
	cancel()              // No further need for the context

	now := time.Now().Round(time.Second)
	date := resp.Header.Get("Date")
	if date == "" {
		t.Fatal("no Date header received fro time.gov")
	}
	t.Logf("Date: %s", date)

	dt, err := time.Parse(time.RFC1123, date)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("time.gov: %s (skew %s)", dt, now.Sub(dt))
}

/*

=== RUN   TestHeadTimeWithTimeout
    time_test.go:104: Date: Sat, 03 Jun 2023 05:17:52 GMT
    time_test.go:111: time.gov: 2023-06-03 05:17:52 +0000 GMT (skew 0s)
--- PASS: TestHeadTimeWithTimeout (0.65s)

*/
