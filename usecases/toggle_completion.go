package usecases

import (
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// ToggleCompletion toggles the completion status of a todo at the specified position
func ToggleCompletion(writer TodoWriter, m matrix.Matrix, quadrant matrix.QuadrantType, index int) (matrix.Matrix, error) {
	// Get the todos from the quadrant
	var todos []todo.Todo
	switch quadrant {
	case matrix.DoFirstQuadrant:
		todos = m.DoFirst()
	case matrix.ScheduleQuadrant:
		todos = m.Schedule()
	case matrix.DelegateQuadrant:
		todos = m.Delegate()
	case matrix.EliminateQuadrant:
		todos = m.Eliminate()
	}

	// Validate index
	if index < 0 || index >= len(todos) {
		return m, nil // No-op if invalid index
	}

	// Use rich domain model to toggle completion
	selectedTodo := todos[index]
	updatedTodo := selectedTodo.ToggleCompletion()

	// Update matrix with the changed todo
	updatedMatrix := m.UpdateTodoAtIndex(quadrant, index, updatedTodo)

	// Persist changes (implementation detail)
	err := saveAllTodos(writer, updatedMatrix)
	if err != nil {
		return m, err // Return original matrix if save fails
	}

	return updatedMatrix, nil
}
