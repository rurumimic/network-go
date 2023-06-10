# HTTP Server

_back to [/README.md](/README.md)_

---

---

## Server

```bash
go get github.com/awoodbeck/gnp/ch09/handlers
```

- [server/server_test.go](server/server_test.go)

### TLS

```go
srv := &http.Server{
  Addr:        ":8443",
}

go func() {
  err := srv.ServeTLS(l, "cert.pem", "key.pem")
  if err != http.ErrServerClosed {
    t.Error(err)
  }
 }()
```

### Handlers

- [handlers/pitfall_test.go](handlers/pitfall_test.go)
- [handlers/methods.go](handlers/methods.go)
