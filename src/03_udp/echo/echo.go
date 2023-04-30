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
			<-ctx.Done()
			_ = s.Close()
		}()

		buf := make([]byte, 1024)

		for {
			n, clientAddr, err := s.ReadFrom(buf) // client -> server
			// _, clientAddr, err := s.ReadFrom(buf) // client -> server
			if err != nil {
				return
			}

			_, err = s.WriteTo(buf[:n], clientAddr) // server -> client
			// _, err = s.WriteTo([]bytes("cake"), clientAddr) // server -> client
			if err != nil {
				return
			}
		}
	}()

	return s.LocalAddr(), nil
}
