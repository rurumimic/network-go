# Send Data

_back to [/README.md](/README.md)_

---

- [net](https://pkg.go.dev/net)
  - net.[Conn](https://pkg.go.dev/net#Conn)
- [io](https://pkg.go.dev/io)
  - io.[Reader](https://pkg.go.dev/io#Reader)
  - io.[Writer](https://pkg.go.dev/io#Writer)
  - io.[ReadWriteCloser](https://pkg.go.dev/io#ReadWriteCloser)
  - io.[Copy](https://pkg.go.dev/io#Copy)
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
