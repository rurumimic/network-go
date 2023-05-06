// go test -v -race echo_posix_test.go echo.go

package echo

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestEchoServerUnixDatagram(t *testing.T) {
	dir, err := os.MkdirTemp("", "echo_unixgram")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if rErr := os.RemoveAll(dir); rErr != nil {
			t.Error(rErr)
		}
	}()

	t.Logf("using temp dir: %s", dir)

	ctx, cancel := context.WithCancel(context.Background())
	sSocket := filepath.Join(dir, fmt.Sprintf("s%d.sock", os.Getpid()))
	serverAddr, err := datagramEchoServer(ctx, "unixgram", sSocket)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	err = os.Chmod(sSocket, os.ModeSocket|0622)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("server listening on %s", serverAddr.String())
	t.Logf("server socket file: %s", sSocket)

	cSocket := filepath.Join(dir, fmt.Sprintf("c%d.sock", os.Getpid()))
	client, err := net.ListenPacket("unixgram", cSocket)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = client.Close()
	}()

	err = os.Chmod(sSocket, os.ModeSocket|0622)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("client socket file: %s", cSocket)

	msg := []byte("ping")
	for i := 0; i < 3; i++ {
		_, err = client.WriteTo(msg, serverAddr)
		if err != nil {
			t.Fatal(err)
		}
	}

	buf := make([]byte, 1024)
	for i := 0; i < 3; i++ {
		n, addr, err := client.ReadFrom(buf)
		if err != nil {
			t.Fatal(err)
		}

		if addr.String() != serverAddr.String() {
			t.Fatalf("received reply from %q instead of %q", addr, serverAddr)
		}

		if !bytes.Equal(msg, buf[:n]) {
			t.Fatalf("expected %q, got %q", msg, buf[:n])
		}

		t.Logf("received %q", buf[:n]) // ping
	}

}

/*

=== RUN   TestEchoServerUnixDatagram
    echo_posix_test.go:26: using temp dir: /var/folders/8r/00xknx_d2194lrv4p873m97h0000gn/T/echo_unixgram1936980904
    echo_posix_test.go:41: server listening on /var/folders/8r/00xknx_d2194lrv4p873m97h0000gn/T/echo_unixgram1936980904/s82744.sock
    echo_posix_test.go:42: server socket file: /var/folders/8r/00xknx_d2194lrv4p873m97h0000gn/T/echo_unixgram1936980904/s82744.sock
    echo_posix_test.go:58: client socket file: /var/folders/8r/00xknx_d2194lrv4p873m97h0000gn/T/echo_unixgram1936980904/c82744.sock
    echo_posix_test.go:83: received "ping"
    echo_posix_test.go:83: received "ping"
    echo_posix_test.go:83: received "ping"
--- PASS: TestEchoServerUnixDatagram (0.00s)
PASS
ok  	command-line-arguments	0.406s

*/
