package usecases

import (
	"strings"

	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// TodoWriter is the interface for writing todos
type TodoWriter interface {
	SaveTodo(line string) error
	ReplaceAll(content string) error
}

// saveTodo appends a single todo to the file (private helper)
func saveTodo(writer TodoWriter, t todo.Todo) error {
	return writer.SaveTodo(t.String())
}

// saveAllTodos writes all todos from the matrix back to the file (private helper)
func saveAllTodos(writer TodoWriter, m matrix.Matrix) error {
	var b strings.Builder

	// Format all todos from all quadrants
	for _, t := range m.AllTodos() {
		b.WriteString(t.String())
	}

	return writer.ReplaceAll(b.String())
}
