// go test -v dial_cancel_test.go
package dial

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestDialContextCancel(t *testing.T) {
	// with cancel
	ctx, cancel := context.WithCancel(context.Background())

	// or not cancel
	// ctx, _ := context.WithCancel(context.Background())
	sync := make(chan struct{})

	go func() {
		defer func() {
			sync <- struct{}{}
		}()

		var d net.Dialer // DialContext is Dialer's method
		d.Control = func(_, _ string, _ syscall.RawConn) error {
			time.Sleep(time.Second)
			return nil
		}

		conn, err := d.DialContext(ctx, "tcp", "142.250.207.36:443")
		if err != nil {
			t.Log(err) // dial tcp 142.250.207.36:443: operation was canceled
			return
		}

		conn.Close()
		t.Error("Connection did not time out")
	}()

	cancel()
	<-sync

	if ctx.Err() != context.Canceled {
		t.Errorf("Expected canceled context; actual: %v", ctx.Err())
	}
}

/*

=== RUN   TestDialContextCancel
    dial_cancel_test.go:29: dial tcp 142.250.207.36:443: operation was canceled
--- PASS: TestDialContextCancel (0.00s)
PASS
ok  	command-line-arguments	0.363s


=== RUN   TestDialContextCancel
    dial_cancel_test.go:37: Connection did not time out
    dial_cancel_test.go:44: Expected canceled context; actual: <nil>
--- FAIL: TestDialContextCancel (0.04s)
FAIL
FAIL	command-line-arguments	0.408s
FAIL

*/
