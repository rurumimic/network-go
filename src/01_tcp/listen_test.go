// go test -v listen_test.go

package main

import (
	"net"
	"testing"
)

func TestListener(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")

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
    listen_test.go:21: bound to "127.0.0.1:56411"
--- PASS: TestListener (0.00s)
PASS
ok  	command-line-arguments	0.454s

*/
