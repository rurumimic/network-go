package dial

import (
	"net"
	"syscall"
	"time"
)

func DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{ // mimic net.Dialer interface
		Control: func(_, addr string, _ syscall.RawConn) error { // override Control
			return &net.DNSError{ // mimic DNS timeout error
				Err:         "Connection timed out",
				Name:        addr,
				Server:      "127.0.0.1",
				IsTimeout:   true, // or false
				IsTemporary: true,
			}
		},
		Timeout: timeout,
	}
	return d.Dial(network, address)
}
