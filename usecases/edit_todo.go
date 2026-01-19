package usecases

import (
	"github.com/quii/todo-eisenhower/domain/matrix"
)

// EditTodo updates a todo at the specified position with new description and tags
func EditTodo(repo TodoRepository, m matrix.Matrix, quadrant matrix.QuadrantType, index int, newDescription string) (matrix.Matrix, error) {
	updatedMatrix := m.EditTodo(quadrant, index, newDescription)

	err := saveAllTodos(repo, updatedMatrix)
	if err != nil {
		return m, err
	}

	return updatedMatrix, nil
}
