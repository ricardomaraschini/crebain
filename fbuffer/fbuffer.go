package fbuffer

import "sync"

// FileFilterFn returns true if a path passes its internal filter logic.
type FileFilterFn func(path string) bool

// New returns a new file buffer populated with provided filters.
func New(filters ...FileFilterFn) *FBuffer {
	// remove possible nil FileFilterFn
	prepFilters := make([]FileFilterFn, 0, len(filters))
	for _, fn := range filters {
		if fn == nil {
			continue
		}
		prepFilters = append(prepFilters, fn)
	}
	return &FBuffer{
		db:      make(map[string]bool),
		filters: prepFilters,
	}
}

// FBuffer keeps all changed files since last interaction.
type FBuffer struct {
	sync.Mutex
	filters []FileFilterFn
	db      map[string]bool
}

// Push pushes a new file on the buffer. Some systems generate multiple
// events for every single file change, this buffer keeps an unified
// list where files are presented only once.
func (f *FBuffer) Push(filePath string) {
	if !f.accept(filePath) {
		return
	}

	f.Lock()
	f.db[filePath] = true
	f.Unlock()
}

// Drain returns a list of all files within the buffer, cleaning it up
// afterwards.
func (f *FBuffer) Drain() []string {
	f.Lock()
	files := make([]string, 0, len(f.db))
	for f := range f.db {
		files = append(files, f)
	}
	f.db = make(map[string]bool)
	f.Unlock()
	return files
}

// accept returns true if filePath is accepted by all filters.
func (f *FBuffer) accept(filePath string) bool {
	if len(f.filters) == 0 {
		return true
	}
	for _, fn := range f.filters {
		if !fn(filePath) {
			return false
		}
	}
	return true
}
