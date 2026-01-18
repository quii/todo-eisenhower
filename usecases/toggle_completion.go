package usecases

import (
	"github.com/quii/todo-eisenhower/domain/matrix"
)

// ToggleCompletion toggles the completion status of a todo at the specified position
func ToggleCompletion(writer TodoWriter, m matrix.Matrix, quadrant matrix.QuadrantType, index int) (matrix.Matrix, error) {
	// Tell the matrix to toggle completion at the specified position
	updatedMatrix, changed := m.ToggleCompletionAt(quadrant, index)

	// Only persist if something changed
	if changed {
		err := saveAllTodos(writer, updatedMatrix)
		if err != nil {
			return m, err // Return original matrix if save fails
		}
	}

	return updatedMatrix, nil
}
