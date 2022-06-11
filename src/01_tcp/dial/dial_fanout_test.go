// go test -v dial_fanout_test.go

package dial

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"
)

func TestDialContextCancelFanOut(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))

	// or context.DeadExceeded
	// ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))

	// or nil <= case response<-id

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	go func() {
		conn, err := listener.Accept()
		if err == nil {
			conn.Close()
		}
	}()

	dial := func(ctx context.Context, address string, response chan int, id int, wg *sync.WaitGroup) {
		defer func() {
			wg.Done()
			// t.Logf("End: %d", id)
		}()

		var d net.Dialer
		c, err := d.DialContext(ctx, "tcp", address)
		if err != nil {
			// t.Logf("%d: %s", id, err)
			return
		}
		c.Close()

		// t.Logf("Close: %d", id)

		select {
		case <-ctx.Done():
			// t.Logf("Done: %d", id)
		case response <- id:
			// t.Logf("Response: %d", id)
		}
	}

	res := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go dial(ctx, listener.Addr().String(), res, i+1, &wg)
	}

	response := <-res
	// t.Logf("ID: %d", response)
	cancel()

	wg.Wait()
	close(res)

	if ctx.Err() != context.Canceled {
		t.Errorf("Expected canceled context; actual: %s", ctx.Err())
	}

	t.Logf("Dialer %d retrieved the resource", response)
}

/*

=== RUN   TestDialContextCancelFanOut
    dial_fanout_test.go:42: Close: 10
    dial_fanout_test.go:31: End 10
    dial_fanout_test.go:42: Close: 8
    dial_fanout_test.go:31: End 8
    dial_fanout_test.go:42: Close: 3
    dial_fanout_test.go:31: End 3
    dial_fanout_test.go:37: 2: dial tcp 127.0.0.1:58809: operation was canceled
    dial_fanout_test.go:37: 4: dial tcp 127.0.0.1:58809: operation was canceled
    dial_fanout_test.go:31: End 4
    dial_fanout_test.go:31: End 2
    dial_fanout_test.go:37: 9: dial tcp 127.0.0.1:58809: operation was canceled
    dial_fanout_test.go:31: End 9
    dial_fanout_test.go:37: 6: dial tcp 127.0.0.1:58809: operation was canceled
    dial_fanout_test.go:31: End 6
    dial_fanout_test.go:42: Close: 7
    dial_fanout_test.go:31: End 7
    dial_fanout_test.go:37: 5: dial tcp 127.0.0.1:58809: operation was canceled
    dial_fanout_test.go:31: End 5
    dial_fanout_test.go:37: 1: dial tcp 127.0.0.1:58809: operation was canceled
    dial_fanout_test.go:31: End 1
    dial_fanout_test.go:67: Dialer 10 retrieved the resource
--- PASS: TestDialContextCancelFanOut (0.00s)
PASS
ok  	dial	0.246s



=== RUN   TestDialContextCancelFanOut
    dial_fanout_test.go:50: Response: 1
    dial_fanout_test.go:63: ID: 1
    dial_fanout_test.go:48: Done: 8
    dial_fanout_test.go:48: Done: 6
    dial_fanout_test.go:48: Done: 2
    dial_fanout_test.go:48: Done: 10
    dial_fanout_test.go:48: Done: 4
    dial_fanout_test.go:48: Done: 3
    dial_fanout_test.go:48: Done: 7
    dial_fanout_test.go:48: Done: 5
    dial_fanout_test.go:48: Done: 9
    dial_fanout_test.go:70: Expected canceled context; actual: context deadline exceeded
    dial_fanout_test.go:73: Dialer 1 retrieved the resource
--- FAIL: TestDialContextCancelFanOut (10.00s)
FAIL
FAIL	command-line-arguments	10.371s
FAIL

*/
