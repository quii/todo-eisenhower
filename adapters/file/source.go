// Package file provides file-based adapters for reading and writing todo.txt files.
package file

import (
	"os"
	"path/filepath"

	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/domain/todotxt"
)

// Repository is a file-based implementation of TodoRepository
// It handles all the file mechanics and marshaling/unmarshaling
type Repository struct {
	path string
}

// NewRepository creates a new file-based todo repository
func NewRepository(path string) *Repository {
	return &Repository{path: path}
}

// LoadAll reads todos from the file
func (r *Repository) LoadAll() ([]todo.Todo, error) {
	// Create file if it doesn't exist
	if _, err := os.Stat(r.path); os.IsNotExist(err) {
		//nolint:gosec // G306: todo.txt files are intentionally world-readable (0o644 per todo.txt spec)
		if err := os.WriteFile(r.path, []byte(""), 0o644); err != nil {
			return nil, err
		}
	}

	f, err := os.Open(r.path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()

	return todotxt.Unmarshal(f)
}

// SaveAll writes todos to the file (full rewrite)
func (r *Repository) SaveAll(todos []todo.Todo) error {
	//nolint:gosec // G302: todo.txt files are intentionally world-readable (0o644 per todo.txt spec)
	f, err := os.Create(r.path) // Truncates automatically
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	return todotxt.Marshal(f, todos)
}

// AppendToArchive appends a todo to the archive file (done.txt)
// The archive file is in the same directory as todo.txt, named done.txt
func (r *Repository) AppendToArchive(t todo.Todo) error {
	// Determine archive file path (done.txt in same directory as todo.txt)
	archivePath := r.archivePath()

	// Open file in append mode, create if doesn't exist
	//nolint:gosec // G302,G304: done.txt files are intentionally world-readable (0o644 per todo.txt spec)
	f, err := os.OpenFile(archivePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	// Write the todo as a single line
	return todotxt.Marshal(f, []todo.Todo{t})
}

// archivePath returns the path to the archive file (done.txt)
// Following todo.txt convention: always use done.txt in the same directory
func (r *Repository) archivePath() string {
	dir := filepath.Dir(r.path)
	return filepath.Join(dir, "done.txt")
}
