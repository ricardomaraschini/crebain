package fbuffer

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	buf := New()
	if len(buf.filters) != 0 {
		t.Fatalf("expected lenght of filters to be zero")
	}

	buf = New(func(f string) bool {
		return true
	})

	if buf.db == nil {
		t.Fatal("expected internal map not to be nil")
	}

	if len(buf.filters) != 1 {
		t.Fatal("expected lenght of filters to be zero")
	}
}

func TestPush(t *testing.T) {
	for _, tt := range []struct {
		name  string
		fns   []FileFilterFn
		exp   int
		paths []string
	}{
		{
			name: "no filter",
			fns:  []FileFilterFn{},
			paths: []string{
				"/file0",
				"/file2",
			},
			exp: 2,
		},
		{
			name: "nil filter",
			fns:  []FileFilterFn{nil},
			paths: []string{
				"/file0",
				"/file1",
				"/file2",
			},
			exp: 3,
		},
		{
			name: "deny all",
			fns: []FileFilterFn{
				func(string) bool {
					return false
				},
			},
			paths: []string{
				"/file0",
				"/file1",
				"/file2",
			},
		},
		{
			name: "accept only go suffix",
			fns: []FileFilterFn{
				func(p string) bool {
					return strings.HasSuffix(p, "go")
				},
			},
			paths: []string{
				"/file0",
				"/file/something_go",
				"/file1",
				"/file/something_else_go",
				"/file2",
			},
			exp: 2,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			buf := New(tt.fns...)
			for _, path := range tt.paths {
				buf.Push(path)
			}

			fs := buf.Drain()
			if len(fs) != tt.exp {
				t.Fatalf("expected db len %d,  db len %d instead", tt.exp, len(buf.db))
			}
		})
	}
	// using a filter that denies everything.
	buf := New(func(string) bool {
		return false
	})
	buf.Push("/test")
	if len(buf.db) != 0 {
		t.Fatal("expected len of db to be zero")
	}

	buf = New()
	buf.Push("/test")
	buf.Push("/test")
	buf.Push("/test")
	if len(buf.db) != 1 {
		t.Fatal("expected len of db to be 1")
	}
}
