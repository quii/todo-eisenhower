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
	// Parse description to extract clean text and tags
	cleanDescription, projects, contexts := todotxt.ParseDescription(description)

	// Set creation date to now
	now := time.Now()
	creationDate := &now

	// Create the todo using rich domain model with creation date
	var newTodo todo.Todo
	if len(projects) > 0 || len(contexts) > 0 {
		newTodo = todo.NewWithTagsAndDates(cleanDescription, priority, creationDate, projects, contexts)
	} else {
		newTodo = todo.NewWithCreationDate(cleanDescription, priority, creationDate)
	}

	// Add todo to matrix
	updatedMatrix := m.AddTodo(newTodo)

	// Persist changes
	err := saveTodo(repo, newTodo)
	if err != nil {
		return m, err // Return original matrix if save fails
	}

	return updatedMatrix, nil
}
