// go test -v -run TestDialTimeout

package dial

import (
	"net"
	"reflect"
	"testing"
	"time"
)

func TestDialTimeout(t *testing.T) {
	conn, err := DialTimeout("tcp", "10.0.0.0:http", 5*time.Second)
	if err == nil {
		conn.Close()
		t.Fatal("Connection did not time out")
	}

	nErr, ok := err.(net.Error)
	// nErr (*net.OpError) = dial tcp 10.0.0.0:80: lookup 10.0.0.0:80 on 127.0.0.1: Connection timed out
	// ok = true

	t.Logf("Error Type: %s", reflect.TypeOf(err)) // *net.OpError
	t.Logf("Is Timeout: %t", nErr.Timeout())      // true

	if !ok {
		t.Fatal(err)
	}

	if !nErr.Timeout() {
		t.Fatal("Error is not a timeout")
	}

	t.Log("ok")
}

/*

=== RUN   TestDialTimeout
    dial_timeout_test.go:22: Error Type: *net.OpError
    dial_timeout_test.go:23: Is Timeout: true
    dial_timeout_test.go:33: ok
--- PASS: TestDialTimeout (0.00s)
PASS
ok  	dial	0.485s

=== RUN   TestDialTimeout
    dial_timeout_test.go:22: Error Type: *net.OpError
    dial_timeout_test.go:23: Is Timeout: false
    dial_timeout_test.go:30: Error is not a timeout
--- FAIL: TestDialTimeout (0.01s)
FAIL
exit status 1
FAIL	dial	0.384s
*/
