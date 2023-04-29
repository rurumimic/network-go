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
			t.Log("Wait ...")
			time.Sleep(time.Second)
			t.Log("... 1 Second")
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

with Cancel:

=== RUN   TestDialContextCancel
    dial_cancel_test.go:37: dial tcp 142.250.207.36:443: operation was canceled
--- PASS: TestDialContextCancel (0.00s)
PASS
ok  	command-line-arguments	0.274s


without Cancel:

=== RUN   TestDialContextCancel
    dial_cancel_test.go:29: Wait ...
    dial_cancel_test.go:31: ... 1 Second
    dial_cancel_test.go:42: Connection did not time out
    dial_cancel_test.go:49: Expected canceled context; actual: <nil>
--- FAIL: TestDialContextCancel (1.03s)
FAIL
FAIL	command-line-arguments	1.307s
FAIL

*/
