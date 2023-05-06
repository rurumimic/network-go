# TFTP

_back to [/README.md](/README.md)_

---

- [crypto](https://pkg.go.dev/crypto)
  - crypto/[sha512](https://pkg.go.dev/crypto/sha512)
    - crypto/sha512.[Sum512_256](https://pkg.go.dev/crypto/sha512#Sum512_256)

---

## TFTP Data Download

- [main.go](main.go)
- [tftp/types.go](tftp/types.go)
- [tftp/server.go](tftp/server.go)

### Server

#### go.mod

```bash
go mod edit -replace tftp=./tftp
go mod tidy
```

#### Build and Run Server

```bash
sudo go run main.go
# or
go build -o main main.go
sudo ./main
```

```bash
2023/05/06 17:00:42 Listening on 127.0.0.1:69
2023/05/06 17:00:55 [127.0.0.1:57581] requested file: motorcycle.jpg
2023/05/06 17:00:55 [127.0.0.1:57581] sent 276 blocks
```

### Client

```bash
cd client
```

```bash
tftp -e 127.0.0.1

tftp> verbose
Verbose mode on.

tftp> binary
mode set to octet

tftp> get motorcycle.jpg
getting from 127.0.0.1:motorcycle.jpg to motorcycle.jpg [octet]
Received 140847 bytes in 0.0 seconds [inf bits/sec]

tftp> quit
```

```bash
ls -al motorcycle.jpg
```

<img alt="motorcycle" src="motorcycle.jpg" width="300px">

---

## Checksum

- [checksum/checksum.go](checksum/checksum.go)

### Compare Image

```bash
go run checksum/checksum.go motorcycle.jpg client/motorcycle.jpg

76c7045a9d4e2103482fcc186b3e8c0b8f1bcb64d183dcde1c8f82481f07eb4c motorcycle.jpg
76c7045a9d4e2103482fcc186b3e8c0b8f1bcb64d183dcde1c8f82481f07eb4c client/motorcycle.jpg
```
