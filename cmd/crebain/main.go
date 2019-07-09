package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/ricardomaraschini/crebain/fbuffer"
	"github.com/ricardomaraschini/crebain/match"
	"github.com/ricardomaraschini/crebain/watcher"
)

func main() {
	var exclude match.Multi

	dpath, err := os.Getwd()
	if err != nil {
		log.Fatal("Getwd:", err)
	}
	path := flag.String("path", dpath, "the path to be watched")
	flag.Var(&exclude, "exclude", "regex rules for excluding paths from watching")
	flag.Parse()

	if err := os.Chdir(*path); err != nil {
		log.Fatal("Chdir:", err)
	}

	buf := fbuffer.New(hasTestFile)
	watcher, err := watcher.New(*path, exclude, buf)
	if err != nil {
		log.Fatal("NewWatcher:", err)
	}
	defer watcher.Close()
	drainLoop(buf, time.Second)
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

		test(modDirs)
	}
}

func test(dirs []string) {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	for _, dir := range dirs {
		cmd := exec.Command("go", "test", "-cover", dir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
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
