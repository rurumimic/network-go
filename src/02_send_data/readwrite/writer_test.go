package readwrite

import (
	"crypto/rand"
	"net"
	"testing"
)

func TestWriter(t *testing.T) {
	payload := make([]byte, 1<<24) // 16 MB
	_, err := rand.Read(payload)
	if err != nil {
		t.Fatal(err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:5678")
	if err != nil {
		t.Fatal(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		t.Log(err)
		return
	}

	_, err = conn.Write(payload)
	if err != nil {
		t.Error(err)
	}

	conn.Close()
}

/*

go test -v writer_test.go

=== RUN   TestWriter
--- PASS: TestWriter (14.77s)
PASS
ok      command-line-arguments  15.378s

*/
