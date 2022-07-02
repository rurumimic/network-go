// go test -race -v proxy_test.go

package send_data

import (
	"io"
	"net"
	"sync"
	"testing"
)

func proxy(t *testing.T, from io.Reader, to io.Writer) error {
	fromWriter, fromIsWriter := from.(io.Writer)
	toReader, toIsReader := to.(io.Reader)

	if toIsReader && fromIsWriter {
		go func() {
			t.Logf("proxy %T -> echo %T", fromWriter, toReader)
			_, _ = io.Copy(fromWriter, toReader)
		}()
	}

	t.Logf("echo %T -> proxy %T", to, from)
	_, err := io.Copy(to, from)

	return err
}

////////////////////////////////////////////////////////////////////////////////

func pingpong(t *testing.T, c net.Conn) {
	defer c.Close()

	for {
		buf := make([]byte, 1024)
		n, err := c.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			} else {
				t.Log("Read EOF")
			}
			return
		}

		switch msg := string(buf[:n]); msg {
		case "ping":
			t.Logf("[Echo] Read %s", msg)
			_, err = c.Write([]byte("pong"))
		default:
			t.Logf("[Echo] Read %s", msg)
			_, err = c.Write(buf[:n])
		}

		if err != nil { // c.Write err check
			if err != io.EOF {
				t.Error(err)
			} else {
				t.Log("Write EOF")
			}
			return
		}
	}
}

func echoserver(t *testing.T, wg *sync.WaitGroup, server net.Listener) {
	defer wg.Done()

	for {
		conn, err := server.Accept()
		if err != nil {
			t.Log("Echo server close")
			return
		}

		go pingpong(t, conn)
	}

}

////////////////////////////////////////////////////////////////////////////////

func dial(t *testing.T, from net.Conn, server net.Listener) {
	defer from.Close() // End proxy server

	to, err := net.Dial("tcp", server.Addr().String()) // To echo server
	if err != nil {
		t.Error(err)
		return
	}

	defer to.Close() // End echo server

	t.Log("===PROXY START===")
	err = proxy(t, from, to) // from:proxy -> from:echo
	t.Log("===PROXY END=====")
	if err != nil && err != io.EOF {
		t.Error(err)
	}
}

func proxyserver(t *testing.T, wg *sync.WaitGroup, proxyServer net.Listener, server net.Listener) {
	defer wg.Done()

	for {
		conn, err := proxyServer.Accept()
		if err != nil {
			t.Log("Proxy server close")
			return
		}

		go dial(t, conn, server)
	}
}

////////////////////////////////////////////////////////////////////////////////

func TestProxy(t *testing.T) {
	var wg sync.WaitGroup

	// Setup Echo
	server, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	wg.Add(1)
	go echoserver(t, &wg, server)

	// Setup Proxy
	proxyServer, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	wg.Add(1)
	go proxyserver(t, &wg, proxyServer, server)

	// upstream -> downstream
	conn, err := net.Dial("tcp", proxyServer.Addr().String()) // To proxy server
	if err != nil {
		t.Fatal(err)
	}

	msgs := []struct{ Message, Reply string }{
		{"ping", "pong"},
		{"pong", "pong"},
		{"echo", "echo"},
		{"ping", "pong"},
	}

	for i, m := range msgs {
		_, err = conn.Write([]byte(m.Message)) // message -> proxy
		if err != nil {
			t.Fatal(err)
		}

		buf := make([]byte, 1024)

		n, err := conn.Read(buf) // proxy -> message
		if err != nil {
			t.Fatal(err)
		}

		actual := string(buf[:n])
		t.Logf("%q -> proxy -> %q", m.Message, actual)

		if actual != m.Reply {
			t.Errorf("%d: expoected reply: %q; actual %q", i, m.Reply, actual)
		}
	}

	_ = conn.Close()        // End connection to proxy
	_ = proxyServer.Close() // End proxy server
	_ = server.Close()      // End echo server

	wg.Wait()

}

/*

=== RUN   TestProxy
    proxy_test.go:94: ===PROXY START===
    proxy_test.go:23: echo *net.TCPConn -> proxy *net.TCPConn
    proxy_test.go:18: proxy *net.TCPConn -> echo *net.TCPConn
    proxy_test.go:48: [Echo] Read ping
    proxy_test.go:164: "ping" -> proxy -> "pong"
    proxy_test.go:51: [Echo] Read pong
    proxy_test.go:164: "pong" -> proxy -> "pong"
    proxy_test.go:51: [Echo] Read echo
    proxy_test.go:164: "echo" -> proxy -> "echo"
    proxy_test.go:48: [Echo] Read ping
    proxy_test.go:164: "ping" -> proxy -> "pong"
    proxy_test.go:96: ===PROXY END=====
    proxy_test.go:72: Echo server close
    proxy_test.go:108: Proxy server close
    proxy_test.go:41: Read EOF
--- PASS: TestProxy (0.00s)

*/

func TestEcho(t *testing.T) {
	var wg sync.WaitGroup

	// Setup Echo
	server, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	wg.Add(1)
	go echoserver(t, &wg, server)

	// upstream -> downstream
	conn, err := net.Dial("tcp", server.Addr().String()) // To echo server
	if err != nil {
		t.Fatal(err)
	}

	msgs := []struct{ Message, Reply string }{
		{"echo", "echo"},
		{"ping", "pong"},
	}

	for i, m := range msgs {
		_, err = conn.Write([]byte(m.Message)) // message -> proxy
		if err != nil {
			t.Fatal(err)
		}

		buf := make([]byte, 1024)

		n, err := conn.Read(buf) // proxy -> message
		if err != nil {
			t.Fatal(err)
		}

		actual := string(buf[:n])
		t.Logf("%q -> proxy -> %q", m.Message, actual)

		if actual != m.Reply {
			t.Errorf("%d: expoected reply: %q; actual %q", i, m.Reply, actual)
		}
	}

	_ = conn.Close()   // End connection to echo
	_ = server.Close() // End echo server

	wg.Wait()
}

/*

=== RUN   TestEcho
    proxy_test.go:51: [Echo] Read echo
    proxy_test.go:235: "echo" -> proxy -> "echo"
    proxy_test.go:48: [Echo] Read ping
    proxy_test.go:235: "ping" -> proxy -> "pong"
    proxy_test.go:41: Read EOF
    proxy_test.go:72: Echo server close
--- PASS: TestEcho (0.00s)

*/
