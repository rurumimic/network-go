# TCP

_back to [/README.md](/README.md)_

---

- net
  - net.Listen
- time
  - time.NewTimer
- Ping

---

## net

- [pkg net](https://pkg.go.dev/net)

```go
import "net"
```

### net.Listen

```go
func Listen(network, address string) (Listener, error)
```

- [func Listen](https://pkg.go.dev/net#Listen)
  - [type Listener](https://pkg.go.dev/net#Listener)
    - `Accept() (Conn, error)`: Accept waits for and returns the next connection to the listener.
      - [net/tcpsock.go](https://cs.opensource.google/go/go/+/refs/tags/go1.18.3:src/net/tcpsock.go;drc=386245b68ef4a24450a12d4f85d1835779dfef86;bpv=1;bpt=1;l=284?gsn=Accept&gs=kythe%3A%2F%2Fgo.googlesource.com%2Fgo%3Flang%3Dgo%3Fpath%3Dnet%23method%2520TCPListener.Accept)
      - [net/unixsock.go](https://cs.opensource.google/go/go/+/refs/tags/go1.18.3:src/net/unixsock.go;drc=386245b68ef4a24450a12d4f85d1835779dfef86;bpv=1;bpt=1;l=256?gsn=Accept&gs=kythe%3A%2F%2Fgo.googlesource.com%2Fgo%3Flang%3Dgo%3Fpath%3Dnet%23method%2520UnixListener.Accept)
      - `accept()`: [net/fd_unix.go](https://cs.opensource.google/go/go/+/refs/tags/go1.18.3:src/net/fd_unix.go;drc=386245b68ef4a24450a12d4f85d1835779dfef86;bpv=1;bpt=1;l=171?gsn=accept&gs=kythe%3A%2F%2Fgo.googlesource.com%2Fgo%3Flang%3Dgo%3Fpath%3Dnet%23method%2520netFD.accept)
    - `Close() error`
    - `Addr() Addr`
  - [type error](https://pkg.go.dev/builtin#error)


#### examples

##### Listen

- [listen_test.go](./listen_test.go)

```bash
go test listen_test.go -v
```

#### Dial: Accept and Receive

- [dial/dial.go](dial/dial.go)
- [dial/dial_test.go](dial/dial_test.go)

1. Start
   1. Listen `127.0.0.1:0`
   2. Error check
   3. Bound to `127.0.0.1:56758`
2. Wait a connection: blocking...
3. Dial to `127.0.0.1:56758` and Close
4. Receive `EOF`
5. Close `Listener`


### net.OpError

- [type OpError](https://pkg.go.dev/net#OpError)
  - [func Timeout](https://cs.opensource.google/go/go/+/refs/tags/go1.18.3:src/net/net.go;bpv=1;bpt=1;l=506?gsn=Timeout&gs=kythe%3A%2F%2Fgo.googlesource.com%2Fgo%3Flang%3Dgo%3Fpath%3Dnet%23method%2520OpError.Timeout)
  - [func Temporary](https://cs.opensource.google/go/go/+/refs/tags/go1.18.3:src/net/net.go;bpv=1;bpt=1;l=519?gsn=Temporary&gs=kythe%3A%2F%2Fgo.googlesource.com%2Fgo%3Flang%3Dgo%3Fpath%3Dnet%23method%2520OpError.Temporary)

```go
if nErr, ok := err.(net.Error); ok && !nErr.Temporary() {
  return err
}
```

### net.DialTimeout

- [net.DialTimeout](https://pkg.go.dev/net#DialTimeout)
  - [net/dial.go](https://cs.opensource.google/go/go/+/refs/tags/go1.18.3:src/net/dial.go;bpv=1;bpt=1;l=332?gsn=DialTimeout&gs=kythe%3A%2F%2Fgo.googlesource.com%2Fgo%3Flang%3Dgo%3Fpath%3Dnet%23func%2520DialTimeout)
- [net.Dial](https://pkg.go.dev/net#Dial)

```go
Dialer{Timeout: timeout}.Dial(network, address)
```

#### example: timeout

- [dial/dial_timeout.go](dial/dial_timeout.go)
- [dial/dial_timeout_test.go](dial/dial_timeout_test.go)

### net.DialContext

- [net.DialContext](https://pkg.go.dev/net#Dialer.DialContext)
  - [func WithDeadline](https://cs.opensource.google/go/go/+/refs/tags/go1.18.3:src/context/context.go;drc=386245b68ef4a24450a12d4f85d1835779dfef86;bpv=1;bpt=1;l=434?gsn=WithDeadline&gs=kythe%3A%2F%2Fgo.googlesource.com%2Fgo%3Flang%3Dgo%3Fpath%3Dcontext%23func%2520WithDeadline)
  - [func WithCancel](https://cs.opensource.google/go/go/+/refs/tags/go1.18.3:src/context/context.go;drc=386245b68ef4a24450a12d4f85d1835779dfef86;bpv=1;bpt=1;l=232?gsn=WithCancel&gs=kythe%3A%2F%2Fgo.googlesource.com%2Fgo%3Flang%3Dgo%3Fpath%3Dcontext%23func%2520WithCancel)

#### example: dial context

- [dial/dial_context_test.go](dial/dial_context_test.go)
- [dial/dial_cancel_test.go](dial/dial_cancel_test.go)
- [dial/dial_fanout_test.go](dial/dial_fanout_test.go)

### net.Conn

- [net.Conn](https://pkg.go.dev/net#Conn)
  - `SetDeadline(t time.Time) error`
  - `SetReadDeadline(t time.Time) error`
  - `SetWriteDeadline(t time.Time) error`

---

## Time

- [type Time](https://pkg.go.dev/time#Time)
  - [func NewTimer](https://pkg.go.dev/time#NewTimer)
  - [func Stop](https://pkg.go.dev/time#Timer.Stop)
    - Stop does not close the channel

```go
// if the call stops the timer
if t.Stop() {
}

// if the timer has already expired or been stopped
if !t.Stop() {
	<-t.C
}
```

#### example: timer

- [timer.go](./timer.go)

---

## Ping

- [ping.go](./ping/ping.go)
- [ping_test.go](./ping/ping_test.go)
- [ping_heartbeat_test.go](./ping/ping_heartbeat_test.go)

#### example: ping.go

1. Set `interval = 30s`
2. Loop: Write "ping" at `w`

#### example: ping_test.go

1. Init: r/w pipe, interval, timer channels
2. Loop: receivePing
   1. Reset Timer
   2. Read "Ping"
3. Terminate the context

#### example: ping_heartbeat_test.go

1. Start Listener
2. Start Pinger
3. Set deadline: 5s (4 pings + EOF)
4. CLIENT: Read 4 Pings + Send 1 Pong
5. SERVER: delay the deadline
6. CLIENT: Read 4 Pings
7. SERVER: Disconnect
8. CLIENT: Done
