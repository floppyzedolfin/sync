package main

import (
	"flag"
	"log"

	"github.com/floppyzedolfin/sync/reference/client"
	"github.com/floppyzedolfin/sync/watcher/internal/watcher"
)

func main() {
	var watchedDir string
	flag.StringVar(&watchedDir, "watchedDir", "example/playground", "path from which changes will be propagated")
	flag.Parse()

	c := client.NewClient()
	w, err := watcher.New(c)
	if err != nil {
		log.Fatalf("error while instantiating the watcher: %s", err.Error())
	}
	w.Watch(watchedDir)
}
