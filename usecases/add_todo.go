// Package usecases contains application use cases and business logic.
package usecases

import (
	"time"

	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todotxt"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// AddTodo creates a new todo and adds it to the matrix
func AddTodo(repo TodoRepository, m matrix.Matrix, description string, priority todo.Priority) (matrix.Matrix, error) {
	newTodo := todotxt.ParseNew(description, priority, time.Now())
	updatedMatrix := m.AddTodo(newTodo)

	err := saveAllTodos(repo, updatedMatrix)
	if err != nil {
		return m, err
	}

	return updatedMatrix, nil
}
