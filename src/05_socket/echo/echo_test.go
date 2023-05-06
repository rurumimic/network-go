// go test -v -race echo_test.go echo.go

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

func TestEchoServerUnix(t *testing.T) {
	dir, err := os.MkdirTemp("", "echo_unix")
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
	socket := filepath.Join(dir, fmt.Sprintf("%d.sock", os.Getpid()))
	rAddr, err := streamingEchoServer(ctx, "unix", socket)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	err = os.Chmod(socket, os.ModeSocket|0666)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("listening on %s", rAddr.String())
	t.Logf("socket file: %s", socket)

	conn, err := net.Dial("unix", rAddr.String())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = conn.Close()
	}()

	msg := []byte("ping")
	for i := 0; i < 3; i++ {
		_, err = conn.Write(msg)
		if err != nil {
			t.Fatal(err)
		}
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	expected := bytes.Repeat(msg, 3)
	if !bytes.Equal(expected, buf[:n]) {
		t.Fatalf("expected %q, got %q", expected, buf[:n])
	}

	t.Logf("received %q", buf[:n]) // pingpingping
}

/*

=== RUN   TestEchoServerUnix
    echo_test.go:26: using temp dir: /var/folders/8r/00xknx_d2194lrv4p873m97h0000gn/T/echo_unix1155505486
    echo_test.go:41: listening on /var/folders/8r/00xknx_d2194lrv4p873m97h0000gn/T/echo_unix1155505486/79352.sock
    echo_test.go:42: socket file: /var/folders/8r/00xknx_d2194lrv4p873m97h0000gn/T/echo_unix1155505486/79352.sock
    echo_test.go:71: received "cakeccakeccakec"
--- PASS: TestEchoServerUnix (0.00s)
PASS
ok  	command-line-arguments	0.308s

*/
