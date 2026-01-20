// Package memory provides an in-memory adapter for testing
package memory

import (
	"bytes"

	"github.com/quii/todo-eisenhower/domain/todo"
	"github.com/quii/todo-eisenhower/domain/todotxt"
)

// Repository is an in-memory implementation of TodoRepository
// Perfect for testing - backed by bytes.Buffer
type Repository struct {
	buffer        *bytes.Buffer
	archiveBuffer *bytes.Buffer
}

// NewRepository creates a new empty in-memory repository
func NewRepository() *Repository {
	return &Repository{
		buffer:        &bytes.Buffer{},
		archiveBuffer: &bytes.Buffer{},
	}
}

// LoadAll reads todos from the buffer
func (r *Repository) LoadAll() ([]todo.Todo, error) {
	return todotxt.Unmarshal(bytes.NewReader(r.buffer.Bytes()))
}

// SaveAll writes todos to the buffer (full rewrite)
func (r *Repository) SaveAll(todos []todo.Todo) error {
	r.buffer.Reset()
	return todotxt.Marshal(r.buffer, todos)
}

// String returns the current buffer contents (useful for test assertions)
func (r *Repository) String() string {
	return r.buffer.String()
}

// AppendToArchive appends a todo to the archive buffer
func (r *Repository) AppendToArchive(t todo.Todo) error {
	return todotxt.Marshal(r.archiveBuffer, []todo.Todo{t})
}

// ArchiveString returns the current archive buffer contents (useful for test assertions)
func (r *Repository) ArchiveString() string {
	return r.archiveBuffer.String()
}
