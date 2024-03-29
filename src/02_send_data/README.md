# Send Data

_back to [/README.md](/README.md)_

---

- [net](https://pkg.go.dev/net)
  - net.[Conn](https://pkg.go.dev/net#Conn)
  - net.[TCPConn](https://pkg.go.dev/net#TCPConn)
- [io](https://pkg.go.dev/io)
  - io.[Reader](https://pkg.go.dev/io#Reader)
  - io.[Writer](https://pkg.go.dev/io#Writer)
  - io.[ReadWriteCloser](https://pkg.go.dev/io#ReadWriteCloser)
  - io.[Copy](https://pkg.go.dev/io#Copy)
  - io.[TeeReader](https://pkg.go.dev/io#TeeReader)
  - io.[MultiWriter](https://pkg.go.dev/io#MultiWriter)
- [bufio](https://pkg.go.dev/bufio)
  - bufio.[Scanner](https://pkg.go.dev/bufio#Scanner)
- iota
  - [test/iota.go](https://go.dev/test/iota.go)
- [reflect](https://pkg.go.dev/reflect)
  - reflect.[DeepEqual](https://pkg.go.dev/reflect#DeepEqual)
- [fmt](https://pkg.go.dev/fmt)
  - [Printing](https://pkg.go.dev/fmt#hdr-Printing)
- go test [-race](https://go.dev/doc/articles/race_detector)

---

## net

### net.Conn.Read & net.Conn.Write

#### examples

##### Read Test

- [read_test.go](read_test.go)
- [readwrite/reader_test.go](readwrite/reader_test.go)
- [readwrite/writer_test.go](readwrite/writer_test.go)

###### run

```bash
go test -v read_test.go

# readwrite
go test -v writer_test.go
go test -v reader_test.go
```

###### write

```go
payload := make([]byte, 1<<24)

_, err = conn.Write(payload)
if err != nil {
  t.Error(err)
}
```

###### read

```go
buf := make([]byte, 1<<19)

for {
  n, err := conn.Read(buf)
  if err != nil {
    if err != io.EOF {
      t.Error(err)
    }
    break
  }
  t.Logf("read %d bytes", n)
}
```

---

## bufio

### bufio.Scanner

#### examples

##### Scanner

- [scanner_test.go](scanner_test.go)

```go
scanner := bufio.NewScanner(conn)
scanner.Split(bufio.ScanWords)

var words []string

for scanner.Scan() {
  words = append(words, scanner.Text())
}

err = scanner.Err()
if err != nil {
  t.Error(err)
}
```

---

## TLV, Type-Length-Value: encoding system

#### examples

- [types.go](types.go)
- [types_test.go](types_test.go)

##### Read and Write

```go
const (
	BinaryType uint8 = iota + 1
	StringType
	MaxPayloadSize uint32 = 10 << 20 // 10 MB
)

type Payload interface {
	fmt.Stringer
	io.ReaderFrom
	io.WriterTo
	Bytes() []byte
}
```

```go
type Binary []byte
func (m Binary) Bytes() []byte
func (m Binary) String() string
func (m Binary) WriteTo(w io.Writer) (int64, error)
func (m *Binary) ReadFrom(r io.Reader) (int64, error)
```

```go
func decode(r io.Reader) (Payload, error)
```

##### TCP Data Send

```go
conn, err := listener.Accept()

p := Payload(&Binary("Don't panic."))
_, err = p.WriteTo(conn)
```

```go
conn, err := net.Dial("tcp", listener.Addr().String())
actual, err := decode(conn) // "Don't panic."
```

```go
t.Logf("[%T] %[1]q", actual)
```

- `%T`: a Go-syntax representation of the type of the value
- `%q`: a double-quoted string safely escaped with Go syntax
- notation `[n]`: immediately before the verb indicates that the nth one-indexed argument

---

## Error Handling

- [ping.go](ping/ping.go)

```bash
nc -kl 127.0.0.1 8888
```

```bash
go run ping.go -c 10 127.0.0.1:8888

PING 127.0.0.1:8888
1 414.179µs
2 561.307µs
3 621.029µs
```

```go
var (
  err error
  n   int
  i   = 7 // max retry
)

for ; i > 0; i-- {
  n, err = conn.Write([]byte("hello world"))
  if err != nil {
    // check net.Error interface implementation and Temporary error
    if nErr, ok := err.(net.Error); ok && nErr.Temporary() {
      log.Println("temporary error:", nErr)
      time.Sleep(10 * time.second)
      continue
    }
    return err
  }
  break
}

if i == 0 {
  return erros.New("temporary write failure threshold exceeded")
}

log.Printf("wrote %d bytes to %s\n", n, conn.RemoteAddr())
```

---

## Proxy

### io.Copy

#### examples

##### Copy

- [proxy_test.go](proxy_test.go)
- go test [-race](https://go.dev/doc/articles/race_detector)

```go
func proxy(from io.Reader, to io.Writer) {
  _, err := io.Copy(to, from) // echo -> proxy
  _, _ = io.Copy(fromWriter, toReader) // proxy -> echo
}

// from:proxy     to: echo
// from net.Conn, to net.Conn
proxy(from, to)
```

- `net.Conn` ⟷ `io.Reader`
- `net.Conn` ⟷ `io.Writer`

---

## Monitor

### io.TeeReader, io.MultiWriter

#### examples

- [monitor_test.go](monitor_test.go)

---

## TCPConn

### net.TCPConn

#### examples

- [tcp_conn_test.go](tcp_conn_test.go)

```go
tcpConn, ok := conn.(*net.TCPConn)
```

##### KeepAlive

```go
err := tcpConn.SetKeepAlive(true)
err := tcpConn.SetKeepAlivePeriod(time.Minute)
```

##### Linger

```go
err := tcpConn.SetLinger(-1) // -1: no timeout, 0: close immediately, >0: timeout
```

##### Buffer

```go
if err := tcpConn.SetReadBuffer(212992); err != nil {
  return err
}

if err := tcpConn.SetWriteBuffer(212992); err != nil {
  return err
}
```

##### Zero Window Error

```go
buf := make([]byte, 1024)

for {
  n, err := conn.Read(buf) // 1. recieve buffer is empty
  if err != nil {
    return err
  }

  handle(buf[:n]) // 2. blocking...
  // 3. recieve buffer is full
}
```

##### CLOSE_WAIT

When a TCP socket is CLOSE_WAIT state and client send data other than FIN packet,  
server send RST and the connection is closed.

```go
for {
  conn, err := listener.Accept()
  if err != nil {
    return err
  }

  go func(c net.Conn) {
    // must
    // c.Close()
    buf := make([]byte, 1024)

    for {
      n, err := c.Read(buf)
      if err != nil {
        return // c.Close() not called -> CLOSE_WAIT
      }

      handle(buf[:n])
    }
  }(conn)
}
```
