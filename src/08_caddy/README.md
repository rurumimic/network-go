# Caddy

_back to [/README.md](/README.md)_

---

---

## Caddy Server

- [caddyserver/caddy](https://github.com/caddyserver/caddy)

### Build from source

```bash
git clone https://github.com/caddyserver/caddy.git
cd caddy/cmd/caddy/
go build
```

### Start Caddy

```bash
./caddy start
```

```bash
2023/07/08 03:52:46.924 INFO admin admin endpoint started {"address": "localhost:2019", "enforce_origin": false, "origins": ["//localhost:2019", "//[::1]:2019", "//127.0.0.1:2019"]}
2023/07/08 03:52:46.924 INFO serving initial configuration
Successfully started Caddy (pid=33005) - Caddy is running in the background
```

#### Update a endpoint

```bash
curl localhost:2019/load \
-X POST -H "Content-Type: application/json" \
-d '
{
  "apps": {
    "http": {
      "servers": {
        "hello": {
          "listen": ["localhost:2020"],
          "routes": [{
            "handle": [{
              "handler": "static_response",
              "body": "Hello, World!"
            }]
          }]
        }
      }
    }
  }
}
'
```

#### Listen ports

```bash
lsof -Pi :2019-2025
COMMAND     PID USER    FD   TYPE             DEVICE SIZE/OFF NODE NAME
caddy     33005 keanu   12u  IPv4 0xfc12deda094844f1      0t0  TCP localhost:2019 (LISTEN)
caddy     33005 keanu   13u  IPv4 0xfc12deda08780bf1      0t0  TCP localhost:2020 (LISTEN)
```

#### Configuration List

```bash
curl localhost:2019/config/

{"apps":{"http":{"servers":{"hello":{"listen":["localhost:2020"],"routes":[{"handle":[{"body":"Hello, World!","handler":"static_response"}]}]}}}}}
```

#### Configuration traversal

```bash
curl localhost:2019/config/apps/http/servers/hello/listen

["localhost:2020"]
```

#### Get Response

```bash
curl -X GET localhost:2020

Hello, World!
```

#### Add a endpoint

```bash
curl localhost:2019/config/apps/http/servers/hello/listen \
-X POST -H "Content-Type: application/json" -d '"localhost:2021"'
```

```bash
lsof -Pi :2019-2025

COMMAND   PID USER    FD   TYPE             DEVICE SIZE/OFF NODE NAME
caddy   33005 keanu   12u  IPv4 0xfc12deda0948c4f1      0t0  TCP localhost:2020 (LISTEN)
caddy   33005 keanu   14u  IPv4 0xfc12deda082564f1      0t0  TCP localhost:2019 (LISTEN)
caddy   33005 keanu   15u  IPv4 0xfc12deda0948a391      0t0  TCP localhost:2021 (LISTEN)
```

```bash
curl localhost:2019/config/apps/http/servers/hello/listen \
-X PATCH -H "Content-Type: application/json" -d '["localhost:2020-2025"]'
```

```bash
lsof -Pi :2019-2025

COMMAND   PID USER    FD   TYPE             DEVICE SIZE/OFF NODE NAME
caddy   33005 keanu   13u  IPv4 0xfc12deda0949b9d1      0t0  TCP localhost:2019 (LISTEN)
caddy   33005 keanu   14u  IPv4 0xfc12deda08790711      0t0  TCP localhost:2020 (LISTEN)
caddy   33005 keanu   16u  IPv4 0xfc12deda09491871      0t0  TCP localhost:2021 (LISTEN)
caddy   33005 keanu   17u  IPv4 0xfc12deda087670d1      0t0  TCP localhost:2022 (LISTEN)
caddy   33005 keanu   18u  IPv4 0xfc12deda08783871      0t0  TCP localhost:2023 (LISTEN)
caddy   33005 keanu   19u  IPv4 0xfc12deda08792871      0t0  TCP localhost:2024 (LISTEN)
caddy   33005 keanu   20u  IPv4 0xfc12deda08780bf1      0t0  TCP localhost:2025 (LISTEN)
```

```bash
curl localhost:2019/config/apps/http/servers/test \
-X POST -H "Content-Type: application/json" \
-d '
{
  "listen": ["localhost:2030"],
  "routes": [{
    "handle": [{
      "handler": "static_response",
      "body": "Welcome to my temporary test server."
    }]
  }]
}
'
```

```bash
curl localhost:2030

Welcome to my temporary test server.
```

```bash
curl localhost:2019/config/apps/http/servers/test -X DELETE
curl localhost:2030

curl: (7) Failed to connect to localhost port 2030 after 8 ms: Couldn't connect to server
```

```bash
./caddy stop
```

#### Save configurations

```bash
vi caddy.json
```

```json
{
  "apps": {
    "http": {
      "servers": {
        "hello": {
          "listen": ["localhost:2020-2025"],
          "routes": [{
            "handle": [{
              "handler": "static_response",
              "body": "Hello, World!"
            }]
          }]
        },
        "test": {
          "listen": ["localhost:2030"],
          "routes": [{
            "handle": [{
              "handler": "static_response",
              "body": "Welcome to my temporary test server."
            }]
          }]
        }
      }
    }
  }
}
```

```bash
./caddy start --config caddy.json
```

```bash
lsof -Pi :2019-2035

COMMAND   PID USER    FD   TYPE             DEVICE SIZE/OFF NODE NAME
caddy   41721 keanu   10u  IPv4 0xfc12deda09491871      0t0  TCP localhost:2019 (LISTEN)
caddy   41721 keanu   11u  IPv4 0xfc12deda094939d1      0t0  TCP localhost:2030 (LISTEN)
caddy   41721 keanu   12u  IPv4 0xfc12deda09492eb1      0t0  TCP localhost:2020 (LISTEN)
caddy   41721 keanu   13u  IPv4 0xfc12deda08782231      0t0  TCP localhost:2021 (LISTEN)
caddy   41721 keanu   14u  IPv4 0xfc12deda08783871      0t0  TCP localhost:2022 (LISTEN)
caddy   41721 keanu   15u  IPv4 0xfc12deda09489871      0t0  TCP localhost:2023 (LISTEN)
caddy   41721 keanu   16u  IPv4 0xfc12deda09498d51      0t0  TCP localhost:2024 (LISTEN)
caddy   41721 keanu   17u  IPv4 0xfc12deda08792871      0t0  TCP localhost:2025 (LISTEN)
```

```bash
curl localhost:2019/config/

{"apps":{"http":{"servers":{"hello":{"listen":["localhost:2020-2025"],"routes":[{"handle":[{"body":"Hello, World!","handler":"static_response"}]}]},"test":{"listen":["localhost:2030"],"routes":[{"handle":[{"body":"Welcome to my temporary test server.","handler":"static_response"}]}]}}}}}
```

```bash
curl localhost:2020 # Hello, World!
curl localhost:2030 # Welcome to my temporary test server.
```

```bash
./caddy stop
```

---

## Caddyfile adapter

ref:

- [awoodbeck/caddy-toml-adapter](https://github.com/awoodbeck/caddy-toml-adapter)
- [awoodbeck/caddy-restrict-prefix](https://github.com/awoodbeck/caddy-restrict-prefix)

code:

- [toml-adapter/toml.go](toml-adapter/toml.go)
- [restrict-prefix/toml.go](restrict-prefix/restrict_prefix.go)

```bash
cd toml-adapter
go mod init toml-adapter
go mod tidy

cd restrict-prefix
go mod init restrict-prefix
go mod tidy
```

### Build Caddy

```bash
mkdir caddy
cd caddy
go mod init caddy
```

```bash
vi go.mod
```

```go
module caddy

go 1.20

replace toml-adapter => ../toml-adapter
replace restrict-prefix => ../restrict-prefix
```

```bash
vi main.go
```

```bash
package main

import (
 caddycmd "github.com/caddyserver/caddy/v2/cmd"

 // plug in Caddy modules here
 _ "github.com/caddyserver/caddy/v2/modules/standard"

 _ "restrict-prefix"
 _ "toml-adapter"
)

func main() {
 caddycmd.Main()
}
```

```bash
go mod tidy
go build
```

```bash
./caddy list-modules | grep "toml\|restrict_prefix"

caddy.adapters.toml
http.handlers.restrict_prefix
```

---

## With Backend Service

- [caddy/backend/main.go](caddy/backend/main.go)
- [caddy/caddy.toml](caddy/caddy.toml)

```bash
./caddy start --config caddy.toml --adapter toml
```

```bash
curl localhost:2019/config/

{"apps":{"http":{"servers":{"test_server":{"listen":["localhost:2020"],"routes":[{"handle":[{"handler":"reverse_proxy","upstreams":[{"dial":"localhost:8080"}]}],"match":[{"path":["/backend","/backend/*"]}]},{"handle":[{"handler":"restrict_prefix","prefix":"."},{"handler":"file_server","index_names":["index.html"],"root":"./files"}]}]}}}}}
```

```bash
go run main.go

Listening on localhost:8080 ...
```

- Open [localhost:2020](http://localhost:2020)
- Open [localhost:2020/backend](http://localhost:2020/backend)
- Open [localhost:2020/files](http://localhost:2020/backend)

Turn on logging:

```bash
curl localhost:2019/config/logging -X POST -H "Content-Type: application/json" -d '{"logs": {"": {"level": "DEBUG"}}}';
curl localhost:2019/config/logging -X POST -H "Content-Type: application/json" -d '{"logs": {"access": {"level": "debug"}}}';
```

```bash
curl -i http://localhost:2020/.dir/secret

HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
Server: Caddy
X-Content-Type-Options: nosniff
Date: Sat, 08 Jul 2023 05:39:30 GMT
Content-Length: 10

Not Found

2023/07/08 05:39:19.250 DEBUG http.handlers.restrict_prefix restricted prefix: ".dir" in /.dir/secret
```

```bash
./caddy stop
```
