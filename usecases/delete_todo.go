package usecases

import (
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// DeleteTodo removes a todo from the matrix and persists the change
func DeleteTodo(repo TodoRepository, m matrix.Matrix, todoToDelete todo.Todo) (matrix.Matrix, error) {
	updatedMatrix := m.RemoveTodo(todoToDelete)

	err := saveAllTodos(repo, updatedMatrix)
	if err != nil {
		return m, err
	}

	return updatedMatrix, nil
}
