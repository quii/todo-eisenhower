package usecases

import (
	"time"

	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// ToggleCompletion toggles the completion status of a todo at the specified position.
// When completing a recurring task, it automatically creates a successor with an updated due date.
// The now parameter allows deterministic testing.
func ToggleCompletion(repo TodoRepository, m matrix.Matrix, quadrant matrix.QuadrantType, index int, now time.Time) (matrix.Matrix, error) {
	// Get original todo before toggling (need recurrence + due date info)
	todos := m.GetTodosForQuadrant(quadrant)
	if index < 0 || index >= len(todos) {
		return m, nil
	}
	originalTodo := todos[index]

	updatedMatrix, changed := m.ToggleCompletionAt(quadrant, index, now)
	if !changed {
		return m, nil
	}

	// Check if we're completing (not uncompleting) a recurring task
	isCompleting := !originalTodo.IsCompleted()
	if isCompleting && originalTodo.Recurrence() != nil && originalTodo.DueDate() != nil {
		updatedMatrix = createSuccessor(updatedMatrix, originalTodo, quadrant, index, now)
	}

	err := saveAllTodos(repo, updatedMatrix)
	if err != nil {
		return m, err
	}

	return updatedMatrix, nil
}

// createSuccessor creates a successor task for a recurring todo that was just completed.
// The completed original has its rec: tag stripped, and the new successor gets an updated due date.
func createSuccessor(m matrix.Matrix, originalTodo todo.Todo, quadrant matrix.QuadrantType, index int, now time.Time) matrix.Matrix {
	rec := originalTodo.Recurrence()
	dueDate := originalTodo.DueDate()

	// Strip rec: from the completed original
	completedTodo := m.GetTodosForQuadrant(quadrant)[index]
	completedWithoutRec := completedTodo.WithRecurrence(nil)
	m = m.UpdateTodoAtIndex(quadrant, index, completedWithoutRec)

	// Calculate successor due date
	nextDue := rec.NextDueDate(*dueDate, now, now)

	// Create successor in Schedule quadrant - recurring tasks are future work by definition
	successor := todo.NewFull(
		originalTodo.Description(),
		todo.PriorityB,
		false,
		nil,
		&now,
		&nextDue,
		nil,
		rec,
		originalTodo.Projects(),
		originalTodo.Contexts(),
	)

	return m.AddTodo(successor)
}
