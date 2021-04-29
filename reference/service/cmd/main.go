package main

import "github.com/floppyzedolfin/sync/reference/service/internal"

const (
	port = 8405
)

func main() {
	server := internal.New()
	server.Listen(port)
}
