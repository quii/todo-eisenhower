package usecases

import (
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// ArchiveCompletedInQuadrant archives all completed todos in the specified quadrant
func ArchiveCompletedInQuadrant(repo TodoRepository, m matrix.Matrix, quadrant matrix.QuadrantType) (matrix.Matrix, error) {
	archived, updatedMatrix := m.ArchiveCompletedInQuadrant(quadrant)
	return persistArchivedTodos(repo, m, updatedMatrix, archived)
}

// ArchiveAllCompleted archives all completed todos across all quadrants
func ArchiveAllCompleted(repo TodoRepository, m matrix.Matrix) (matrix.Matrix, error) {
	archived, updatedMatrix := m.ArchiveAllCompleted()
	return persistArchivedTodos(repo, m, updatedMatrix, archived)
}

// persistArchivedTodos handles the common persistence logic for archiving todos
func persistArchivedTodos(repo TodoRepository, original, updated matrix.Matrix, archived []todo.Todo) (matrix.Matrix, error) {
	if len(archived) == 0 {
		return original, nil
	}

	for _, t := range archived {
		if err := repo.AppendToArchive(t); err != nil {
			return original, err
		}
	}

	if err := saveAllTodos(repo, updated); err != nil {
		return original, err
	}

	return updated, nil
}
