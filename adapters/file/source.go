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
func (s Source) GetTodos() (io.ReadCloser, error) {
	return os.Open(s.path)
}

// SaveTodo appends a todo line to the file
func (s Source) SaveTodo(line string) error {
	f, err := os.OpenFile(s.path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(line)
	return err
}
