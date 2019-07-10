package watcher

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type matcher interface {
	Match(value string) bool
}

type buffer interface {
	Push(path string)
}

// New returns a Watcher that monitors file changes on path,
// subdirectories are also monitored for changes as they got created.
func New(path string, exclude matcher, buf buffer) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watcher := &Watcher{
		Watcher:  fsw,
		buf:      buf,
		exclude:  exclude,
		rootPath: path,
	}
	finfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// hook the root directory.
	if watcher.hookDir(path, finfo, nil); err != nil {
		_ = fsw.Close()
		return nil, err
	}

	// recursively hook all sub directories.
	if err := filepath.Walk(path, watcher.hookDir); err != nil {
		_ = fsw.Close()
		return nil, err
	}

	return watcher, nil
}

func (w *Watcher) Loop() {
	go w.loop()
}

// Watcher monitors changes on the filesystem.
type Watcher struct {
	*fsnotify.Watcher
	buf      buffer
	exclude  matcher
	rootPath string
}

// hookDir enables watcher on path, it complies with filepath.WalkFunc
// definition. If provided path does not point to a directory it
// simply ignores it.
func (w *Watcher) hookDir(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !w.isWatchable(path, info) {
		return nil
	}

	return w.Add(path)
}

func (w *Watcher) isWatchable(path string, info os.FileInfo) bool {
	if w.isPathExcluded(path) {
		return false
	}
	return info.IsDir()
}

// isPathExcluded checks whether the path matches against the exclusion rules.
// Check is performed in relation of the root path.
func (w *Watcher) isPathExcluded(path string) bool {
	relative := strings.TrimPrefix(path, w.rootPath)
	relative = strings.TrimPrefix(relative, "/")
	return w.exclude.Match(relative)
}

// loop awaits for file write operations. Everytime a write happens on monitored
// path it pushes the monitored file towards its internal FileDB.
func (w *Watcher) loop() {
	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				log.Fatal("watcher channel closed.")
			}
			w.processEvent(event)
		case err, ok := <-w.Errors:
			if !ok {
				log.Fatal("watcher errors channel closed.")
			}
			log.Println("watcher error:", err)
		}
	}
}

// processEvent is called everytime we detect a change on the filesystem.
func (w *Watcher) processEvent(event fsnotify.Event) {
	if event.Op&fsnotify.Create == fsnotify.Create {
		// if something got created we need to check if it is a file or
		// a directory, in case of file we add it to our internal buffer
		// and if it is a directory we hook ourselves on it to capture
		// future events.
		finfo, err := os.Stat(event.Name)
		if err != nil {
			log.Println("Stat:", event.Name, err)
			return
		}

		// try to hook on this new file/directory. If it is not a dir
		// it will be a no-op anyways.
		if err := w.hookDir(event.Name, finfo, nil); err != nil {
			log.Println("hookDir:", event.Name, err)
			return
		}

		// we only push files to buffer, never directories.
		if !finfo.IsDir() {
			w.buf.Push(event.Name)
		}
		return
	}

	// Ignore events that are only chmod.
	// Write, rename and remove are acceptable.
	if event.Op == fsnotify.Chmod {
		log.Println("Chmod:", event.Name)
		return
	}

	w.buf.Push(event.Name)
}
