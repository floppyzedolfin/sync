package main

import (
	"flag"
	"log"
	"path"

	"github.com/rjeczalik/notify"
)

func main() {
	var watchedDir string
	flag.StringVar(&watchedDir, "watchedDir", "example/playground", "path from which changes will be propagated")

	flag.Parse()

	c := make(chan notify.EventInfo, 1)
	if err := notify.Watch(path.Join(watchedDir, "..."), c, notify.All); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)
	for {
		select {
		case ei := <-c:
			log.Println("event: ", ei)
		}
	}
}
