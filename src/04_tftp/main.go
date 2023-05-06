package main

import (
	"flag"
	"log"
	"os"

	"tftp"
)

var (
	address = flag.String("a", "127.0.0.1:69", "listen address")
	payload = flag.String("p", "motorcycle.jpg", "file to serve to clients")
)

func main() {
	flag.Parse()

	p, err := os.ReadFile(*payload)
	if err != nil {
		log.Fatal(err)
	}

	s := tftp.Server{Payload: p}
	log.Fatal(s.ListenAndServe(*address))
}
