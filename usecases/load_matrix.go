package usecases

import (
	"io"

	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/parser"
)

// TodoSource is a port for retrieving todo data
type TodoSource interface {
	GetTodos() (io.ReadCloser, error)
}

// LoadMatrix loads todos from the given source and returns an Eisenhower matrix
func LoadMatrix(source TodoSource) (matrix.Matrix, error) {
	reader, err := source.GetTodos()
	if err != nil {
		return matrix.Matrix{}, err
	}
	defer reader.Close()

	todos, err := parser.Parse(reader)
	if err != nil {
		return matrix.Matrix{}, err
	}

	return matrix.New(todos), nil
}
