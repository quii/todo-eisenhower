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
	// Let the domain parse user input and create the todo
	newTodo := todotxt.ParseNew(description, priority, time.Now())

	// Add todo to matrix
	updatedMatrix := m.AddTodo(newTodo)

	// Persist changes
	err := saveAllTodos(repo, updatedMatrix)
	if err != nil {
		return m, err // Return original matrix if save fails
	}

	return updatedMatrix, nil
}
