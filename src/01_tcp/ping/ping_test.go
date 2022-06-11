// go test -v ping.go ping_test.go
package ping

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"
)

func TestPinger(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	r, w := io.Pipe()
	done := make(chan struct{})
	resetTimer := make(chan time.Duration, 1) // interval channel
	resetTimer <- time.Second                 // init interval

	go func() {
		Pinger(ctx, w, resetTimer)
		close(done)
	}()

	receivePing := func(d time.Duration, r io.Reader) {
		if d >= 0 {
			fmt.Printf("Resetting timer (%s)\n", d)
			resetTimer <- d
		}

		now := time.Now()
		buf := make([]byte, 1024)
		n, err := r.Read(buf) // blocking...
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Received %q (%s)\n", buf[:n], time.Since(now).Round(100*time.Millisecond))
	}

	for i, v := range []int64{0, 200, 300, 0, -1, -1, -1} {
		fmt.Printf("Run %d\n", i+1)
		receivePing(time.Duration(v)*time.Millisecond, r)
	}

	cancel()
	<-done
}

/*

=== RUN   TestPinger
Run 1
Resetting timer (0s)
Received "ping" (1s)
Run 2
Resetting timer (200ms)
Received "ping" (200ms)
Run 3
Resetting timer (300ms)
Received "ping" (300ms)
Run 4
Resetting timer (0s)
Received "ping" (300ms)
Run 5
Received "ping" (300ms)
Run 6
Received "ping" (300ms)
Run 7
Received "ping" (300ms)
--- PASS: TestPinger (2.71s)
PASS
ok  	command-line-arguments	3.077s

*/
