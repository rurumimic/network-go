# TLS

_back to [/README.md](/README.md)_

---

---

## TLS Test

- [tls_client_test.go](tls_client_test.go)
- [tls_echo.go](tls_echo.go)

### Generate a self-signed certificate

```bash
go run $GOROOT/src/crypto/tls/generate_cert.go -host localhost -ecdsa-curve P256
```

### Make Certs

```bash
go run cert/generate.go -cert serverCert.pem -key serverKey.pem -host localhost

2023/07/08 16:24:32 wrote serverCert.pem
2023/07/08 16:24:32 wrote serverKey.pem

go run cert/generate.go -cert clientCert.pem -key clientKey.pem -host localhost

2023/07/08 16:25:15 wrote clientCert.pem
2023/07/08 16:25:15 wrote clientKey.pem
```
