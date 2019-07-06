package main

import "sync"

// FileFilterFn returns true if a path passes its internal filter logic.
type FileFilterFn func(path string) bool

// NewFileDB returns a new file database populated with provided filters.
func NewFileDB(filters ...FileFilterFn) *FileDB {
	return &FileDB{
		db:      make(map[string]bool),
		filters: filters,
	}
}

// FileDB keeps all changed files since last interaction.
type FileDB struct {
	sync.Mutex
	filters []FileFilterFn
	db      map[string]bool
}

// Push pushes a new file on the database. Some systems generate multiple
// events for every single file change, this database keeps an unified
// list where files are presented only once.
func (f *FileDB) Push(filePath string) {
	if !f.accept(filePath) {
		return
	}

	f.Lock()
	f.db[filePath] = true
	f.Unlock()
}

// Drain returns a list of all files within the DB, cleaning it up afterwards.
func (f *FileDB) Drain() []string {
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
func (f *FileDB) accept(filePath string) bool {
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
