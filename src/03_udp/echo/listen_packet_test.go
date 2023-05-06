// go test -v -race listen_packet_test.go echo.go

package echo

import (
	"bytes"
	"context"
	"net"
	"testing"
)

func TestListenPacketUDP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	serverAddr, err := echoServerUDP(ctx, "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()
	t.Logf("Server: %s", serverAddr)

	client, err := net.ListenPacket("udp", "127.0.0.1:")
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
	_, err = client.WriteTo(ping, serverAddr)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)
	n, addr, err := client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Client: Read '%s' from %q", buf[:n], addr)

	if !bytes.Equal(interrupt, buf[:n]) {
		t.Fatalf("expected %q, got %q", interrupt, buf[:n])
	}

	if addr.String() != interloper.LocalAddr().String() {
		t.Errorf("expected message from %q; actual sender is %q", interloper.LocalAddr(), addr)
	}

	n, addr, err = client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Client: Read '%s' from %q", buf[:n], addr)

	if !bytes.Equal(ping, buf[:n]) {
		t.Errorf("expected %q, got %q", ping, buf[:n])
	}

	if addr.String() != serverAddr.String() {
		t.Errorf("expected message from %q; actual sender is %q", serverAddr, addr)
	}

}

/*

=== RUN   TestListenPacketUDP
    listen_packet_test.go:19: Server: 127.0.0.1:61028
Read Buffer: Wait...
    listen_packet_test.go:30: Client: 127.0.0.1:62191
Context: Wait...
    listen_packet_test.go:37: Interloper: 127.0.0.1:64202
Write Buffer: Wait...
    listen_packet_test.go:62: Client: Read 'pardon me' from "127.0.0.1:64202"
Read Buffer: Wait...
    listen_packet_test.go:77: Client: Read 'ping' from "127.0.0.1:61028"
    listen_packet_test.go:26: Client: Close
Context: Done
--- PASS: TestListenPacketUDP (0.00s)
End Read Buffer
Connection: Close
PASS
ok  	command-line-arguments	0.297s

*/
