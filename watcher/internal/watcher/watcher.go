package watcher

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	pb "github.com/floppyzedolfin/sync/replica/replica"
	"github.com/howeyc/fsnotify"
)

// Watcher is the structure that can call the server to notify changes in the watched directory
type Watcher struct {
	notifier      *fsnotify.Watcher
	rootDir       string
	replicaClient server
	cancelFunc    context.CancelFunc
}

// New returns an operational watcher
func New(replicaClient server, cancelFunc context.CancelFunc) (*Watcher, error) {
	notifier, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("unable to instantiate fswatcher")
	}
	w := new(Watcher)
	w.notifier = notifier
	w.replicaClient = replicaClient
	w.cancelFunc = cancelFunc
	return w, nil
}

// Watch starts the infinite loop on a directory
func (w *Watcher) Watch(ctx context.Context, dirname string) error {
	defer w.notifier.Close()

	w.rootDir = strings.TrimRight(dirname, "/") + "/"

	// start the polling routing
	go w.poll(ctx)

	err := w.notifier.Watch(dirname)
	if err != nil {
		return fmt.Errorf("error while starting watcher on %s: %w", dirname, err)
	}

	// wait for a termination signal
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled")
	}
}

// poll listens to system notifications about changes on files / directories / ...
func (w *Watcher) poll(ctx context.Context) {
	for {
		select {
		case ev := <-w.notifier.Event:
			err := w.sendEvent(ctx, ev)
			if err != nil {
				fmt.Printf("error while sending event for %s to server: %s", ev.Name, err.Error())
				w.cancelFunc()
			}
		case e := <-w.notifier.Error:
			// something went wrong, stop all machines
			fmt.Printf("error in notification system: %s; aborting...", e.Error())
			w.cancelFunc()
		case <-ctx.Done():
			return
		}
	}
}

// sendEvent sends an event (if appropriate) to the server
func (w *Watcher) sendEvent(ctx context.Context, ev *fsnotify.FileEvent) error {
	if ev == nil {
		return nil
	}
	var err error
	update := ev.IsCreate() || ev.IsModify() || ev.IsAttrib()
	if update {
		err = w.createOrUpdate(ctx, ev.Name)
	} else {
		err = w.delete(ctx, ev.Name)
	}
	if err != nil {
		return fmt.Errorf("error while processing cached file: %w", err)
	}
	return nil
}

func (w *Watcher) createOrUpdate(ctx context.Context, path string) error {
	s, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("unable to get stat of %s: %w", path, err)
	}

	switch {
	case s.IsDir():
		err = w.walk(ctx, path)
	default:
		err = w.patchFile(ctx, path)
	}
	if err != nil {
		return fmt.Errorf("unable to update %s: %w", path, err)
	}
	return nil
}

// patchFile performs the creation/update of a file on the remote server
func (w *Watcher) patchFile(ctx context.Context, filePath string) error {
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("unable to read contents of file %s: %w", filePath, err)
	}
	patchRequest := pb.PatchFileRequest{
		FullPath:     strings.TrimPrefix(filePath, w.rootDir),
		FullContents: string(contents),
	}
	_, err = w.replicaClient.PatchFile(ctx, &patchRequest)
	if err != nil {
		return fmt.Errorf("cmd error when patching %s: %w", filePath, err)
	}
	return nil
}

// createDir creates a directory (and its subdirectories) on the remove file server
func (w *Watcher) createDir(ctx context.Context, dirPath string) error {
	fmt.Println("creating dir " + dirPath)
	createDirRequest := pb.CreateDirRequest{
		FullPath: strings.TrimPrefix(dirPath, w.rootDir),
	}
	_, err := w.replicaClient.CreateDir(ctx, &createDirRequest)
	if err != nil {
		return fmt.Errorf("server error while creating dir %s: %w", dirPath, err)
	}

	return nil
}

// delete removes an entity (file, dir, ...) from the server
func (w *Watcher) delete(ctx context.Context, filePath string) error {
	// ignore the error - we can't remember whether is was a dir or a file
	_ = w.notifier.RemoveWatch(filePath)

	deleteRequest := pb.DeleteRequest{FullPath: strings.TrimPrefix(filePath, w.rootDir)}
	_, err := w.replicaClient.Delete(ctx, &deleteRequest)
	if err != nil {
		return fmt.Errorf("server error while deleting %s on remote: %w", filePath, err)
	}
	return nil
}

// walk explores the contents of a directory for interesting things such as subdirectories or unlisted files
func (w *Watcher) walk(ctx context.Context, dirPath string) error {
	// check for subdirectories
	err := filepath.Walk(dirPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("unable to access file %s: %w", path, err)
			}
			if !info.Mode().IsDir() {
				// it's a file
				// this scenario can happen if the user creates a file inside a dir that isn't under scrupulous examination
				err = w.patchFile(ctx, path)
				if err != nil {
					return fmt.Errorf("error while adding file %s in the new directory %s: %w", path, dirPath, err)
				}
				return nil
			}
			// let's add this entry to the list of eligible parents
			err = w.notifier.Watch(path)
			if err != nil {
				// this could happen if we reached the inotify limit
				return fmt.Errorf("unable to add dir %s for watch: %w", dirPath, err)
			}
			err = w.createDir(ctx, path)
			if err != nil {
				return fmt.Errorf("unable to create remote dir at %s: %w", path, err)
			}
			return nil
		})
	if err != nil {
		return fmt.Errorf("error while walking path %s: %w", dirPath, err)
	}
	return nil
}
