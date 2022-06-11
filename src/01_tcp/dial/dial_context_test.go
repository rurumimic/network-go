// go test -v dial_context_test.go

package dial

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestDialContext(t *testing.T) {
	deadline := time.Now().Add(time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel() // garbage control

	var d net.Dialer // DialContext is Dialer's method
	d.Control = func(_, _ string, _ syscall.RawConn) error {
		// timeout:
		time.Sleep(1*time.Second + time.Millisecond)

		// or

		// no timeout
		// time.Sleep(time.Millisecond)
		return nil
	}

	conn, err := d.DialContext(ctx, "tcp", "142.250.207.36:443") // google.com
	if err == nil {
		conn.Close()
		t.Fatal("Connection did not time out")
	}

	nErr, ok := err.(net.Error)
	if !ok {
		t.Error(err)
	} else {
		if nErr.Timeout() {
			t.Logf("Connection did time out: %v", err)
		}
		if !nErr.Timeout() {
			t.Errorf("Error is not a timeout: %v", err)
		}
	}

	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("Expected deadline exceeded; actual: %v", ctx.Err())
	}
}

/*

=== RUN   TestDialContext
    dial_context_test.go:41: Connection did time out: dial tcp 142.250.207.36:443: i/o timeout
--- PASS: TestDialContext (1.00s)
PASS
ok  	command-line-arguments	1.257s


=== RUN   TestDialContext
    dial_context_test.go:33: Connection did not time out
--- FAIL: TestDialContext (0.04s)
FAIL
FAIL	command-line-arguments	0.294s
FAIL

*/
