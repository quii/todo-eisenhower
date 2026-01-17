// Package file provides file-based adapters for reading and writing todo.txt files.
package file

import (
	"io"
	"os"
)

// Source is an adapter that reads and writes todos to a file
type Source struct {
	path string
}

// NewSource creates a new file-based todo source
func NewSource(path string) Source {
	return Source{path: path}
}

// GetTodos opens the file and returns its contents as an io.ReadCloser
// If the file doesn't exist, it creates an empty one
func (s Source) GetTodos() (io.ReadCloser, error) {
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		//nolint:gosec // G306: todo.txt files are intentionally world-readable (0o644 per todo.txt spec)
		if err := os.WriteFile(s.path, []byte(""), 0o644); err != nil {
			return nil, err
		}
	}
	return os.Open(s.path)
}

// SaveTodo appends a todo line to the file
func (s Source) SaveTodo(line string) error {
	//nolint:gosec // G302: todo.txt files are intentionally world-readable (0o644 per todo.txt spec)
	f, err := os.OpenFile(s.path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	_, err = f.WriteString(line)
	return err
}

// ReplaceAll replaces the entire file content
func (s Source) ReplaceAll(content string) error {
	//nolint:gosec // G306: todo.txt files are intentionally world-readable (0o644 per todo.txt spec)
	return os.WriteFile(s.path, []byte(content), 0o644)
}
