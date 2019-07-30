package main

import (
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/ricardomaraschini/crebain/fbuffer"
	"github.com/ricardomaraschini/crebain/match"
	"github.com/ricardomaraschini/crebain/trunner"
	"github.com/ricardomaraschini/crebain/tui"
	"github.com/ricardomaraschini/crebain/tui/basic"
	"github.com/ricardomaraschini/crebain/tui/fancy"
	"github.com/ricardomaraschini/crebain/watcher"
)

var userIf tui.UI

func main() {
	var exclude match.Multi
	var err error

	// Ignore hidden files and directories by default.
	exclude.Set("^\\.")

	dpath, err := os.Getwd()
	if err != nil {
		log.Fatal("Getwd:", err)
	}
	path := flag.String("path", dpath, "the path to be watched")
	xif := flag.Bool("tui", false, "enable text user interface")
	flag.Var(&exclude, "exclude", "regex rules for excluding paths from watching")
	flag.Parse()

	*path, err = filepath.Abs(*path)
	if err != nil {
		log.Fatal("Path:", err)
	}

	if err := os.Chdir(*path); err != nil {
		log.Fatal("Chdir:", err)
	}

	buf := fbuffer.New(hasTestFile)
	watcher, err := watcher.New(*path, exclude, buf)
	if err != nil {
		log.Fatal("NewWatcher:", err)
	}

	userIf = basic.New()
	if *xif {
		userIf, err = fancy.New()
		if err != nil {
			log.Fatal("tui.New():", err)
		}
	}

	go drainLoop(buf, time.Second)
	go readWatcherErrors(watcher)
	watcher.Loop()
	userIf.Start()
	watcher.Close()
}

// readWatcherErrors captures all errors and reports them on the interface.
func readWatcherErrors(w *watcher.Watcher) {
	for {
		err := <-w.Errors
		userIf.PushResult(&SysError{err})
	}
}

// drain loop iterates once every interval duration running tests on all
// changed modules. With some editors we may see multiple Write events
// almost at the same time, with this loop we consolidate what we have
// in memory thus running tests only once.
func drainLoop(db *fbuffer.FBuffer, interval time.Duration) {
	for range time.NewTicker(interval).C {
		modFiles := db.Drain()
		if len(modFiles) == 0 {
			continue
		}

		dedup := make(map[string]bool)
		modDirs := make([]string, 0)
		for _, fpath := range modFiles {
			dir := path.Dir(fpath)
			if _, ok := dedup[dir]; ok {
				continue
			}
			dedup[dir] = true
			modDirs = append(modDirs, dir)
		}

		testDirs(modDirs)
	}
}

// testDirs run go test on provided slice of directories.
func testDirs(dirs []string) {
	runner := trunner.New()
	for _, dir := range dirs {
		result, err := runner.Run(dir)
		if err != nil {
			userIf.PushResult(&SysError{err})
			return
		}
		userIf.PushResult(result)
	}
}

// hasTestFile is a filter that returns true if filePath belongs to a directory
// containing test files.
func hasTestFile(filePath string) bool {
	if filepath.Ext(filePath) != ".go" {
		return false
	}

	dir := path.Dir(filePath)
	testFiles, _ := filepath.Glob(dir + "/*_test.go")
	return len(testFiles) > 0
}
