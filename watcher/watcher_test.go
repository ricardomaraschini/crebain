package watcher

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/fsnotify/fsnotify"
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
			expected:  false,
		},
		{
			name:      "directory",
			path:      dir,
			exclusion: func(_ string) bool { return false },
			expected:  true,
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
				exclude: &dummyMatcher{st.exclusion},
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
		exclude: &dummyMatcher{func(relative string) bool {
			return relative == "orthanc"
		}},
		rootPath: "/mytowers",
	}
	expected := true

	if got := watcher.isPathExcluded(path); got != expected {
		t.Fatal("Unexpected result:", got)
	}
}

type dummyBuffer struct {
	element string
}

func (d *dummyBuffer) Push(path string) {
	d.element = path
}

func TestNewWithInvalidDir(t *testing.T) {
	matcher := &dummyMatcher{
		func(_ string) bool { return false },
	}
	buf := &dummyBuffer{}
	_, err := New(
		"/tmp/does-not-exist",
		matcher,
		buf,
	)

	if err == nil {
		t.Fatal("expected error, received nil instead")
	}

}

func TestProcessEvent(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "processEvent")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile, err := os.Create(tmpDir + "/confusion.go")
	if err != nil {
		t.Fatal(err)
	}
	defer tmpFile.Close()

	name := tmpFile.Name()
	buf := &dummyBuffer{}
	matcher := &dummyMatcher{
		func(_ string) bool { return false },
	}

	w, err := New(
		tmpDir,
		matcher,
		buf,
	)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Watcher.Close()

	t.Run("accepted", func(t *testing.T) {
		ops := []fsnotify.Op{
			fsnotify.Create,
			fsnotify.Create | fsnotify.Write,
			fsnotify.Write,
			fsnotify.Rename | fsnotify.Write,
			fsnotify.Rename,
			fsnotify.Chmod | fsnotify.Write,
		}
		for _, op := range ops {
			t.Run(op.String(), func(t *testing.T) {
				buf.element = ""

				e := fsnotify.Event{
					Name: name,
					Op:   op,
				}

				w.processEvent(e)
				if buf.element != name {
					t.Fatalf("Ignored %s event", op)
				}
			})
		}
	})

	t.Run("ignored", func(t *testing.T) {
		t.Run("chmod", func(t *testing.T) {
			buf.element = ""

			op := fsnotify.Chmod
			e := fsnotify.Event{
				Name: name,
				Op:   op,
			}

			w.processEvent(e)
			if buf.element == name {
				t.Fatalf("Accepted %s event", op)
			}
		})
		t.Run("dir", func(t *testing.T) {
			buf.element = ""
			newDir, err := ioutil.TempDir("", "processEvent2")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(newDir)

			op := fsnotify.Create
			e := fsnotify.Event{
				Name: newDir,
				Op:   op,
			}

			w.processEvent(e)
			if buf.element == newDir {
				t.Fatalf("Accepted dir %s", newDir)
			}
		})
		t.Run("excluded", func(t *testing.T) {
			buf.element = ""
			matcher.match = func(_ string) bool { return true }

			op := fsnotify.Create
			e := fsnotify.Event{
				Name: name,
				Op:   op,
			}

			w.processEvent(e)
			if buf.element == name {
				t.Fatalf("Accepted %s excluded", name)
			}
		})
	})
}
