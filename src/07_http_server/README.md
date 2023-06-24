# HTTP Server

_back to [/README.md](/README.md)_

---

---

## Server

```bash
go mod init http_server
```

### libs

```bash
go get github.com/awoodbeck/gnp/ch09/handlers
```

- [server/server_test.go](server/server_test.go)

### TLS

- [server/tls_test.go](server/tls_test.go)

#### Make Certs

##### Auto Gen

- [FiloSottile/mkcert](https://github.com/FiloSottile/mkcert)

```bash
# Brew
brew install mkcert
brew install nss # if you use Firefox

# Ports
sudo port selfupdate
sudo port install mkcert
sudo port install nss # if you use Firefox
```

```bash
mkcert -key-file key.pem -cert-file cert.pem "$HOST"
openssl x509 -text -noout -in cert.pem
```

```go
srv := &http.Server{
  Addr:        "localhost:8443",
}

go func() {
  err := srv.ServeTLS(l, "cert.pem", "key.pem")
  if err != http.ErrServerClosed {
    t.Error(err)
  }
 }()
```

```go
certPEM, err := os.ReadFile("./certs/cert.pem")
if err != nil {
  t.Fatal(err)
}

rootCAs, _ := x509.SystemCertPool()
if ok := rootCAs.AppendCertsFromPEM(certPEM); !ok {
  t.Fatal(err)
}
config := &tls.Config{
  InsecureSkipVerify: false,
  RootCAs:            rootCAs,
}
tr := &http.Transport{TLSClientConfig: config}
client := new(http.Client)
client.Transport = tr
path := fmt.Sprintf("https://%s/", srv.Addr)
```

### Handlers

- [handlers/default.go](handlers/default.go)
- [handlers/pitfall_test.go](handlers/pitfall_test.go)
- [handlers/methods.go](handlers/methods.go)

#### Replace Package

```bash
# in http_server/server
go mod edit -replace http_server/handlers=../handlers
go mod tidy
```

- [server/go.mod](server/go.mod):

```go
module server

go 1.20

require (
 github.com/awoodbeck/gnp v0.0.0-20230225045246-30fd6b8da810
 http_server/handlers v0.0.0-00010101000000-000000000000
)

replace http_server/handlers => ../handlers
```

`http_server/server/server_test.go`:

```go
package server

import (
  // "github.com/awoodbeck/gnp/ch09/handlers"
 "http_server/handlers"
)
```

#### Dependency Injection

```go
dbHandler := func(db *sql.DB) http.Handler {
  return http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
      err := db.Ping()
      // db jobs...
    },
  )
}
```

```go
type Handlers struct {
  db * sql.DB
  log *log.Logger
}

func (h *Handlers) Handler1() http.Handler {
  return http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
      err := h.db.Ping()
      if err != nil {
        h.log.Printf("db ping: %v", err)
      }
      // db jobs...
    },
  )
}

func (h *Handlers) Handler2() http.Handler {
  return http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
      // other jobs...
    },
  )
}
```

```go
h := &Handlers{
  db: db,
  log: log.New(os.Stderr, "handlers: ", log.Lshortfile),
}
http.Handle("/one", h.Handler1())
http.Handle("/two", h.Handler2())
```

### Middleware

- [middleware/timeout_test.go](middleware/timeout_test.go)
- [middleware/restrict_prefix.go](middleware/restrict_prefix.go)
- [middleware/restrict_prefix_test.go](middleware/restrict_prefix_test.go)
- [middleware/mux_test.go](middleware/mux_test.go)

```go
func Middleware(next http.Handler) http.Handler {
  return http.HandlerFunc(
    func(w. http.ResponseWriter, r *http.Request) {
      if r.Method == http.MethodTrace {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
      }
      w.Header().Set("X-Content-Type-Options", "nosniff")

      start := time.Now()
      next.ServeHTTP(w, r)
      log.Printf("Next handler duration %v", time.Now().Sub(start))
    },
  )
}
```

### HTTP/2 Server Push

- [server.go](server.go)

```bash
# in http_server
go mod edit -replace http_server/handlers=./handlers
go mod edit -replace http_server/middleware=./middleware
go mod tidy
```

- [go.mod](go.mod)

```go
module http_server

go 1.20

replace http_server/handlers => ./handlers

replace http_server/middleware => ./middleware

require (
 http_server/handlers v0.0.0-00010101000000-000000000000
 http_server/middleware v0.0.0-00010101000000-000000000000
)
```

```bash
go run server.go
# 2023/06/24 17:56:51 Serving files in "./files" over localhost:8080

go run server.go -listen $HOST:8443 -cert server/certs/cert.pem -key server/certs/key.pem
# 2023/06/24 17:58:16 Serving files in "./files" over localhost:8080
# 2023/06/24 17:58:16 TLS enabled
```

1. [http://localhost:8080/](http://localhost:8080/)
2. https://$HOST:8443/

---

## mkcert

```bash
# list java's trust store
keytool -list -v -keystore /Library/Java/JavaVirtualMachines/jdk-17.0.2.jdk/Contents/Home/lib/security/cacerts | less

# delete java's trust store
keytool -delete -alias "mkcert development ca 246652642225566010557230547171438441732" \
-keystore /Library/Java/JavaVirtualMachines/jdk-17.0.2.jdk/Contents/Home/lib/security/cacerts
```
