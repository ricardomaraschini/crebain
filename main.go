package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/ricardomaraschini/crebain/fbuffer"
	"github.com/ricardomaraschini/crebain/match"
)

func main() {
	exclusionRules := match.Multi{}
	dpath, err := os.Getwd()
	if err != nil {
		log.Fatal("Getwd:", err)
	}
	path := flag.String("path", dpath, "the path to be watched")
	flag.Var(&exclusionRules, "e", "regex rules for excluding some path from watching")
	flag.Parse()

	if err := os.Chdir(*path); err != nil {
		log.Fatal("Chdir:", err)
	}

	buf := fbuffer.New(hasTestFile)
	watcher, err := NewWatcher(*path, exclusionRules, buf)
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
		modifiedFiles := db.Drain()
		if len(modifiedFiles) == 0 {
			continue
		}

		dedup := make(map[string]bool)
		modifiedDirs := make([]string, 0)
		for _, fpath := range modifiedFiles {
			dir := path.Dir(fpath)
			if _, ok := dedup[dir]; ok {
				continue
			}
			dedup[dir] = true
			modifiedDirs = append(modifiedDirs, dir)
		}

		test(modifiedDirs)
	}
}

func test(dirs []string) {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	for _, dir := range dirs {
		cmd := exec.Command("go", "test", "-v", "-cover", dir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}
