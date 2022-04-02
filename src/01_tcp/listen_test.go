// go test listen_test.go -v

package main

import (
	"net"
	"testing"
)

func TestListener(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	// listener, err := net.Listen("tcp", ":")
	// listener, err := net.Listen("tcp", "192.168.101.101:8888")

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = listener.Close()
	}()

	t.Logf("bound to %q", listener.Addr())
}

/*

=== RUN   TestListener
    listen_test.go:23: bound to "[::]:56471"
--- PASS: TestListener (0.00s)
PASS
ok  	command-line-arguments	0.300s

=== RUN   TestListener
    listen_test.go:21: bound to "127.0.0.1:56411"
--- PASS: TestListener (0.00s)
PASS
ok  	command-line-arguments	0.454s

=== RUN   TestListener
    listen_test.go:16: listen tcp 192.168.101.101:8888: bind: can't assign requested address
--- FAIL: TestListener (0.00s)
FAIL
FAIL	command-line-arguments	0.298s
FAIL
*/
