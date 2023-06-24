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
mkcert -key-file key.pem -cert-file cert.pem 127.0.0.1
openssl x509 -text -noout -in cert.pem
```

##### OpenSSL

```bash
openssl req -x509 -newkey rsa:4096 \
-keyout key.pem \
-out cert.pem \
-sha256 -days 365 -nodes \
-subj "/C=XX/ST=StateName/L=CityName/O=CompanyName/OU=CompanySectionName/CN=CommonNameOrHostname" \
-addext "subjectAltName=IP:127.0.0.1"

chmod 440 key.pem
openssl x509 -text -noout -in cert.pem
```

```go
srv := &http.Server{
  Addr:        "127.0.0.1:8443",
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

`http_server/server/go.mod`:

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

func (h *Handlers) dbHandler() http.Handler {
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
