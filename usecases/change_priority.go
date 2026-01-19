package usecases

import (
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// ChangePriority changes the priority of a todo at the specified position
func ChangePriority(repo TodoRepository, m matrix.Matrix, quadrant matrix.QuadrantType, index int, newPriority todo.Priority) (matrix.Matrix, error) {
	// Tell the matrix to change priority at the specified position
	updatedMatrix, changed := m.ChangePriorityAt(quadrant, index, newPriority)

	// Only persist if something changed
	if changed {
		err := saveAllTodos(repo, updatedMatrix)
		if err != nil {
			return m, err // Return original matrix if save fails
		}
	}

	return updatedMatrix, nil
}
