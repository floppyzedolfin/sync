package watcher

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/floppyzedolfin/sync/reference/client"
	"github.com/floppyzedolfin/sync/reference/client/reference"
	"github.com/howeyc/fsnotify"
)

type Watcher struct {
	notifier  *fsnotify.Watcher
	refClient client.API
}

func New(refClient client.API) (*Watcher, error) {
	notifier, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("unable to instantiate fswatcher")
	}
	w := new(Watcher)
	w.notifier = notifier
	w.refClient = refClient
	return w, nil
}

func (w *Watcher) Watch(dirname string) error {
	defer w.notifier.Close()

	ctx, cancelFunc := context.WithCancel(context.Background())
	// start the polling routing
	go w.poll(ctx, cancelFunc)

	err := w.notifier.Watch(dirname)
	if err != nil {
		return fmt.Errorf("error while watching %s: %w", dirname, err)
	}

	// wait for a termination signal
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled")
	}

	return nil
}

func (w *Watcher) poll(ctx context.Context, cancelFunc context.CancelFunc) {
	const (
		notificationFrequency = time.Second
		cacheLimit            = 1
	)
	ticker := time.NewTicker(notificationFrequency)
	cache := make([]*fsnotify.FileEvent, 0, cacheLimit)
	for {
		select {
		case ev := <-w.notifier.Event:
			fmt.Printf("notif on : %s\n", ev.Name)
			err := w.checkDir(ev)
			if err != nil {
				fmt.Printf("error while checking path %s: %s; aborting...", ev.Name, err.Error())
				cancelFunc()
			}
			cache = append(cache, ev)
			if len(cache) == cacheLimit {
				ticker.Stop()
				w.sendCache(ctx, cache)
				cache = cache[:0]
				ticker = time.NewTicker(notificationFrequency)
			}
		case e := <-w.notifier.Error:
			// something went wrong, stop all machines
			fmt.Printf("error in notification system: %s; aborting...", e.Error())
			cancelFunc()
		case <-ticker.C:
			w.sendCache(ctx, cache)
			// reset cache
			cache = cache[:0]
		case <-ctx.Done():
			ticker.Stop()
		}
	}
}

func (w *Watcher) sendCache(ctx context.Context, cache []*fsnotify.FileEvent) error {
	fmt.Printf("sending cache of %n\n", len(cache))
	if len(cache) == 0 {
		return nil
	}
	// list all the files that have been updated and keep their final status (true to update, false to delete)
	updatedFiles := make(map[string]bool)
	for _, ev := range cache {
		if ev == nil {
			continue
		}
		toUpdate := ev.IsCreate() || ev.IsModify() || ev.IsAttrib()
		updatedFiles[ev.Name] = toUpdate
	}

	var err error
	for filePath, update := range updatedFiles {
		if update {
			err = w.patch(ctx, filePath)
		} else {
			err = w.delete(ctx, filePath)
		}
		if err != nil {
			return fmt.Errorf("error while processing cached file: %w", err)
		}
	}
	return nil
}

func (w *Watcher) patch(ctx context.Context, filePath string) error {
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("unable to read contents of file %s: %w", filePath, err)
	}
	patchRequest := reference.PatchRequest{
		FullPath:     filePath,
		FullContents: string(contents),
		Rights:       755,
	}
	fmt.Printf("sending patch request %#v\n", patchRequest)
	_, err = w.refClient.Patch(ctx, &patchRequest)
	if err != nil {
		return fmt.Errorf("server error when patching %s: %w", filePath, err)
	}
	return nil
}

func (w *Watcher) delete(ctx context.Context, filePath string) error {
	deleteRequest := reference.DeleteRequest{FullPath: filePath}
	_, err := w.refClient.Delete(ctx, &deleteRequest)
	if err != nil {
		return fmt.Errorf("server error while deleting %s on remote: %w", filePath, err)
	}
	return nil
}

// checkDir makes sure we properly plug on notifications for eventual subdirs
func (w *Watcher) checkDir(ev *fsnotify.FileEvent) error {
	path := ev.Name
	switch {
	case ev.IsDelete() || ev.IsRename():
		// ignore the error - we can't remember whether is was a dir or a file
		_ = w.notifier.RemoveWatch(path)
	case ev.IsCreate():
		// if it's a dir that's been created, we need to add subdirs - fsnotify isn't recursive, this is how we circumvent it
		err := w.walk(path)
		if err != nil {
			return fmt.Errorf("error while checking dir for new event: %w", err)
		}
	}
	return nil
}

func (w *Watcher) walk(dirPath string) error {
	// check for subdirs
	err := filepath.Walk(dirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("unable to access file %s: %w", path, err)
			}
			if !info.Mode().IsDir() {
				return nil
			}
			err = w.notifier.Watch(path)
			if err != nil {
				return fmt.Errorf("unable to add dir %s for watch: %w", dirPath, err)
			}
			return nil
		})
	if err != nil {
		return fmt.Errorf("error while walking path %s: %w", dirPath, err)
	}
	return nil
}
