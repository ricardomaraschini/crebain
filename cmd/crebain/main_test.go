package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestHasTestFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "hasTestFile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	t.Run(".txt", func(t *testing.T) {
		f, err := ioutil.TempFile(dir, "myfavouritesingers.*.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		defer f.Close()

		if got := hasTestFile(f.Name()); got != false {
			t.Fatal("File is not a valid go file, but it's recognized as it is!")
		}
	})

	t.Run("_test.go", func(t *testing.T) {
		f, err := os.Create(dir + "/awesome_test.go")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(f.Name())
		defer f.Close()

		if got := hasTestFile(f.Name()); got != true {
			t.Fatal("File not recognized as testable")
		}
	})

	t.Run(".go", func(t *testing.T) {
		dotGo, err := os.Create(dir + "/could_be_better.go")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(dotGo.Name())
		defer dotGo.Close()

		t.Run("no tests", func(t *testing.T) {
			if got := hasTestFile(dotGo.Name()); got != false {
				t.Fatal("File recognized as belonging to a directory with packages")
			}
		})
		t.Run("test present", func(t *testing.T) {
			tf, err := os.Create(dir + "/could_be_better_test.go")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tf.Name())
			defer tf.Close()

			if got := hasTestFile(dotGo.Name()); got != true {
				t.Fatal("File not recognized as testable")
			}
		})
		t.Run("test present but different name", func(t *testing.T) {
			tf, err := os.Create(dir + "/plot_twist_test.go")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tf.Name())
			defer tf.Close()

			if got := hasTestFile(dotGo.Name()); got != true {
				t.Fatal("File not recognized as testable")
			}
		})
	})
}
