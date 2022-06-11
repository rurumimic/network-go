// go test -v deadline_test.go
package main

import (
	"io"
	"net"
	"testing"
	"time"
)

func TestDeadline(t *testing.T) {
	sync := make(chan struct{})

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}
		defer func() {
			conn.Close()
			close(sync)
		}()

		err = conn.SetDeadline(time.Now().Add(5 * time.Second))
		// or timeout
		// err = conn.SetDeadline(time.Now())
		if err != nil {
			t.Error(err)
			return
		}

		buf := make([]byte, 1)
		_, err = conn.Read(buf) // blocking...
		t.Log(buf)              // [0]
		nErr, ok := err.(net.Error)
		if !ok || !nErr.Timeout() {
			t.Errorf("Expected timeout error; actual: %v", err)
		}

		sync <- struct{}{}

		err = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			t.Error(err)
			return
		}

		_, err = conn.Read(buf)
		t.Log(buf) // [65] == A
		if err != nil {
			t.Error(err)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	<-sync
	_, err = conn.Write([]byte("A"))
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1)
	_, err = conn.Read(buf) // blocking...
	t.Log(buf)              // [0]
	if err != io.EOF {
		t.Errorf("Expected server termination; actual: %v", err)
	}
}

/*

=== RUN   TestDeadline
    deadline_test.go:38: [0]
    deadline_test.go:53: [65]
    deadline_test.go:73: [0]
--- PASS: TestDeadline (5.00s)
PASS
ok  	command-line-arguments	5.375s


=== RUN   TestDeadline
    deadline_test.go:39: [0]
    deadline_test.go:55: [0]
    deadline_test.go:57: read tcp 127.0.0.1:50978->127.0.0.1:50979: i/o timeout
    deadline_test.go:75: [0]
--- FAIL: TestDeadline (0.00s)
FAIL
FAIL	command-line-arguments	0.265s
FAIL

*/
