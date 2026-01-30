package usecases

import (
	"time"

	"github.com/quii/todo-eisenhower/domain/matrix"
	"github.com/quii/todo-eisenhower/domain/todo"
)

// ChangePriority changes the priority of a todo at the specified position
// Manages prioritised: tag automatically:
// - Moving TO Priority A: adds prioritised date (current date)
// - Moving FROM Priority A: removes prioritised date
func ChangePriority(repo TodoRepository, m matrix.Matrix, quadrant matrix.QuadrantType, index int, newPriority todo.Priority) (matrix.Matrix, error) {
	// Get the todo at this index to check priorities
	todos := getTodosForQuadrant(m, quadrant)
	if index < 0 || index >= len(todos) {
		return m, nil
	}

	selectedTodo := todos[index]
	oldPriority := selectedTodo.Priority()

	// Use the existing ChangePriorityAt method
	updatedMatrix, changed := m.ChangePriorityAt(quadrant, index, newPriority)
	if !changed {
		return m, nil
	}

	// Now manage the prioritised date if priority changed to/from A
	if needsPrioritisedDateUpdate(oldPriority, newPriority) {
		updatedMatrix = updatePrioritisedDate(updatedMatrix, selectedTodo, oldPriority, newPriority)
	}

	err := saveAllTodos(repo, updatedMatrix)
	if err != nil {
		return m, err
	}

	return updatedMatrix, nil
}

func getTodosForQuadrant(m matrix.Matrix, quadrant matrix.QuadrantType) []todo.Todo {
	switch quadrant {
	case matrix.DoFirstQuadrant:
		return m.DoFirst()
	case matrix.ScheduleQuadrant:
		return m.Schedule()
	case matrix.DelegateQuadrant:
		return m.Delegate()
	case matrix.EliminateQuadrant:
		return m.Eliminate()
	case matrix.BacklogQuadrant:
		return m.Backlog()
	default:
		return []todo.Todo{}
	}
}

func needsPrioritisedDateUpdate(oldPriority, newPriority todo.Priority) bool {
	return (newPriority == todo.PriorityA && oldPriority != todo.PriorityA) ||
		(oldPriority == todo.PriorityA && newPriority != todo.PriorityA)
}

func updatePrioritisedDate(m matrix.Matrix, originalTodo todo.Todo, oldPriority, newPriority todo.Priority) matrix.Matrix {
	// Collect all todos from the matrix (including backlog)
	allTodos := m.AllTodosIncludingBacklog()

	// Find and update the todo that matches original
	var updatedTodos []todo.Todo
	for _, t := range allTodos {
		if t.Description() == originalTodo.Description() &&
			datesEqual(t.CreationDate(), originalTodo.CreationDate()) &&
			t.Priority() == newPriority {

			var prioritisedDate *time.Time
			if newPriority == todo.PriorityA {
				// Moving TO Priority A: add prioritised date
				now := time.Now()
				prioritisedDate = &now
			}
			// else: Moving FROM Priority A: prioritisedDate stays nil

			updatedTodo := todo.NewFull(
				t.Description(),
				t.Priority(),
				t.IsCompleted(),
				t.CompletionDate(),
				t.CreationDate(),
				t.DueDate(),
				prioritisedDate,
				t.Projects(),
				t.Contexts(),
			)
			updatedTodos = append(updatedTodos, updatedTodo)
		} else {
			updatedTodos = append(updatedTodos, t)
		}
	}

	return matrix.New(updatedTodos)
}

func datesEqual(d1, d2 *time.Time) bool {
	if d1 == nil && d2 == nil {
		return true
	}
	if d1 == nil || d2 == nil {
		return false
	}
	return d1.Equal(*d2)
}
