package watcher

import (
	"log"
	"path"

	"github.com/rjeczalik/notify"
)

type Watcher struct {
	path string
	eventsChan chan notify.EventInfo
}

func New(path string) *Watcher {
	w := new(Watcher)

	// buffered channel to ensure we dont' miss events
	w.eventsChan = make(chan notify.EventInfo, 1)
	w.path = path

	return w
}

func (w *Watcher) Watch() {
	if err := notify.Watch(path.Join(w.path, "..."), w.eventsChan); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(w.eventsChan)
	w.handleChanges()
}

func (w *Watcher) handleChanges() {
	ei := <-w.eventsChan
	handleEvent(ei)
}

func handleEvent(ei notify.EventInfo) {
	log.Println("Got event:", ei)
}