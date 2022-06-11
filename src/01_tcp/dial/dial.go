package dial

import (
	"fmt"
	"io"
	"net"
	"os"
)

func Starter() net.Listener {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return listener
}

func Accepter(listener net.Listener, done chan<- struct{}) {
	defer func() {
		fmt.Println("Close Accepter")
		done <- struct{}{}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		go Receiver(conn, done)
	}
}

func Receiver(c net.Conn, done chan<- struct{}) {
	defer func() {
		fmt.Println("Close Receiver")
		c.Close()
		done <- struct{}{}
	}()

	fmt.Println("Reading...")

	buf := make([]byte, 1024)
	for {
		n, err := c.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Println("Read EOF")
			}
			return
		}
		fmt.Printf("Received: %q", buf[:n])
	}
}
