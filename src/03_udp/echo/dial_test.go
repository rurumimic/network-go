// go test -v dial_test.go echo.go

package echo

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"
)

func TestDialUDP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	serverAddr, err := echoServerUDP(ctx, "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()
	t.Logf("Server: %s", serverAddr)

	client, err := net.Dial("udp", serverAddr.String())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		t.Log("Client: Close")
		_ = client.Close()
	}()

	t.Logf("Client: %s", client.LocalAddr())

	interloper, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Interloper: %s", interloper.LocalAddr())

	interrupt := []byte("pardon me")
	n, err := interloper.WriteTo(interrupt, client.LocalAddr())
	if err != nil {
		t.Fatal(err)
	}
	_ = interloper.Close()

	if l := len(interrupt); l != n {
		t.Fatalf("wrote %d bytes of %d", n, l)
	}

	ping := []byte("ping")
	_, err = client.Write(ping)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)
	n, err = client.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Client: Read '%s' from %q", buf[:n], serverAddr.String())

	if !bytes.Equal(ping, buf[:n]) {
		t.Errorf("expected %q, got %q", ping, buf[:n])
	}

	err = client.SetDeadline(time.Now().Add(time.Second))
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Read(buf)
	if err == nil {
		t.Fatal("unexpected packet")
	}

}

/*

=== RUN   TestDialUDP
    dial_test.go:20: Server: 127.0.0.1:49701
Read Buffer: Wait...
Context: Wait...
    dial_test.go:31: Client: 127.0.0.1:55294
    dial_test.go:38: Interloper: 127.0.0.1:59914
Write Buffer: Wait...
Read Buffer: Wait...
    dial_test.go:63: Client: Read 'ping' from "127.0.0.1:49701"
    dial_test.go:27: Client: Close
--- PASS: TestDialUDP (1.00s)
PASS
Context: Done
ok  	command-line-arguments	1.264s

*/
