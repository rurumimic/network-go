// go test -v -race echo_test.go echo.go

package echo

import (
	"context"
	"fmt"
	"net"
)

func echoServerUDP(ctx context.Context, addr string) (net.Addr, error) {
	s, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("binding to udp %s: %w", addr, err)
	}

	go func() {
		go func() {
			fmt.Println("Context: Wait...")
			<-ctx.Done()
			fmt.Println("Context: Done")
			_ = s.Close()
			fmt.Println("Connection: Close")
		}()

		buf := make([]byte, 1024)

		for {
			fmt.Println("Read Buffer: Wait...")
			n, clientAddr, err := s.ReadFrom(buf) // client -> server
			// _, clientAddr, err := s.ReadFrom(buf) // client -> server
			if err != nil {
				fmt.Println("End Read Buffer")
				return
			}

			fmt.Println("Write Buffer: Wait...")
			_, err = s.WriteTo(buf[:n], clientAddr) // server -> client
			// _, err = s.WriteTo([]byte("cake"), clientAddr) // server -> client
			if err != nil {
				fmt.Println("End Write Buffer")
				return
			}
		}
	}()

	return s.LocalAddr(), nil
}

/*

=== RUN   TestEchoServerUDP
Read Buffer: Wait...
Write Buffer: Wait...
Read Buffer: Wait...
Context: Wait...
    echo_test.go:49: expected "ping", got "ping"
    echo_test.go:51: Client: End
    echo_test.go:25: Client: Close
Context: Done
--- PASS: TestEchoServerUDP (0.00s)
End Read Buffer
Connection: Close
PASS
ok  	command-line-arguments	0.172s

*/
