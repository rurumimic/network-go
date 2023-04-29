// go test -race -v tcp_conn_test.go

package send_data

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

func startServer(done chan<- struct{}, addr *net.TCPAddr) {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer listener.Close()
	fmt.Println("listen...")

	tcpConn, err := listener.AcceptTCP()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer func() {
		tcpConn.Close()
		close(done)
	}()
	fmt.Println("accept...")

	done <- struct{}{}
	fmt.Println("...end server")
}

func TestConn(t *testing.T) {
	done := make(chan struct{})

	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:54321")
	if err != nil {
		t.Fatal(err)
	}

	go startServer(done, addr)
	conn, err := net.Dial("tcp", "127.0.0.1:54321")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		t.Fatal("conn is not a TCPConn")
	}
	if tcpConn == nil {
		t.Fatal("conn is nil")
	}
	tcpConn.Close()
	t.Logf("Conn is ok")
	<-done
	t.Log("...end client")
}

func TestTcpConn(t *testing.T) {
	done := make(chan struct{})

	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:54322")
	if err != nil {
		t.Fatal(err)
	}

	go startServer(done, addr)

	tcpConn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		t.Fatal(err)
	}
	if tcpConn == nil {
		t.Fatal("tcpConn is nil")
	}
	defer tcpConn.Close()

	t.Logf("TCPConn is ok")

	// err = tcpConn.SetKeepAlive(true)
	// or
	err = tcpConn.SetKeepAlivePeriod(time.Second * 1)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(3 * time.Second)

	<-done
	t.Log("...end client")
}

/*

=== RUN   TestConn
listen...
accept...
    tcp_conn_test.go:59: Conn is ok
    tcp_conn_test.go:61: ...end client
...end server
--- PASS: TestConn (0.00s)


=== RUN   TestTcpConn
listen...
accept...
    tcp_conn_test.go:83: TCPConn is ok
    tcp_conn_test.go:85: ...end client
...end server
--- PASS: TestTcpConn (1.00s)
PASS
ok  	command-line-arguments	1.166s

*/
