package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// NewWatcher returns a Watcher that monitors file changes on path,
// subdirectories are also monitored for changes as they got created.
func NewWatcher(path string, db *FileDB) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	watcher := &Watcher{fsw, db}
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

	go watcher.loop()
	return watcher, nil
}

// Watcher monitors changes on the filesystem.
type Watcher struct {
	*fsnotify.Watcher
	db *FileDB
}

// hookDir enables watcher on path, it complies with filepath.WalkFunc
// definition. If provided path does not point to a directory it
// simply ignores it.
func (w *Watcher) hookDir(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return nil
	}

	if err := w.Add(path); err != nil {
		return err
	}
	return nil
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
		// a directory, in case of file we add it to our internal db
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

		// only adds to our database if it is a file.
		if !finfo.IsDir() {
			w.db.Push(event.Name)
		}
		return
	}

	// we only care about write changes from this point on.
	if event.Op&fsnotify.Write != fsnotify.Write {
		return
	}

	w.db.Push(event.Name)
}
