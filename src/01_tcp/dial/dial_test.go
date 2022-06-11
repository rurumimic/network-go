// go test -v dial_test.go dial.go
// go test -v -run TestDial

package dial

import (
	"net"
	"testing"
	"time"
)

func TestDial(t *testing.T) {
	listener := Starter()
	t.Logf("Bound to %q", listener.Addr())

	done := make(chan struct{})
	go Accepter(listener, done)
	time.Sleep(1 * time.Second)

	t.Logf("Dial to %q", listener.Addr())
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Log("Dial error")
		t.Fatal(err)
	}

	t.Log("Close Connection: Send FIN")
	conn.Close()
	<-done
	t.Log("Receiver closed")

	t.Log("Close Listener: Send FIN")
	listener.Close()
	<-done
	t.Log("Accepter closed")
}

/*

=== RUN   TestDial
    dial_test.go:18: Bound to "127.0.0.1:56758"
    dial_test.go:23: Dial to "127.0.0.1:56758"
Reading...
    dial_test.go:30: Close Connection: Send FIN
Read EOF
Close Receiver
    dial_test.go:33: Receiver closed
    dial_test.go:35: Close Listener: Send FIN
accept tcp4 127.0.0.1:56758: use of closed network connection
Close Accepter
    dial_test.go:38: Accepter closed
--- PASS: TestDial (1.00s)
PASS
ok  	command-line-arguments	1.363s

*/
