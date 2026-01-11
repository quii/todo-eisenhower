package file

import (
	"io"
	"os"
)

// Source is an adapter that reads todos from a file
type Source struct {
	path string
}

// NewSource creates a new file-based todo source
func NewSource(path string) Source {
	return Source{path: path}
}

// GetTodos opens the file and returns its contents as an io.ReadCloser
func (s Source) GetTodos() (io.ReadCloser, error) {
	return os.Open(s.path)
}
