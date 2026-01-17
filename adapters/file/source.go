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
		if err := os.WriteFile(s.path, []byte(""), 0644); err != nil {
			return nil, err
		}
	}
	return os.Open(s.path)
}

// SaveTodo appends a todo line to the file
func (s Source) SaveTodo(line string) error {
	f, err := os.OpenFile(s.path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
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
	return os.WriteFile(s.path, []byte(content), 0644)
}
