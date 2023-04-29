// go test -v monitor_test.go

package send_data

import (
	"io"
	"log"
	"net"
	"os"
	"testing"
)

type Monitor struct {
	*log.Logger
}

func (m *Monitor) Write(p []byte) (int, error) {
	err := m.Output(2, string(p))
	if err != nil {
		log.Println(err)
	}
	return len(p), nil
}

func TestMonitor(t *testing.T) {
	monitor := &Monitor{Logger: log.New(os.Stdout, "[MONITOR] ", 0)}

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		monitor.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)

		conn, err := listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		b := make([]byte, 1024)
		r := io.TeeReader(conn, monitor)

		n, err := r.Read(b)
		if err != nil && err != io.EOF {
			monitor.Println(err)
			return
		}

		w := io.MultiWriter(conn, monitor)

		_, err = w.Write(b[:n]) // echo back
		if err != nil && err != io.EOF {
			monitor.Println(err)
			return
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		monitor.Fatal(err)
	}

	_, err = conn.Write([]byte("Test 1357\n"))
	if err != nil {
		monitor.Fatal(err)
	}

	_ = conn.Close()
	<-done
}

/*

=== RUN   TestMonitor
[MONITOR] Test 1357
[MONITOR] Test 1357
--- PASS: TestMonitor (0.00s)
PASS
ok  	command-line-arguments	0.364s

*/
