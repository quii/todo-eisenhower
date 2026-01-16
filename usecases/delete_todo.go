package usecases

import (
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// DeleteTodo removes a todo from the matrix and persists the change
func DeleteTodo(writer TodoWriter, m matrix.Matrix, todoToDelete todo.Todo) (matrix.Matrix, error) {
	// Remove todo from matrix
	updatedMatrix := m.RemoveTodo(todoToDelete)

	// Persist changes - rewrite entire file
	err := saveAllTodos(writer, updatedMatrix)
	if err != nil {
		return m, err // Return original matrix if save fails
	}

	return updatedMatrix, nil
}
