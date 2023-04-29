// go test -v read_test.go

package send_data

import (
	"crypto/rand"
	"io"
	"net"
	"testing"
)

func TestReadIntoBuffer(t *testing.T) {

	//   24   20   16   12   8    4
	// 1 0000 0000 0000 0000 0000 0000
	payload := make([]byte, 1<<24) // 16 MB
	_, err := rand.Read(payload)
	if err != nil {
		t.Fatal(err)
	}

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

		defer conn.Close()

		_, err = conn.Write(payload)
		if err != nil {
			t.Error(err)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1<<19) // 512KB

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}

		t.Logf("read %d bytes", n) // buf[:n] == Read Data from `conn`
	}

	conn.Close()

}

/*

=== RUN   TestReadIntoBuffer
    read_test.go:58: read 16332 bytes
    read_test.go:58: read 375792 bytes
    read_test.go:58: read 524288 bytes
    read_test.go:58: read 7256 bytes
    read_test.go:58: read 302636 bytes
    read_test.go:58: read 524288 bytes
    read_test.go:58: read 350748 bytes
    read_test.go:58: read 74144 bytes
    read_test.go:58: read 482600 bytes
    read_test.go:58: read 342972 bytes
    read_test.go:58: read 107120 bytes
    read_test.go:58: read 458544 bytes
    read_test.go:58: read 524288 bytes
    read_test.go:58: read 391396 bytes
    read_test.go:58: read 213408 bytes
    read_test.go:58: read 524288 bytes
    read_test.go:58: read 497684 bytes
    read_test.go:58: read 74716 bytes
    read_test.go:58: read 524288 bytes
    read_test.go:58: read 383360 bytes
    read_test.go:58: read 221704 bytes
    read_test.go:58: read 524288 bytes
    read_test.go:58: read 342556 bytes
    read_test.go:58: read 229896 bytes
    read_test.go:58: read 343596 bytes
    read_test.go:58: read 359928 bytes
    read_test.go:58: read 32664 bytes
    read_test.go:58: read 524288 bytes
    read_test.go:58: read 524288 bytes
    read_test.go:58: read 129824 bytes
    read_test.go:58: read 507748 bytes
    read_test.go:58: read 473628 bytes
    read_test.go:58: read 107536 bytes
    read_test.go:58: read 524288 bytes
    read_test.go:58: read 524288 bytes
    read_test.go:58: read 113752 bytes
    read_test.go:58: read 524288 bytes
    read_test.go:58: read 391968 bytes
    read_test.go:58: read 188988 bytes
    read_test.go:58: read 524288 bytes
    read_test.go:58: read 524288 bytes
    read_test.go:58: read 113752 bytes
    read_test.go:58: read 524288 bytes
    read_test.go:58: read 391968 bytes
    read_test.go:58: read 188988 bytes
    read_test.go:58: read 524288 bytes
    read_test.go:58: read 524288 bytes
    read_test.go:58: read 113752 bytes
    read_test.go:58: read 57652 bytes
--- PASS: TestReadIntoBuffer (0.14s)
PASS
ok      command-line-arguments  0.388s

*/
