package usecases

import (
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// ChangePriority changes the priority of a todo at the specified position
func ChangePriority(repo TodoRepository, m matrix.Matrix, quadrant matrix.QuadrantType, index int, newPriority todo.Priority) (matrix.Matrix, error) {
	updatedMatrix, changed := m.ChangePriorityAt(quadrant, index, newPriority)

	if changed {
		err := saveAllTodos(repo, updatedMatrix)
		if err != nil {
			return m, err
		}
	}

	return updatedMatrix, nil
}
