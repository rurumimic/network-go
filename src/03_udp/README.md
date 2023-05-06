# UDP

_back to [/README.md](/README.md)_

---

- [net](https://pkg.go.dev/net)
  - net.[Addr](https://pkg.go.dev/net#Addr)
  - net.[ListenPacket](https://pkg.go.dev/net#ListenPacket)
  - net.[PacketConn](https://pkg.go.dev/net#PacketConn)
  - net.[Dial](https://pkg.go.dev/net#Dial)

---

## net

### net.ListenPacket & net.PacketConn & net.Dial

#### examples

##### echo

- [echo/echo.go](echo/echo.go)
- [echo/echo_test.go](echo/echo_test.go)
- [echo/listen_packet_test.go](echo/listen_packet_test.go)
- [echo/dial_test.go](echo/dial_test.go)

--

## Fragmentation

- MTU: [Maximum Transmission Unit](https://en.wikipedia.org/wiki/Maximum_transmission_unit)
  - Ethernet: 1500 bytes

### Ping

- ICMP: [Internet Control Message Protocol](https://en.wikipedia.org/wiki/Internet_Control_Message_Protocol)
  - header: 8-byte
- IPv4: [Internet Protocol](https://en.wikipedia.org/wiki/Internet_Protocol_version_4)
  - header: 20-byte

#### Success

```bash
# mac
ping -D -s 1472 1.1.1.1

PING 1.1.1.1 (1.1.1.1): 1472 data bytes
1480 bytes from 1.1.1.1: icmp_seq=0 ttl=55 time=5.097 ms
1480 bytes from 1.1.1.1: icmp_seq=1 ttl=55 time=4.678 ms
1480 bytes from 1.1.1.1: icmp_seq=2 ttl=55 time=5.391 ms
^C
--- 1.1.1.1 ping statistics ---
3 packets transmitted, 3 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 4.678/5.055/5.391/0.293 ms
```

```bash
# linux
ping -M do -s 1472 1.1.1.1
```

```bash
# windows
ping -f -l 1472 1.1.1.1
```

#### Failure

```bash
# mac
ping -D -s 1500 1.1.1.1

PING 1.1.1.1 (1.1.1.1): 1500 data bytes
ping: sendto: Message too long
ping: sendto: Message too long
Request timeout for icmp_seq 0
ping: sendto: Message too long
Request timeout for icmp_seq 1
^C
--- 1.1.1.1 ping statistics ---
3 packets transmitted, 0 packets received, 100.0% packet loss
```

```bash
# linux
ping -M do -s 1500 1.1.1.1
```

```bash
# windows
ping -f -l 1500 1.1.1.1
```
