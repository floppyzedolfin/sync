package watcher

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	pb "github.com/floppyzedolfin/sync/replica/replica"
	"github.com/howeyc/fsnotify"
)

type Watcher struct {
	notifier      *fsnotify.Watcher
	rootDir string
	replicaClient pb.ReplicaClient
}

func New(replicaClient pb.ReplicaClient) (*Watcher, error) {
	notifier, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("unable to instantiate fswatcher")
	}
	w := new(Watcher)
	w.notifier = notifier
	w.replicaClient = replicaClient
	return w, nil
}

func (w *Watcher) Watch(dirname string) error {
	defer w.notifier.Close()
	ctx, cancelFunc := context.WithCancel(context.Background())

	w.rootDir = dirname

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
		notificationFrequency = 5*time.Second
		cacheLimit            = 1
	)
	ticker := time.NewTicker(notificationFrequency)
	cache := make([]*fsnotify.FileEvent, 0, cacheLimit)
	for {
		select {
		case ev := <-w.notifier.Event:
			fmt.Printf("notif on : %s\n", ev.Name)
			err := w.checkDir(ctx, ev)
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
	fmt.Printf("sending cache of %d\n", len(cache))
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
			err = w.patchFile(ctx, filePath)
		} else {
			err = w.delete(ctx, filePath)
		}
		if err != nil {
			return fmt.Errorf("error while processing cached file: %w", err)
		}
	}
	return nil
}

func (w *Watcher) patchFile(ctx context.Context, filePath string) error {
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("unable to read contents of file %s: %w", filePath, err)
	}
	rights, err := getRights(filePath)
	if err != nil {
		return fmt.Errorf("unable to get the rights of file %s: %w", filePath, err)
	}
	patchRequest := pb.PatchFileRequest{
		FullPath:     strings.TrimPrefix(filePath, w.rootDir),
		FullContents: string(contents),
		Rights:       rights,
	}
	fmt.Printf("sending patch request %#v\n", patchRequest)
	res, err := w.replicaClient.PatchFile(ctx, &patchRequest)
	fmt.Printf("res: %#v\n", res)
	if err != nil {
		return fmt.Errorf("cmd error when patching %s: %w", filePath, err)
	}
	return nil
}

func (w *Watcher) createDir(ctx context.Context, dirPath string) error {
	rights, err := getRights(dirPath)
	if err != nil {
		return fmt.Errorf("unable to get the rights of dir %s: %w", dirPath, err)
	}
	createDirRequest := pb.CreateDirRequest{
		FullPath: strings.TrimPrefix(dirPath, w.rootDir),
		Rights:   rights,
	}
	_, err = w.replicaClient.CreateDir(ctx, &createDirRequest)
	if err !=  nil {
		return fmt.Errorf("server error while creating dir %s: %w", dirPath, err)
	}
	return nil
}

func (w *Watcher) delete(ctx context.Context, filePath string) error {
	deleteRequest := pb.DeleteRequest{FullPath: strings.TrimPrefix(filePath, w.rootDir)}
	_, err := w.replicaClient.Delete(ctx, &deleteRequest)
	if err != nil {
		return fmt.Errorf("cmd error while deleting %s on remote: %w", filePath, err)
	}
	return nil
}

// checkDir makes sure we properly plug on notifications for eventual subdirs
func (w *Watcher) checkDir(ctx context.Context, ev *fsnotify.FileEvent) error {
	path := ev.Name
	switch {
	case ev.IsDelete() || ev.IsRename():
		// ignore the error - we can't remember whether is was a dir or a file
		_ = w.notifier.RemoveWatch(path)
	case ev.IsCreate():
		// if it's a dir that's been created, we need to add subdirs - fsnotify isn't recursive, this is how we circumvent it
		err := w.walk(ctx, path)
		if err != nil {
			return fmt.Errorf("error while checking dir for new event: %w", err)
		}
	}
	return nil
}

func (w *Watcher) walk(ctx context.Context, dirPath string) error {
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
			w.createDir(ctx, path)
			return nil
		})
	if err != nil {
		return fmt.Errorf("error while walking path %s: %w", dirPath, err)
	}
	return nil
}

func getRights(path string) (uint32, error) {
	stats, err := os.Lstat(path)
	if err != nil {
		return 0, fmt.Errorf("can't lstat %s: %w", path, err)
	}
	return uint32(stats.Mode().Perm()), nil
}
