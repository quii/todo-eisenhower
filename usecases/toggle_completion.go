package usecases

import (
	"time"

	"github.com/quii/todo-eisenhower/domain/matrix"
)

// ToggleCompletion toggles the completion status of a todo at the specified position
func ToggleCompletion(repo TodoRepository, m matrix.Matrix, quadrant matrix.QuadrantType, index int) (matrix.Matrix, error) {
	// Tell the matrix to toggle completion at the specified position
	// Use case provides "now" (application concern)
	updatedMatrix, changed := m.ToggleCompletionAt(quadrant, index, time.Now())

	// Only persist if something changed
	if changed {
		err := saveAllTodos(repo, updatedMatrix)
		if err != nil {
			return m, err // Return original matrix if save fails
		}
	}

	return updatedMatrix, nil
}
