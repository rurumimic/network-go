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
		t.Log("Client: Close")
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

	t.Logf("expected %q, got %q", msg, buf[:n])

	t.Log("Client: End")

}

/*

=== RUN   TestEchoServerUDP
Read Buffer: Wait...
Write Buffer: Wait...
Context: Wait...
Read Buffer: Wait...
    echo_test.go:49: expected "ping", got "ping"
    echo_test.go:51: Client: End
    echo_test.go:25: Client: Close
--- PASS: TestEchoServerUDP (0.00s)
PASS
Context: Done
ok  	command-line-arguments	0.377s

when wrong:

=== RUN   TestEchoServerUDP
Read Buffer: Wait...
Write Buffer: Wait...
Context: Wait...
Read Buffer: Wait...
    echo_test.go:46: expected "ping", got "cake"
    echo_test.go:25: Client: Close
Context: Done
--- FAIL: TestEchoServerUDP (0.00s)
FAIL
End Read Buffer
FAIL	command-line-arguments	0.271s
FAIL

*/
