// go test -v echo_test.go echo.go

package echo

import (
	"bytes"
	"context"
	"net"
	"testing"
)

func TestEchoServerUDP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	serverAddr, err := echoServerUDP(ctx, "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	client, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = client.Close()
	}()

	msg := []byte("ping")
	_, err = client.WriteTo(msg, serverAddr)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)
	n, addr, err := client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}

	if addr.String() != serverAddr.String() {
		t.Fatalf("expected %s, got %s", serverAddr, addr)
	}

	if !bytes.Equal(msg, buf[:n]) {
		t.Fatalf("expected %q, got %q", msg, buf[:n])
	}

}

/*

=== RUN   TestEchoServerUDP
--- PASS: TestEchoServerUDP (0.00s)
PASS
ok  	command-line-arguments	0.353s

when wrong:

=== RUN   TestEchoServerUDP
    echo_test.go:45: expected "ping", got "cake"
--- FAIL: TestEchoServerUDP (0.00s)
FAIL
FAIL	command-line-arguments	0.366s
FAIL

*/
