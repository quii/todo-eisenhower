package usecases

import (
	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// ChangePriority changes the priority of a todo at the specified position
func ChangePriority(source TodoSource, writer TodoWriter, m matrix.Matrix, quadrant matrix.QuadrantType, index int, newPriority todo.Priority) (matrix.Matrix, error) {
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

	// Get current todo and check if priority is already the same
	selectedTodo := todos[index]
	if selectedTodo.Priority() == newPriority {
		return m, nil // No-op if priority unchanged
	}

	// Use rich domain model to change priority
	updatedTodo := selectedTodo.ChangePriority(newPriority)

	// Update matrix with the changed todo
	updatedMatrix := m.UpdateTodoAtIndex(quadrant, index, updatedTodo)

	// Persist changes (implementation detail)
	err := saveAllTodos(writer, updatedMatrix)
	if err != nil {
		return m, err // Return original matrix if save fails
	}

	// Reload matrix from file to get proper quadrant organization
	// (changing priority means todo should move to different quadrant)
	return LoadMatrix(source)
}
