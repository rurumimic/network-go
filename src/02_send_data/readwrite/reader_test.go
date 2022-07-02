package readwrite

import (
	"io"
	"log"
	"net"
	"testing"
)

func TestReader(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:5678")
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1<<19) // 512KB

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}

		log.Printf("read %d bytes", n) // buf[:n] == Read Data from `conn`
	}

	conn.Close()

}

/*

go test -v reader_test.go

=== RUN   TestReader
2022/06/25 15:27:04 read 65328 bytes
2022/06/25 15:27:04 read 196348 bytes
2022/06/25 15:27:04 read 359616 bytes
2022/06/25 15:27:04 read 384452 bytes
2022/06/25 15:27:04 read 106600 bytes
2022/06/25 15:27:04 read 409548 bytes
2022/06/25 15:27:04 read 277644 bytes
2022/06/25 15:27:04 read 213408 bytes
2022/06/25 15:27:04 read 327472 bytes
2022/06/25 15:27:04 read 163580 bytes
2022/06/25 15:27:04 read 409548 bytes
2022/06/25 15:27:04 read 163736 bytes
2022/06/25 15:27:04 read 98200 bytes
2022/06/25 15:27:04 read 208 bytes
2022/06/25 15:27:04 read 114168 bytes
2022/06/25 15:27:04 read 65328 bytes
2022/06/25 15:27:04 read 208 bytes
2022/06/25 15:27:04 read 16332 bytes
2022/06/25 15:27:04 read 48996 bytes
2022/06/25 15:27:04 read 32872 bytes
2022/06/25 15:27:04 read 32872 bytes
2022/06/25 15:27:04 read 65328 bytes
2022/06/25 15:27:04 read 208 bytes
2022/06/25 15:27:04 read 65328 bytes
2022/06/25 15:27:04 read 49412 bytes
2022/06/25 15:27:04 read 408976 bytes
2022/06/25 15:27:04 read 342972 bytes
2022/06/25 15:27:04 read 148080 bytes
2022/06/25 15:27:04 read 408300 bytes
2022/06/25 15:27:04 read 123712 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 46812 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 350956 bytes
2022/06/25 15:27:04 read 106288 bytes
2022/06/25 15:27:04 read 114324 bytes
2022/06/25 15:27:04 read 98408 bytes
2022/06/25 15:27:04 read 64912 bytes
2022/06/25 15:27:04 read 81660 bytes
2022/06/25 15:27:04 read 75184 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 243784 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 334260 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 524288 bytes
2022/06/25 15:27:04 read 200376 bytes
--- PASS: TestReader (0.01s)
PASS
ok      command-line-arguments  0.127s

*/
