package main

import (
	"io/ioutil"
	"os"
	"testing"
)

type dummyMatcher struct {
	match func(path string) bool
}

func (dm *dummyMatcher) Match(path string) bool {
	return dm.match(path)
}

func TestIsWatchable(t *testing.T) {
	dir, err := ioutil.TempDir("", "crebain")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	tmpFile, err := ioutil.TempFile("", "crebain")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpFile.Name())

	subtests := []struct {
		name      string
		path      string
		exclusion func(path string) bool
		expected  bool
	}{
		{
			name:      "file",
			path:      tmpFile.Name(),
			exclusion: func(_ string) bool { return false },
			expected:  true,
		},
		{
			name:      "directory",
			path:      dir,
			exclusion: func(_ string) bool { return false },
			expected:  false,
		},
		{
			name:      "excluded file",
			path:      tmpFile.Name(),
			exclusion: func(_ string) bool { return true },
			expected:  false,
		},
	}

	for _, st := range subtests {
		t.Run(st.name, func(t *testing.T) {

			info, err := os.Stat(st.path)
			if err != nil {
				t.Fatal(err)
			}

			watcher := Watcher{
				exclusionRules: &dummyMatcher{st.exclusion},
			}
			if got := watcher.isWatchable(dir, info); got != st.expected {
				t.Fatal("Unexpected result:", got)
			}

		})
	}
}

func TestIsPathExcluded(t *testing.T) {
	path := "/mytowers/orthanc"
	watcher := Watcher{
		exclusionRules: &dummyMatcher{func(relative string) bool {
			return relative == "orthanc"
		}},
		rootPath: "/mytowers",
	}
	expected := true

	if got := watcher.isPathExcluded(path); got != expected {
		t.Fatal("Unexpected result:", got)
	}
}
