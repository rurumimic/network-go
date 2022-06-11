// go test -v ping.go ping_heartbeat_test.go
package ping

import (
	"context"
	"io"
	"net"
	"testing"
	"time"
)

func TestHeartbeatPinger(t *testing.T) {
	done := make(chan struct{})
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	begin := time.Now()
	go func() {
		defer func() {
			close(done)
		}()

		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer func() {
			cancel()
			conn.Close()
		}()

		resetTimer := make(chan time.Duration, 1)
		resetTimer <- time.Second
		go Pinger(ctx, conn, resetTimer)

		err = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			t.Error(err)
			return
		}

		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				// t.Error(err)
				t.Logf("Error: %v", err)
				return
			}
			t.Logf("SERVER: [%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])

			// reset Timer and delay deadline
			resetTimer <- 0
			err = conn.SetDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				t.Error(err)
				return
			}
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for i := 0; i < 4; i++ { // Read 4 Pings
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("CLIENT: [%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])
	}
	_, err = conn.Write([]byte("PONG!!!")) // Send Message to server
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 4; i++ { // Read 4 Pings
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Fatal(err)
			}
			break
		}
		t.Logf("CLIENT: [%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])
	}

	<-done

	end := time.Since(begin).Truncate(time.Second)
	t.Logf("CLIENT: [%s] done", end)
	if end != 9*time.Second {
		t.Fatalf("expected EOF at 9 seconds; actual %s", end)
	}

}

/*

=== RUN   TestHeartbeatPinger
    ping_heartbeat_test.go:78: CLIENT: [1s] ping
    ping_heartbeat_test.go:78: CLIENT: [2s] ping
    ping_heartbeat_test.go:78: CLIENT: [3s] ping
    ping_heartbeat_test.go:78: CLIENT: [4s] ping
    ping_heartbeat_test.go:54: SERVER: [4s] PONG!!!
    ping_heartbeat_test.go:93: CLIENT: [5s] ping
    ping_heartbeat_test.go:93: CLIENT: [6s] ping
    ping_heartbeat_test.go:93: CLIENT: [7s] ping
    ping_heartbeat_test.go:93: CLIENT: [8s] ping
    ping_heartbeat_test.go:51: Error: read tcp 127.0.0.1:51482->127.0.0.1:51483: i/o timeout
    ping_heartbeat_test.go:99: CLIENT: [9s] done
--- PASS: TestHeartbeatPinger (9.01s)
PASS
ok  	command-line-arguments	9.267s

*/
