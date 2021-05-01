package main

import (
	"github.com/floppyzedolfin/sync/reference/client"
	"github.com/floppyzedolfin/sync/reference/service/internal"
)


func main() {
	// start and run the server
	server := internal.New()
	server.Listen(client.PortNumber)
}
