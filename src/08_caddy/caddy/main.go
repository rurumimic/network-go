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
