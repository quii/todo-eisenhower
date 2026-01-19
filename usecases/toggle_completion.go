package usecases

import (
	"time"

	"github.com/quii/todo-eisenhower/domain/matrix"
)

// ToggleCompletion toggles the completion status of a todo at the specified position
func ToggleCompletion(repo TodoRepository, m matrix.Matrix, quadrant matrix.QuadrantType, index int) (matrix.Matrix, error) {
	updatedMatrix, changed := m.ToggleCompletionAt(quadrant, index, time.Now())

	if changed {
		err := saveAllTodos(repo, updatedMatrix)
		if err != nil {
			return m, err
		}
	}

	return updatedMatrix, nil
}
