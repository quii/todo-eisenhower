package usecases

import (
	"github.com/quii/todo-eisenhower/domain/matrix"
)

// ArchiveTodo archives a completed todo at the specified position by moving it to the archive file
func ArchiveTodo(repo TodoRepository, m matrix.Matrix, quadrant matrix.QuadrantType, index int) (matrix.Matrix, error) {
	archivedTodo, updatedMatrix, archived := m.ArchiveTodoAt(quadrant, index)

	if !archived {
		return m, nil // No-op if todo is not completed or index is invalid
	}

	// Append to archive file
	err := repo.AppendToArchive(archivedTodo)
	if err != nil {
		return m, err
	}

	// Save the updated todos (with archived todo removed)
	err = saveAllTodos(repo, updatedMatrix)
	if err != nil {
		return m, err
	}

	return updatedMatrix, nil
}
